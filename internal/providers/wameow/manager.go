package wameow

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"

	"fiozap/internal/core"
	"fiozap/internal/integrations/webhook"
	"fiozap/internal/repository"

	"github.com/google/uuid"
	"github.com/mdp/qrterminal/v3"
	"github.com/rs/zerolog"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// Manager gerencia sessoes WhatsApp usando whatsmeow
type Manager struct {
	sessions  map[string]*Session
	mu        sync.RWMutex
	container *sqlstore.Container
	repo      repository.SessionRepository
	webhook   *webhook.Dispatcher
	log       zerolog.Logger
}

// New cria um novo Manager
func New(container *sqlstore.Container, repo repository.SessionRepository, log zerolog.Logger, webhookDispatcher *webhook.Dispatcher) *Manager {
	m := &Manager{
		sessions:  make(map[string]*Session),
		container: container,
		repo:      repo,
		webhook:   webhookDispatcher,
		log:       log.With().Str("component", "wameow").Logger(),
	}
	m.loadSessionsFromDB()
	return m
}

// loadSessionsFromDB carrega sessoes existentes do banco
func (m *Manager) loadSessionsFromDB() {
	sessions, err := m.repo.List(context.Background())
	if err != nil {
		m.log.Error().Err(err).Msg("Failed to load sessions from DB")
		return
	}

	var sessionsToReconnect []string

	for _, s := range sessions {
		// Busca device do whatsmeow se existir JID
		var device *store.Device
		if s.JID.Valid && s.JID.String != "" {
			parsedJID, _ := types.ParseJID(s.JID.String)
			if !parsedJID.IsEmpty() {
				device, _ = m.container.GetDevice(context.Background(), parsedJID)
			}
		}
		if device == nil {
			device = m.container.NewDevice()
		}

		session := &Session{
			ID:     s.ID,
			Name:   s.Name,
			Token:  s.Token,
			Device: device,
		}
		m.sessions[s.Name] = session
		m.log.Info().Str("name", s.Name).Msg("Session loaded from DB")

		// Marca para reconexão se estava conectada e tem JID (já pareada)
		if s.Connected && s.JID.Valid && s.JID.String != "" {
			sessionsToReconnect = append(sessionsToReconnect, s.Name)
		}
	}

	// Reconecta sessões em background
	for _, name := range sessionsToReconnect {
		go m.reconnectSession(name)
	}
}

// reconnectSession tenta reconectar uma sessão
func (m *Manager) reconnectSession(name string) {
	m.log.Info().Str("name", name).Msg("Attempting to reconnect session")
	_, err := m.Connect(context.Background(), name)
	if err != nil {
		m.log.Error().Err(err).Str("name", name).Msg("Failed to reconnect session")
	}
}

func generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// CreateSession cria uma nova sessao
func (m *Manager) CreateSession(ctx context.Context, name string) (core.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[name]; exists {
		return nil, fmt.Errorf("session %s already exists", name)
	}

	session := &Session{
		ID:     uuid.New().String(),
		Name:   name,
		Token:  generateToken(),
		Device: m.container.NewDevice(),
	}

	// Persiste no banco
	model := &repository.SessionModel{
		ID:        session.ID,
		Name:      session.Name,
		Token:     session.Token,
		Connected: false,
	}
	if err := m.repo.Create(ctx, model); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	m.sessions[name] = session
	m.log.Info().Str("name", name).Msg("Session created")
	return session, nil
}

// GetSession retorna uma sessao existente
func (m *Manager) GetSession(name string) (core.Session, error) {
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
func (m *Manager) ListSessions() []core.Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]core.Session, 0, len(m.sessions))
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

	// Remove do banco
	if err := m.repo.Delete(ctx, name); err != nil {
		return fmt.Errorf("failed to delete session from DB: %w", err)
	}

	delete(m.sessions, name)
	m.log.Info().Str("name", name).Msg("Session deleted")
	return nil
}

// Connect conecta uma sessao
func (m *Manager) Connect(ctx context.Context, name string) (core.Session, error) {
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
		// Usa Background context para o QR channel não ser cancelado quando a requisição HTTP terminar
		qrChan, _ := client.GetQRChannel(context.Background())
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
	qrCount := 0
	for evt := range qrChan {
		switch evt.Event {
		case "code":
			qrCount++
			session.setQRCode(evt.Code)
			m.log.Info().Str("name", session.Name).Int("qr_number", qrCount).Msg("QR code received")
			fmt.Printf("\n=== QR Code #%d for session '%s' (expires in ~20s) ===\n", qrCount, session.Name)
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			fmt.Println("====================================================")
		case "timeout":
			session.setQRCode("")
			m.log.Warn().Str("name", session.Name).Int("qr_count", qrCount).Msg("QR code timeout - no more codes will be generated")
		case "success":
			session.setQRCode("")
			m.log.Info().Str("name", session.Name).Msg("QR code scanned successfully")
		}
	}
	m.log.Debug().Str("name", session.Name).Msg("QR channel closed")
}

