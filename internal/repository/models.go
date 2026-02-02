package repository

import (
	"database/sql"
	"time"
)

// SessionModel representa uma sessao no banco de dados
type SessionModel struct {
	ID        string
	Name      string
	Token     string
	JID       sql.NullString
	Phone     sql.NullString
	PushName  sql.NullString
	Connected bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetJID retorna JID como string (vazio se null)
func (s *SessionModel) GetJID() string {
	if s.JID.Valid {
		return s.JID.String
	}
	return ""
}

// GetPhone retorna Phone como string (vazio se null)
func (s *SessionModel) GetPhone() string {
	if s.Phone.Valid {
		return s.Phone.String
	}
	return ""
}

// GetPushName retorna PushName como string (vazio se null)
func (s *SessionModel) GetPushName() string {
	if s.PushName.Valid {
		return s.PushName.String
	}
	return ""
}

// NullString converte string para sql.NullString
func NullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
