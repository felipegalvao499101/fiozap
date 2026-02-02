package repository

import (
	"context"
	"database/sql"
)

// SessionRepository define operacoes de persistencia de sessoes
type SessionRepository interface {
	Create(ctx context.Context, session *SessionModel) error
	GetByName(ctx context.Context, name string) (*SessionModel, error)
	GetByToken(ctx context.Context, token string) (*SessionModel, error)
	List(ctx context.Context) ([]*SessionModel, error)
	Update(ctx context.Context, session *SessionModel) error
	Delete(ctx context.Context, name string) error
	UpdateConnection(ctx context.Context, name string, connected bool, jid, phone, pushName string) error
}

// sessionRepository implementa SessionRepository usando PostgreSQL
type sessionRepository struct {
	db *sql.DB
}

// NewSessionRepository cria um novo SessionRepository
func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *SessionModel) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO "sessions" ("id", "name", "token", "jid", "phone", "pushName", "connected")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, session.ID, session.Name, session.Token, session.JID, session.Phone, session.PushName, session.Connected)
	return err
}

func (r *sessionRepository) GetByName(ctx context.Context, name string) (*SessionModel, error) {
	session := &SessionModel{}
	err := r.db.QueryRowContext(ctx, `
		SELECT "id", "name", "token", "jid", "phone", "pushName", "connected", "createdAt", "updatedAt"
		FROM "sessions" WHERE "name" = $1
	`, name).Scan(
		&session.ID, &session.Name, &session.Token, &session.JID,
		&session.Phone, &session.PushName, &session.Connected,
		&session.CreatedAt, &session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*SessionModel, error) {
	session := &SessionModel{}
	err := r.db.QueryRowContext(ctx, `
		SELECT "id", "name", "token", "jid", "phone", "pushName", "connected", "createdAt", "updatedAt"
		FROM "sessions" WHERE "token" = $1
	`, token).Scan(
		&session.ID, &session.Name, &session.Token, &session.JID,
		&session.Phone, &session.PushName, &session.Connected,
		&session.CreatedAt, &session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

func (r *sessionRepository) List(ctx context.Context) ([]*SessionModel, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT "id", "name", "token", "jid", "phone", "pushName", "connected", "createdAt", "updatedAt"
		FROM "sessions"
		ORDER BY "createdAt" ASC
	`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var sessions []*SessionModel
	for rows.Next() {
		s := &SessionModel{}
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Token, &s.JID,
			&s.Phone, &s.PushName, &s.Connected,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

func (r *sessionRepository) Update(ctx context.Context, session *SessionModel) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE "sessions" SET 
			"jid" = $1, 
			"phone" = $2, 
			"pushName" = $3, 
			"connected" = $4,
			"updatedAt" = CURRENT_TIMESTAMP
		WHERE "name" = $5
	`, session.JID, session.Phone, session.PushName, session.Connected, session.Name)
	return err
}

func (r *sessionRepository) Delete(ctx context.Context, name string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM "sessions" WHERE "name" = $1`, name)
	return err
}

func (r *sessionRepository) UpdateConnection(ctx context.Context, name string, connected bool, jid, phone, pushName string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE "sessions" SET 
			"jid" = $1, 
			"phone" = $2, 
			"pushName" = $3, 
			"connected" = $4,
			"updatedAt" = CURRENT_TIMESTAMP
		WHERE "name" = $5
	`, NullString(jid), NullString(phone), NullString(pushName), connected, name)
	return err
}