func (m *Manager) handleEvent(session *Session, evt interface{}) {
	ctx := context.Background()

	switch e := evt.(type) {
	case *events.Connected:
		session.setConnected(true)
		session.setQRCode("")
		m.updateSessionInDB(session)
		m.log.Info().Str("name", session.Name).Msg("Connected")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventConnected, e)

	case *events.PairSuccess:
		if session.Client != nil && session.Client.Store.ID != nil {
			session.mu.Lock()
			session.jid = session.Client.Store.ID.String()
			session.mu.Unlock()
		}
		m.updateSessionInDB(session)
		m.log.Info().Str("name", session.Name).Str("jid", e.ID.String()).Msg("Pair success")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventPairSuccess, e)

	case *events.Disconnected:
		session.setConnected(false)
		m.updateSessionInDB(session)
		m.log.Info().Str("name", session.Name).Msg("Disconnected")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventDisconnected, e)

	case *events.LoggedOut:
		session.setConnected(false)
		m.updateSessionInDB(session)
		m.log.Warn().Str("name", session.Name).Msg("Logged out")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventLoggedOut, e)

	case *events.Message:
		m.log.Debug().Str("name", session.Name).Str("from", e.Info.Sender.String()).Msg("Message received")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventMessage, e)

	case *events.Receipt:
		m.log.Debug().Str("name", session.Name).Strs("ids", e.MessageIDs).Msg("Receipt received")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventReceipt, e)

	case *events.Presence:
		m.log.Debug().Str("name", session.Name).Str("from", e.From.String()).Msg("Presence received")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventPresence, e)

	case *events.ChatPresence:
		m.log.Debug().Str("name", session.Name).Str("chat", e.Chat.String()).Msg("Chat presence received")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventChatPresence, e)

	case *events.HistorySync:
		m.log.Debug().Str("name", session.Name).Msg("History sync received")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventHistorySync, e)

	case *events.GroupInfo:
		m.log.Debug().Str("name", session.Name).Str("group", e.JID.String()).Msg("Group info received")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventGroupInfo, e)

	case *events.JoinedGroup:
		m.log.Debug().Str("name", session.Name).Str("group", e.JID.String()).Msg("Joined group")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventJoinedGroup, e)

	case *events.Picture:
		m.log.Debug().Str("name", session.Name).Str("jid", e.JID.String()).Msg("Picture updated")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventPicture, e)

	case *events.CallOffer:
		m.log.Info().Str("name", session.Name).Str("from", e.From.String()).Msg("Call offer received")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventCallOffer, e)

	case *events.CallAccept:
		m.log.Debug().Str("name", session.Name).Msg("Call accepted")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventCallAccept, e)

	case *events.CallTerminate:
		m.log.Debug().Str("name", session.Name).Msg("Call terminated")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventCallTerminate, e)

	case *events.KeepAliveTimeout:
		m.log.Warn().Str("name", session.Name).Msg("Keep alive timeout")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventKeepAliveTimeout, e)

	case *events.KeepAliveRestored:
		m.log.Info().Str("name", session.Name).Msg("Keep alive restored")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventKeepAliveRestored, e)

	case *events.ConnectFailure:
		m.log.Error().Str("name", session.Name).Str("reason", e.Reason.String()).Msg("Connect failure")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventConnectFailure, e)

	case *events.StreamError:
		m.log.Error().Str("name", session.Name).Str("code", e.Code).Msg("Stream error")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventStreamError, e)

	case *events.TemporaryBan:
		m.log.Warn().Str("name", session.Name).Str("reason", e.String()).Msg("Temporary ban")
		m.webhook.Dispatch(ctx, session.Name, webhook.EventTemporaryBan, e)
	}
}

// updateSessionInDB atualiza os dados da sessao no banco
func (m *Manager) updateSessionInDB(session *Session) {
	jid := session.GetJID()
	phone := session.GetPhone()
	pushName := session.GetPushName()
	connected := session.IsConnected()

	err := m.repo.UpdateConnection(context.Background(), session.Name, connected, jid, phone, pushName)
	if err != nil {
		m.log.Error().Err(err).Str("name", session.Name).Msg("Failed to update session in DB")
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
	// Se não contém @, é apenas um número de telefone
	if !strings.Contains(phone, "@") {
		return types.NewJID(phone, types.DefaultUserServer)
	}
	jid, _ := types.ParseJID(phone)
	return jid
}
