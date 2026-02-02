package core

import "time"

// Session representa uma sessao de mensageria
type Session interface {
	GetName() string
	GetToken() string
	GetJID() string
	GetPhone() string
	GetPushName() string
	GetQRCode() string
	IsConnected() bool
}

// MessageResponse resposta de envio de mensagem
type MessageResponse struct {
	ID        string
	Timestamp time.Time
}

// PhoneCheck resultado da verificacao de telefone
type PhoneCheck struct {
	Phone        string
	IsOnWhatsApp bool
	JID          string
}

// UserInfo informacoes do usuario
type UserInfo struct {
	Status    string
	PictureID string
	Devices   []string
}

// ProfilePicture foto de perfil
type ProfilePicture struct {
	URL string
	ID  string
}

// GroupInfo informacoes do grupo
type GroupInfo struct {
	JID          string
	Name         string
	Topic        string
	GroupCreated time.Time
	OwnerJID     string
	Participants []GroupParticipant
	IsAnnounce   bool
	IsLocked     bool
}

// GroupParticipant participante do grupo
type GroupParticipant struct {
	JID          string
	IsAdmin      bool
	IsSuperAdmin bool
}
