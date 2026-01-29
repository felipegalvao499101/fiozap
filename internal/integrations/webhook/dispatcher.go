package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Event struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type Dispatcher struct {
	client *http.Client
	logger zerolog.Logger
}

func NewDispatcher(logger zerolog.Logger) *Dispatcher {
	return &Dispatcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger.With().Str("component", "webhook-dispatcher").Logger(),
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, webhookURL string, event Event) error {
	if webhookURL == "" {
		return nil
	}

	event.Timestamp = time.Now().Unix()

	body, err := json.Marshal(event)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to marshal event")
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to create request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		d.logger.Error().Err(err).Str("url", webhookURL).Msg("Failed to send webhook")
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		d.logger.Warn().
			Str("url", webhookURL).
			Int("status", resp.StatusCode).
			Msg("Webhook returned error status")
	} else {
		d.logger.Debug().
			Str("url", webhookURL).
			Int("status", resp.StatusCode).
			Msg("Webhook sent successfully")
	}

	return nil
}

type MessageEventData struct {
	MessageID string `json:"message_id"`
	From      string `json:"from"`
	Chat      string `json:"chat"`
	IsFromMe  bool   `json:"is_from_me"`
	IsGroup   bool   `json:"is_group"`
	Type      string `json:"type"`
	Text      string `json:"text,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

type StatusEventData struct {
	Status string `json:"status"`
	Phone  string `json:"phone,omitempty"`
}
