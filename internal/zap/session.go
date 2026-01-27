package zap

import (
	"sync"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
)

type Session struct {
	Name      string
	Token     string
	Client    *whatsmeow.Client
	Device    *store.Device
	Connected bool
	QRCode    string
	mu        sync.RWMutex
}

func (s *Session) IsConnected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Connected
}

func (s *Session) SetConnected(connected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Connected = connected
}

func (s *Session) GetQRCode() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.QRCode
}

func (s *Session) SetQRCode(qr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.QRCode = qr
}

func (s *Session) GetPhone() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.Client != nil && s.Client.Store.ID != nil {
		return s.Client.Store.ID.User
	}
	return ""
}

func (s *Session) GetPushName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.Client != nil && s.Client.Store.ID != nil {
		return s.Client.Store.PushName
	}
	return ""
}

func (s *Session) GetJID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.Client != nil && s.Client.Store.ID != nil {
		return s.Client.Store.ID.String()
	}
	return ""
}
