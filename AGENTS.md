# FioZap

API REST multi-session para WhatsApp usando whatsmeow, com integracao Chatwoot.
Baseado na WuzAPI (https://github.com/asternic/wuzapi) - mesma convencao de campos JSON PascalCase.

## Core Commands

- Build: `go build ./...`
- Run: `go run ./cmd/server/`
- Test: `go test ./...`
- Lint: `golangci-lint run`

## Project Layout

```
fiozap/
├── cmd/server/          → Entrypoint da aplicacao
├── internal/
│   ├── api/
│   │   ├── dto/         → Request/Response structs
│   │   ├── handlers/    → HTTP handlers (session, message, group, chat)
│   │   ├── middleware/  → Middlewares (auth)
│   │   └── router/      → Configuracao de rotas Chi
│   ├── config/          → Configuracao via env vars
│   ├── database/        → PostgreSQL + migrations
│   │   └── upgrades/    → SQL migrations
│   ├── logger/          → zerolog + whatsmeow adapter
│   ├── webhook/         → Dispatcher de webhooks (futuro)
│   ├── chatwoot/        → Integracao Chatwoot (futuro)
│   └── zap/             → WhatsApp Manager (whatsmeow wrapper)
├── docker-compose.yml   → Postgres, Redis, NATS, Chatwoot
└── Dockerfile
```

## JSON Conventions (WuzAPI Style)

**IMPORTANTE:** Todos os campos JSON usam PascalCase para compatibilidade com WuzAPI.

```json
{
  "code": 200,
  "success": true,
  "data": {
    "Phone": "5511999999999",
    "Body": "Hello World",
    "JID": "5511999999999@s.whatsapp.net"
  }
}
```

## Database Schema

### Tabela `sessions`

| Coluna      | Tipo         | Descricao                     |
|-------------|--------------|-------------------------------|
| `id`        | VARCHAR(255) | Primary key (UUID)            |
| `name`      | VARCHAR(255) | Identificador unico da sessao |
| `token`     | VARCHAR(255) | Token de autenticacao         |
| `jid`       | VARCHAR(255) | WhatsApp JID                  |
| `phone`     | VARCHAR(50)  | Numero de telefone            |
| `pushName`  | VARCHAR(255) | Nome no WhatsApp              |
| `connected` | BOOLEAN      | Status de conexao             |
| `createdAt` | TIMESTAMP    | Data de criacao               |
| `updatedAt` | TIMESTAMP    | Data de atualizacao           |

## Autenticacao

Dois niveis de autenticacao via header `Authorization`:

| Token | Acesso |
|-------|--------|
| **GLOBAL_API_TOKEN** (env) | Todas as rotas, todas as sessoes |
| **Session Token** (DB) | Apenas rotas da sessao especifica |

```bash
# Usando token global
curl -H "Authorization: seu_token_global" http://localhost:8080/sessions

# Usando token da sessao
curl -H "Authorization: token_da_sessao" http://localhost:8080/sessions/minha-sessao
```

## API Endpoints

### Session (requer GLOBAL_API_TOKEN)
- `POST /sessions` - Criar sessao `{"Name": "session1"}` → retorna Token
- `GET /sessions` - Listar sessoes

### Session (aceita GLOBAL ou Session Token)
- `GET /sessions/:name` - Status da sessao
- `GET /sessions/:name/qr` - QR code (`?format=image` para PNG)
- `POST /sessions/:name/connect` - Conectar
- `POST /sessions/:name/disconnect` - Desconectar (mantem sessao)
- `POST /sessions/:name/logout` - Logout (remove sessao, requer QR)
- `DELETE /sessions/:name` - Remover sessao

### Messages
- `POST /sessions/:name/messages/text` - `{"Phone":"...", "Body":"..."}`
- `POST /sessions/:name/messages/image` - `{"Phone":"...", "Image":"base64...", "Caption":"..."}`
- `POST /sessions/:name/messages/video` - `{"Phone":"...", "Video":"base64...", "Caption":"..."}`
- `POST /sessions/:name/messages/audio` - `{"Phone":"...", "Audio":"base64..."}`
- `POST /sessions/:name/messages/document` - `{"Phone":"...", "Document":"base64...", "FileName":"..."}`
- `POST /sessions/:name/messages/sticker` - `{"Phone":"...", "Sticker":"base64..."}`
- `POST /sessions/:name/messages/location` - `{"Phone":"...", "Latitude":..., "Longitude":..., "Name":"..."}`
- `POST /sessions/:name/messages/contact` - `{"Phone":"...", "Name":"...", "Vcard":"..."}`
- `POST /sessions/:name/messages/poll` - `{"Phone":"...", "Question":"...", "Options":[...]}`
- `POST /sessions/:name/messages/reaction` - `{"Phone":"...", "Id":"msgId", "Body":"❤️"}`
- `PUT /sessions/:name/messages/:messageId` - `{"Phone":"...", "Body":"new text"}`
- `DELETE /sessions/:name/messages/:messageId` - `{"Phone":"..."}`

### Users
- `POST /sessions/:name/users/check` - `{"Phone":["5511...","5521..."]}`
- `GET /sessions/:name/users/:phone` - Info do usuario
- `GET /sessions/:name/users/:phone/avatar` - Foto de perfil

### Groups
- `POST /sessions/:name/groups` - `{"Name":"...", "Participants":[...]}`
- `GET /sessions/:name/groups` - Listar grupos
- `GET /sessions/:name/groups/:groupJid` - Info do grupo
- `PUT /sessions/:name/groups/:groupJid/name` - `{"Name":"..."}`
- `PUT /sessions/:name/groups/:groupJid/topic` - `{"Topic":"..."}`
- `PUT /sessions/:name/groups/:groupJid/photo` - `{"Image":"base64..."}`
- `POST /sessions/:name/groups/:groupJid/leave`
- `GET /sessions/:name/groups/:groupJid/invite` - Link de convite
- `POST /sessions/:name/groups/:groupJid/invite/revoke` - Revogar link
- `POST /sessions/:name/groups/join` - `{"Code":"..."}`
- `GET /sessions/:name/groups/invite/:code` - Info do convite
- `POST /sessions/:name/groups/:groupJid/participants` - `{"Phone":[...]}`
- `DELETE /sessions/:name/groups/:groupJid/participants` - `{"Phone":[...]}`
- `POST /sessions/:name/groups/:groupJid/participants/promote` - `{"Phone":[...]}`
- `POST /sessions/:name/groups/:groupJid/participants/demote` - `{"Phone":[...]}`
- `PUT /sessions/:name/groups/:groupJid/settings/announce` - `{"Value":true}`
- `PUT /sessions/:name/groups/:groupJid/settings/locked` - `{"Value":true}`
- `PUT /sessions/:name/groups/:groupJid/settings/approval` - `{"Value":true}`

### Chat
- `POST /sessions/:name/chat/markread` - `{"Id":["msg1","msg2"], "ChatPhone":"..."}`
- `POST /sessions/:name/chat/presence` - `{"Phone":"...", "State":"composing", "Media":"audio"}`
- `PUT /sessions/:name/chat/:chatJid/disappearing` - `{"Duration":"24h|7d|90d|off"}`

### Presence (global)
- `POST /sessions/:name/presence` - `{"Type":"available|unavailable"}`
- `POST /sessions/:name/presence/subscribe` - `{"Phone":"..."}`

### Blocklist
- `GET /sessions/:name/blocklist` - Lista de bloqueados
- `POST /sessions/:name/blocklist/block` - `{"Phone":"..."}`
- `POST /sessions/:name/blocklist/unblock` - `{"Phone":"..."}`

## Modulos e Responsabilidades

### `internal/api/dto/`
**Convencoes PascalCase (WuzAPI style):**
- Campos JSON: `Phone`, `Body`, `Caption`, `Image`, `Video`, `Document`, etc.
- Request structs terminam com `Request`
- Response structs terminam com `Response`

### `internal/api/handlers/`
**Arquivos:**
- `session.go` - CRUD de sessoes, QR, connect/disconnect/logout
- `message.go` - Envio de mensagens, check phone, user info
- `group.go` - CRUD de grupos, participantes, settings
- `chat.go` - MarkRead, Presence, Disappearing, Blocklist

### `internal/zap/`
**Arquivos:**
- `session.go` - Struct Session com estado da conexao
- `manager.go` - Gerenciador de sessoes (CRUD, Connect, Disconnect)
- `send.go` - Metodos de envio + Presence + Blocklist
- `group.go` - Metodos de grupo
- `chat.go` - MarkRead, Typing, Disappearing
- `events.go` - Handler de eventos whatsmeow

## Development Patterns

### JSON Response Format
```json
{
  "code": 200,
  "success": true,
  "data": {...}
}
```

### Error Response Format
```json
{
  "code": 400,
  "success": false,
  "error": "missing Phone in Payload"
}
```

### Coding Style
- Go 1.22+
- Usar `internal/` para codigo privado
- DTOs em `internal/api/dto/`
- Logs com zerolog
- Erros wrappados com `fmt.Errorf("context: %w", err)`

### Database
- PostgreSQL com pgx driver
- Migrations em `internal/database/upgrades/`
- Colunas em camelCase com aspas duplas
- whatsmeow usa suas proprias tabelas (prefixo `whatsmeow_`)

## Environment Variables

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
GLOBAL_API_TOKEN=seu_token_super_secreto
DATABASE_URL=postgres://fiozap:fiozap123@localhost:5432/fiozap?sslmode=disable
REDIS_URL=redis://localhost:6379
LOG_LEVEL=debug
LOG_FORMAT=console
WA_DEBUG=false
```

## Docker Services

```bash
# Iniciar infra
docker compose up -d postgres redis

# Iniciar tudo (inclui Chatwoot)
docker compose up -d
```

| Servico         | Porta | Descricao              |
|-----------------|-------|------------------------|
| postgres        | 5432  | PostgreSQL + pgvector  |
| redis           | 6379  | Cache                  |
| nats            | 4222  | Mensageria             |
| dbgate          | 3000  | UI para PostgreSQL     |
| webhook-tester  | 8081  | Testar webhooks        |
| chatwoot        | 3001  | Plataforma atendimento |

## Gotchas

- whatsmeow requer contexto em quase todos os metodos
- QR code expira rapido, cliente deve fazer polling
- Sessoes persistidas no PostgreSQL via sqlstore
- Sempre usar aspas duplas para colunas camelCase no SQL
- Todos os campos JSON usam PascalCase (compatibilidade WuzAPI)
