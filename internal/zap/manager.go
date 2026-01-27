package zap

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/fiozap/fiozap/internal/logger"
	"github.com/rs/zerolog"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type Manager struct {
	sessions  map[string]*Session
	mu        sync.RWMutex
	container *sqlstore.Container
	logger    zerolog.Logger
}

func NewManager(container *sqlstore.Container, log zerolog.Logger) *Manager {
	return &Manager{
		sessions:  make(map[string]*Session),
		container: container,
		logger:    log.With().Str("component", "zap-manager").Logger(),
	}
}

func (m *Manager) CreateSession(ctx context.Context, name string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[name]; exists {
		return nil, fmt.Errorf("session %s already exists", name)
	}

	device := m.container.NewDevice()
	token := generateToken()

	session := &Session{
		Name:      name,
		Token:     token,
		Device:    device,
		Connected: false,
	}

	m.sessions[name] = session
	m.logger.Info().Str("name", name).Msg("Session created")

	return session, nil
}

func (m *Manager) GetSession(name string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[name]
	if !exists {
		return nil, fmt.Errorf("session %s not found", name)
	}

	return session, nil
}

func (m *Manager) GetOrCreateSession(ctx context.Context, name string) (*Session, error) {
	session, err := m.GetSession(name)
	if err == nil {
		return session, nil
	}

	return m.CreateSession(ctx, name)
}

func (m *Manager) ListSessions() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sessions := make([]*Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		sessions = append(sessions, s)
	}
	return sessions
}

func (m *Manager) DeleteSession(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[name]
	if !exists {
		return fmt.Errorf("session %s not found", name)
	}

	if session.Client != nil {
		session.Client.Disconnect()
		if session.Client.Store.ID != nil {
			_ = session.Client.Logout(ctx)
		}
	}

	delete(m.sessions, name)
	m.logger.Info().Str("name", name).Msg("Session deleted")

	return nil
}

func (m *Manager) Connect(ctx context.Context, name string) (*Session, error) {
	session, err := m.GetOrCreateSession(ctx, name)
	if err != nil {
		return nil, err
	}

	if session.Client != nil && session.Client.IsConnected() {
		return session, nil
	}

	waLog := logger.NewWALogger(m.logger, "whatsmeow")
	client := whatsmeow.NewClient(session.Device, waLog)
	session.Client = client

	client.AddEventHandler(func(evt interface{}) {
		handleEvent(session, m.logger, evt)
	})

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(ctx)
		if err := client.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}

		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					session.SetQRCode(evt.Code)
					m.logger.Info().Str("name", name).Msg("QR code received")
				}
			}
		}()
	} else {
		if err := client.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}
	}

	return session, nil
}

func (m *Manager) Disconnect(name string) error {
	session, err := m.GetSession(name)
	if err != nil {
		return err
	}

	if session.Client != nil {
		session.Client.Disconnect()
		session.SetConnected(false)
	}

	return nil
}
