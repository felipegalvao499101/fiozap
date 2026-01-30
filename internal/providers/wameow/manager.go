package wameow

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"fiozap/internal/domain"

	"github.com/mdp/qrterminal/v3"
	"github.com/rs/zerolog"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// Manager gerencia sessoes WhatsApp usando whatsmeow
type Manager struct {
	sessions  map[string]*Session
	mu        sync.RWMutex
	container *sqlstore.Container
	log       zerolog.Logger
}

// New cria um novo Manager
func New(container *sqlstore.Container, log zerolog.Logger) *Manager {
	return &Manager{
		sessions:  make(map[string]*Session),
		container: container,
		log:       log.With().Str("component", "wameow").Logger(),
	}
}

func generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// CreateSession cria uma nova sessao
func (m *Manager) CreateSession(ctx context.Context, name string) (domain.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[name]; exists {
		return nil, fmt.Errorf("session %s already exists", name)
	}

	session := &Session{
		Name:   name,
		Token:  generateToken(),
		Device: m.container.NewDevice(),
	}

	m.sessions[name] = session
	m.log.Info().Str("name", name).Msg("Session created")
	return session, nil
}

// GetSession retorna uma sessao existente
func (m *Manager) GetSession(name string) (domain.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[name]
	if !exists {
		return nil, fmt.Errorf("session %s not found", name)
	}
	return session, nil
}

// getSessionInternal retorna sessao interna (para uso interno)
func (m *Manager) getSessionInternal(name string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[name]
	if !exists {
		return nil, fmt.Errorf("session %s not found", name)
	}
	return session, nil
}

// ListSessions lista todas as sessoes
func (m *Manager) ListSessions() []domain.Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]domain.Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		list = append(list, s)
	}
	return list
}

// DeleteSession remove uma sessao
func (m *Manager) DeleteSession(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[name]
	if !exists {
		return fmt.Errorf("session %s not found", name)
	}

	if session.Client != nil {
		session.Client.Disconnect()
	}

	delete(m.sessions, name)
	m.log.Info().Str("name", name).Msg("Session deleted")
	return nil
}

// Connect conecta uma sessao
func (m *Manager) Connect(ctx context.Context, name string) (domain.Session, error) {
	session, err := m.getSessionInternal(name)
	if err != nil {
		return nil, err
	}

	if session.Client != nil && session.Client.IsConnected() {
		return session, nil
	}

	client := whatsmeow.NewClient(session.Device, nil)
	session.Client = client

	client.AddEventHandler(func(evt interface{}) {
		m.handleEvent(session, evt)
	})

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(ctx)
		if err := client.Connect(); err != nil {
			return nil, fmt.Errorf("connect failed: %w", err)
		}
		go m.handleQR(session, qrChan)
	} else {
		if err := client.Connect(); err != nil {
			return nil, fmt.Errorf("connect failed: %w", err)
		}
	}

	return session, nil
}

func (m *Manager) handleQR(session *Session, qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		if evt.Event == "code" {
			session.setQRCode(evt.Code)
			m.log.Info().Str("name", session.Name).Msg("QR code received")
			fmt.Printf("\n=== QR Code for session '%s' ===\n", session.Name)
			qrterminal.GenerateWithConfig(evt.Code, qrterminal.Config{
				Level:     qrterminal.L,
				Writer:    os.Stdout,
				BlackChar: qrterminal.BLACK,
				WhiteChar: qrterminal.WHITE,
				QuietZone: 1,
			})
			fmt.Println("================================")
		}
	}
}

func (m *Manager) handleEvent(session *Session, evt interface{}) {
	switch evt.(type) {
	case *events.Connected:
		session.setConnected(true)
		session.setQRCode("")
		m.log.Info().Str("name", session.Name).Msg("Connected")
	case *events.Disconnected:
		session.setConnected(false)
		m.log.Info().Str("name", session.Name).Msg("Disconnected")
	case *events.LoggedOut:
		session.setConnected(false)
		m.log.Warn().Str("name", session.Name).Msg("Logged out")
	}
}

// Disconnect desconecta uma sessao
func (m *Manager) Disconnect(name string) error {
	session, err := m.getSessionInternal(name)
	if err != nil {
		return err
	}

	if session.Client != nil {
		session.Client.Disconnect()
		session.setConnected(false)
	}
	return nil
}

// Logout faz logout da sessao (requer novo QR)
func (m *Manager) Logout(ctx context.Context, name string) error {
	session, err := m.getSessionInternal(name)
	if err != nil {
		return err
	}

	if session.Client != nil && session.Client.IsLoggedIn() {
		return session.Client.Logout(ctx)
	}
	return nil
}

// getClient retorna client conectado ou erro
func (m *Manager) getClient(name string) (*whatsmeow.Client, error) {
	session, err := m.getSessionInternal(name)
	if err != nil {
		return nil, err
	}
	if session.Client == nil || !session.Client.IsConnected() {
		return nil, fmt.Errorf("session not connected")
	}
	return session.Client, nil
}

// parseJID converte string para JID
func parseJID(phone string) types.JID {
	jid, _ := types.ParseJID(phone)
	if jid.IsEmpty() {
		return types.NewJID(phone, types.DefaultUserServer)
	}
	return jid
}
