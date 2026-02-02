package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// internalConfig armazena configuracao interna do webhook
type internalConfig struct {
	URL     string
	Events  []EventType
	HMACKey string
}

// Dispatcher gerencia envio de webhooks
type Dispatcher struct {
	client     *http.Client
	logger     zerolog.Logger
	configs    map[string]*internalConfig
	configsMu  sync.RWMutex
	retryCount int
	retryDelay time.Duration
}

// NewDispatcher cria um novo dispatcher de webhooks
func NewDispatcher(logger zerolog.Logger) *Dispatcher {
	return &Dispatcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:     logger.With().Str("component", "webhook").Logger(),
		configs:    make(map[string]*internalConfig),
		retryCount: 3,
		retryDelay: 1 * time.Second,
	}
}

// SetConfig define a configuracao de webhook para uma sessao
func (d *Dispatcher) SetConfig(sessionID string, url string, events []EventType) {
	d.configsMu.Lock()
	defer d.configsMu.Unlock()

	if existing, ok := d.configs[sessionID]; ok {
		existing.URL = url
		existing.Events = events
	} else {
		d.configs[sessionID] = &internalConfig{
			URL:    url,
			Events: events,
		}
	}

	d.logger.Info().
		Str("session", sessionID).
		Str("url", url).
		Int("events_count", len(events)).
		Msg("Webhook configured")
}

// GetConfig retorna a configuracao de webhook de uma sessao
func (d *Dispatcher) GetConfig(sessionID string) *Config {
	d.configsMu.RLock()
	defer d.configsMu.RUnlock()

	cfg, exists := d.configs[sessionID]
	if !exists {
		return nil
	}

	return &Config{
		URL:        cfg.URL,
		Events:     cfg.Events,
		HMACKeySet: cfg.HMACKey != "",
	}
}

// RemoveConfig remove a configuracao de webhook de uma sessao
func (d *Dispatcher) RemoveConfig(sessionID string) {
	d.configsMu.Lock()
	defer d.configsMu.Unlock()
	delete(d.configs, sessionID)
	d.logger.Info().Str("session", sessionID).Msg("Webhook config removed")
}

// SetHMACKey define a chave HMAC para uma sessao
func (d *Dispatcher) SetHMACKey(sessionID string, key string) error {
	if len(key) < 32 {
		return fmt.Errorf("HMAC key must be at least 32 characters")
	}

	d.configsMu.Lock()
	defer d.configsMu.Unlock()

	cfg, exists := d.configs[sessionID]
	if !exists {
		cfg = &internalConfig{}
		d.configs[sessionID] = cfg
	}
	cfg.HMACKey = key

	d.logger.Info().Str("session", sessionID).Msg("HMAC key configured")
	return nil
}

// RemoveHMACKey remove a chave HMAC de uma sessao
func (d *Dispatcher) RemoveHMACKey(sessionID string) {
	d.configsMu.Lock()
	defer d.configsMu.Unlock()

	if cfg, exists := d.configs[sessionID]; exists {
		cfg.HMACKey = ""
	}

	d.logger.Info().Str("session", sessionID).Msg("HMAC key removed")
}

// Dispatch envia um evento raw do whatsmeow para o webhook configurado
func (d *Dispatcher) Dispatch(ctx context.Context, sessionID string, eventType EventType, rawEvent interface{}) {
	d.configsMu.RLock()
	cfg := d.configs[sessionID]
	d.configsMu.RUnlock()

	if cfg == nil || cfg.URL == "" {
		return
	}

	if !d.isSubscribed(cfg, eventType) {
		d.logger.Debug().
			Str("session", sessionID).
			Str("event", string(eventType)).
			Msg("Event not subscribed, skipping webhook")
		return
	}

	event := Event{
		Type:      eventType,
		SessionID: sessionID,
		Event:     rawEvent,
	}

	go d.sendWithRetry(ctx, cfg, event)
}

// isSubscribed verifica se o evento esta na lista de eventos subscritos
func (d *Dispatcher) isSubscribed(cfg *internalConfig, eventType EventType) bool {
	if len(cfg.Events) == 0 {
		return false
	}

	for _, e := range cfg.Events {
		if e == EventAll || e == eventType {
			return true
		}
	}
	return false
}

// sendWithRetry envia o webhook com retentativas
func (d *Dispatcher) sendWithRetry(ctx context.Context, cfg *internalConfig, event Event) {
	var lastErr error

	for attempt := 0; attempt < d.retryCount; attempt++ {
		if attempt > 0 {
			time.Sleep(d.retryDelay * time.Duration(attempt))
		}

		err := d.send(ctx, cfg, event)
		if err == nil {
			return
		}

		lastErr = err
		d.logger.Warn().
			Err(err).
			Int("attempt", attempt+1).
			Int("maxRetries", d.retryCount).
			Str("url", cfg.URL).
			Msg("Webhook failed, retrying")
	}

	d.logger.Error().
		Err(lastErr).
		Str("url", cfg.URL).
		Str("session", event.SessionID).
		Str("event", string(event.Type)).
		Msg("Webhook failed after all retries")
}

// send envia o webhook
func (d *Dispatcher) send(ctx context.Context, cfg *internalConfig, event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "FioZap-Webhook/1.0")

	if cfg.HMACKey != "" {
		signature := d.computeHMAC(body, cfg.HMACKey)
		req.Header.Set("X-Hub-Signature-256", "sha256="+signature)
		req.Header.Set("X-HMAC-Signature", signature)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(respBody))
	}

	d.logger.Debug().
		Str("url", cfg.URL).
		Int("status", resp.StatusCode).
		Str("event", string(event.Type)).
		Msg("Webhook sent successfully")

	return nil
}

// computeHMAC calcula a assinatura HMAC-SHA256
func (d *Dispatcher) computeHMAC(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
