package repository

import "database/sql"

// Repositories agrupa todos os repositories da aplicacao
type Repositories struct {
	Session SessionRepository
}

// New cria todos os repositories
func New(db *sql.DB) *Repositories {
	return &Repositories{
		Session: NewSessionRepository(db),
	}
}
