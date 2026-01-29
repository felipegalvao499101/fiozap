package wameow

import (
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
)

// Session representa uma sessao WhatsApp
type Session struct {
	Name      string
	Token     string
	Client    *whatsmeow.Client
	Device    *store.Device
	connected bool
	qrCode    string
	mu        sync.RWMutex
}

func (s *Session) IsConnected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connected
}

func (s *Session) setConnected(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connected = v
}

func (s *Session) GetQRCode() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.qrCode
}

func (s *Session) setQRCode(qr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.qrCode = qr
}

func (s *Session) GetPhone() string {
	if s.Client != nil && s.Client.Store.ID != nil {
		return s.Client.Store.ID.User
	}
	return ""
}

func (s *Session) GetPushName() string {
	if s.Client != nil && s.Client.Store.ID != nil {
		return s.Client.Store.PushName
	}
	return ""
}

func (s *Session) GetJID() string {
	if s.Client != nil && s.Client.Store.ID != nil {
		return s.Client.Store.ID.String()
	}
	return ""
}

func (s *Session) GetToken() string {
	return s.Token
}

func (s *Session) GetName() string {
	return s.Name
}
