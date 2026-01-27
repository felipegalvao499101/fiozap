package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

//go:embed upgrades/001_create_sessions.sql
var migration001 string

type Database struct {
	DB        *sql.DB
	Container *sqlstore.Container
}

func New(databaseURL string) (*Database, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run FioZap migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// WhatsApp store container
	container := sqlstore.NewWithDB(db, "postgres", nil)

	if err := container.Upgrade(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to upgrade whatsmeow schema: %w", err)
	}

	return &Database{
		DB:        db,
		Container: container,
	}, nil
}

func runMigrations(db *sql.DB) error {
	migrations := []struct {
		version string
		sql     string
	}{
		{"001_create_sessions", migration001},
	}

	for _, m := range migrations {
		applied, err := isMigrationApplied(db, m.version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		if _, err := db.Exec(m.sql); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.version, err)
		}

		if _, err := db.Exec(`INSERT INTO "schema_migrations" ("version") VALUES ($1) ON CONFLICT DO NOTHING`, m.version); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", m.version, err)
		}
	}

	return nil
}

func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'schema_migrations')`).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM "schema_migrations" WHERE "version" = $1)`, version).Scan(&exists)
	return exists, err
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) Ping(ctx context.Context) error {
	return d.DB.PingContext(ctx)
}
