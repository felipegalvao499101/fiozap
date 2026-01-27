-- 001_create_sessions.sql
-- Tabela de sessoes do FioZap

CREATE TABLE IF NOT EXISTS "sessions" (
    "id" VARCHAR(255) PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "token" VARCHAR(255) NOT NULL,
    "jid" VARCHAR(255),
    "phone" VARCHAR(50),
    "pushName" VARCHAR(255),
    "connected" BOOLEAN NOT NULL DEFAULT FALSE,
    "createdAt" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_sessions_name" ON "sessions"("name");
CREATE INDEX IF NOT EXISTS "idx_sessions_token" ON "sessions"("token");
CREATE INDEX IF NOT EXISTS "idx_sessions_connected" ON "sessions"("connected");

-- Tabela de controle de migrations
CREATE TABLE IF NOT EXISTS "schema_migrations" (
    "version" VARCHAR(255) PRIMARY KEY,
    "appliedAt" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
