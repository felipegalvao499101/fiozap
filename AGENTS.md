<coding_guidelines>
# FioZap

API REST multi-session para WhatsApp usando whatsmeow, com integracao Chatwoot.
Baseado na WuzAPI (https://github.com/asternic/wuzapi) - mesma convencao de campos JSON PascalCase.

## Core Commands (Makefile)

```bash
make build         # Compila binario em ./bin/fiozap
make run           # Executa a aplicacao
make test          # Roda os testes
make lint          # Roda golangci-lint
make swagger       # Gera documentacao Swagger
make dev           # Gera swagger e executa
make docker-up     # Sobe containers Docker
make clean         # Remove arquivos gerados
make help          # Lista todos os comandos
```

Comandos diretos (sem Makefile):
- Build: `go build -o bin/fiozap ./cmd/server`
- Run: `go run ./cmd/server/`
- Test: `go test ./...`
- Lint: `golangci-lint run`
- Swagger: `swag init -g cmd/server/main.go -o docs`

## Project Layout

```
fiozap/
├── cmd/server/              → Entrypoint da aplicacao
├── bin/                     → Binarios compilados (gitignore)
├── docs/                    → Swagger docs (gerado)
├── internal/
│   ├── api/
│   │   ├── auth/            → Middleware de autenticacao
│   │   ├── dto/             → Request/Response structs
│   │   ├── handlers/        → HTTP handlers
│   │   └── router/          → Configuracao de rotas Chi
│   ├── config/              → Configuracao via env vars
│   ├── database/            → PostgreSQL + migrations + whatsmeow container
│   │   └── upgrades/        → SQL migrations
│   ├── domain/              → Interfaces e tipos compartilhados
│   ├── integrations/        → Integracoes externas
│   │   ├── chatwoot/        → Integracao Chatwoot (futuro)
│   │   └── webhook/         → Dispatcher de webhooks (futuro)
│   ├── logger/              → zerolog + whatsmeow adapter
│   ├── repository/          → Camada de persistencia (Repository Pattern)
│   └── providers/           → Implementacoes de mensageria
│       └── wameow/          → Provider whatsmeow (nao-oficial)
├── docker-compose.yml       → Postgres, Redis, NATS, Chatwoot
├── Dockerfile
└── Makefile                 → Comandos de build, test, lint, swagger
```

## Architecture

### Domain Layer (`internal/domain/`)

Interface-based design para desacoplar handlers de implementacoes especificas.

**`provider.go`** - Interface Provider para provedores de mensageria:
```go
type Provider interface {
    // Session management
    CreateSession(ctx, name) (Session, error)
    GetSession(name) (Session, error)
    ListSessions() []Session
    Connect(ctx, name) (Session, error)
    Disconnect(name) error
    Logout(ctx, name) error
    
    // Messages, Chat, Groups, Users...
}
```

**`types.go`** - Tipos compartilhados:
- `Session` interface (GetName, GetToken, GetJID, GetQRCode, IsConnected)
- `MessageResponse`, `GroupInfo`, `GroupParticipant`
- `PhoneCheck`, `UserInfo`, `ProfilePicture`

### Providers (`internal/providers/`)

Implementacoes da interface `domain.Provider`.

**`wameow/`** - Provider whatsmeow (API nao-oficial):
- `manager.go` - Gerenciador de sessoes
- `session.go` - Struct Session com estado
- `message.go` - Envio de mensagens
- `group.go` - Operacoes de grupos
- `chat.go` - MarkRead, Typing, Disappearing
- `user.go` - CheckPhone, UserInfo, Blocklist

### Handlers (`internal/api/handlers/`)

Handlers HTTP organizados por dominio:

| Handler | Arquivo | Responsabilidade |
|---------|---------|------------------|
| **SessionHandler** | `session.go` | CRUD sessoes, QR, connect/disconnect |
| **MessageHandler** | `message.go` | Envio de mensagens (text, image, video, etc) |
| **ContactHandler** | `contact.go` | CheckPhone, GetInfo, GetAvatar |
| **GroupHandler** | `group.go` | CRUD grupos, participantes, settings |
| **ChatHandler** | `chat.go` | MarkRead, Presence, Disappearing |
| **BlocklistHandler** | `blocklist.go` | Block/Unblock contatos |
| **CallHandler** | `call.go` | RejectCall |
| **NewsletterHandler** | `newsletter.go` | Canais (not implemented) |
| **PrivacyHandler** | `privacy.go` | Privacidade (not implemented) |
| **ProfileHandler** | `profile.go` | Perfil (not implemented) |

### Repository Layer (`internal/repository/`)

Camada de persistencia usando Repository Pattern para desacoplar acesso ao banco.

**Arquivos:**
- `repository.go` - Agregador `Repositories` struct (facilita injecao de dependencias)
- `models.go` - Models do banco (`SessionModel`) + helpers (`NullString`, getters)
- `session.go` - `SessionRepository` interface + implementacao PostgreSQL

**SessionRepository interface:**
```go
type SessionRepository interface {
    Create(ctx, session) error
    GetByName(ctx, name) (*SessionModel, error)
    GetByToken(ctx, token) (*SessionModel, error)
    List(ctx) ([]*SessionModel, error)
    Update(ctx, session) error
    Delete(ctx, name) error
    UpdateConnection(ctx, name, connected, jid, phone, pushName) error
}
```

**Uso no main.go:**
```go
repos := repository.New(db.DB)
provider := wameow.New(db.Container, repos.Session, log)
```

### Integrations (`internal/integrations/`)

**`webhook/`** - Dispatcher de webhooks (futuro)
**`chatwoot/`** - Integracao Chatwoot (futuro)

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

## Swagger / OpenAPI

Documentacao interativa disponivel em `/swagger/index.html` quando o servidor esta rodando.

```bash
# Gerar/atualizar docs
make swagger

# Arquivos gerados em docs/
docs/
├── docs.go        # Go code para inicializacao
├── swagger.json   # OpenAPI spec (JSON)
└── swagger.yaml   # OpenAPI spec (YAML)
```

**Anotacoes nos handlers:** Todos os endpoints estao documentados com anotacoes `@Summary`, `@Description`, `@Tags`, `@Param`, `@Success`, `@Failure`, `@Security`, `@Router`.

**Tags organizadas por dominio:** sessions, messages, contacts, groups, chat, presence, blocklist, calls, newsletters, privacy, profile, community.

## API Endpoints

### Health
- `GET /health` - Health check
- `GET /swagger/*` - Swagger UI

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
- `POST /sessions/:name/messages/image` - `{"Phone":"...", "Image":"base64...", "Caption":"...", "Mimetype":"..."}`
- `POST /sessions/:name/messages/video` - `{"Phone":"...", "Video":"base64...", "Caption":"...", "Mimetype":"..."}`
- `POST /sessions/:name/messages/audio` - `{"Phone":"...", "Audio":"base64...", "Mimetype":"..."}`
- `POST /sessions/:name/messages/document` - `{"Phone":"...", "Document":"base64...", "FileName":"...", "Mimetype":"..."}`
- `POST /sessions/:name/messages/sticker` - `{"Phone":"...", "Sticker":"base64...", "Mimetype":"..."}`
- `POST /sessions/:name/messages/location` - `{"Phone":"...", "Latitude":..., "Longitude":..., "Name":"...", "Address":"..."}`
- `POST /sessions/:name/messages/contact` - `{"Phone":"...", "Name":"...", "Vcard":"..."}`
- `POST /sessions/:name/messages/poll` - `{"Phone":"...", "Question":"...", "Options":[...], "MultiSelect":false}`
- `POST /sessions/:name/messages/reaction` - `{"Phone":"...", "Id":"msgId", "Body":"❤️"}`
- `PUT /sessions/:name/messages/:messageId` - `{"Phone":"...", "Body":"new text"}`
- `DELETE /sessions/:name/messages/:messageId` - `{"Phone":"..."}`

### Contacts
- `POST /sessions/:name/contacts/check` - `{"Phone":["5511...","5521..."]}`
- `GET /sessions/:name/contacts/:phone` - Info do contato
- `GET /sessions/:name/contacts/:phone/avatar` - Foto de perfil
- `GET /sessions/:name/contacts/:phone/business` - Perfil comercial (not implemented)

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
- `PUT /sessions/:name/groups/:groupJid/settings/approval` - `{"Value":true}` (not implemented)

### Group Request Participants (approval mode)
- `GET /sessions/:name/groups/:groupJid/requests` - Listar solicitacoes (not implemented)
- `POST /sessions/:name/groups/:groupJid/requests/approve` - `{"Phone":[...]}` (not implemented)
- `POST /sessions/:name/groups/:groupJid/requests/reject` - `{"Phone":[...]}` (not implemented)
- `PUT /sessions/:name/groups/:groupJid/settings/memberadd` - `{"Mode":"all_member_add|admin_add"}` (not implemented)

### Community
- `POST /sessions/:name/community/link` - `{"ParentJID":"...", "ChildJID":"..."}` (not implemented)
- `POST /sessions/:name/community/unlink` - `{"ParentJID":"...", "ChildJID":"..."}` (not implemented)
- `GET /sessions/:name/community/:communityJid/subgroups` - Listar subgrupos (not implemented)
- `GET /sessions/:name/community/:communityJid/participants` - Participantes (not implemented)

### Chat
- `POST /sessions/:name/chat/markread` - `{"Id":["msg1","msg2"], "ChatPhone":"..."}`
- `POST /sessions/:name/chat/presence` - `{"Phone":"...", "State":"composing|paused", "Media":"audio"}`
- `PUT /sessions/:name/chat/:chatJid/disappearing` - `{"Duration":"24h|7d|90d|off"}`

### Presence
- `POST /sessions/:name/presence` - `{"Type":"available|unavailable"}`
- `POST /sessions/:name/presence/subscribe` - `{"Phone":"..."}`

### Blocklist
- `GET /sessions/:name/blocklist` - Lista de bloqueados
- `POST /sessions/:name/blocklist/block` - `{"Phone":"..."}`
- `POST /sessions/:name/blocklist/unblock` - `{"Phone":"..."}`

### Newsletter (Channels)
- `POST /sessions/:name/newsletters` - Criar canal (not implemented)
- `GET /sessions/:name/newsletters` - Listar canais (not implemented)
- `GET /sessions/:name/newsletters/:newsletterJid` - Info do canal (not implemented)
- `POST /sessions/:name/newsletters/:newsletterJid/follow` - Seguir (not implemented)
- `POST /sessions/:name/newsletters/:newsletterJid/unfollow` - Deixar de seguir (not implemented)
- `PUT /sessions/:name/newsletters/:newsletterJid/mute` - `{"Mute":true}` (not implemented)
- `POST /sessions/:name/newsletters/:newsletterJid/reaction` - `{"ServerID":"...", "Reaction":"👍"}` (not implemented)

### Privacy
- `GET /sessions/:name/privacy` - Configuracoes de privacidade (not implemented)
- `PUT /sessions/:name/privacy` - `{"Name":"lastSeen", "Value":"all"}` (not implemented)
- `GET /sessions/:name/privacy/status` - Privacidade do status (not implemented)

### Profile
- `GET /sessions/:name/profile/qrlink` - Link QR do contato (not implemented)
- `POST /sessions/:name/profile/qrlink/resolve` - `{"Link":"..."}` (not implemented)
- `PUT /sessions/:name/profile/status` - `{"Status":"..."}` (not implemented)
- `POST /sessions/:name/profile/business/resolve` - `{"Link":"..."}` (not implemented)

### Calls
- `POST /sessions/:name/calls/reject` - `{"CallFrom":"...", "CallID":"..."}` (not implemented)

## Modulos e Responsabilidades

### `internal/api/dto/`
**Convencoes PascalCase (WuzAPI style):**
- Campos JSON: `Phone`, `Body`, `Caption`, `Image`, `Video`, `Document`, etc.
- Request structs terminam com `Request`
- Response structs terminam com `Response`

**Arquivos:**
- `session.go` - CreateSessionRequest, SessionResponse, QRResponse
- `message.go` - SendTextRequest, SendImageRequest, MessageResponse, etc.
- `contact.go` - CheckPhoneRequest, UserInfoResponse, AvatarResponse
- `group.go` - CreateGroupRequest, GroupResponse, ParticipantsRequest, etc.
- `chat.go` - MarkReadRequest, ChatPresenceRequest, DisappearingRequest
- `blocklist.go` - BlocklistResponse, BlockRequest, BlockActionResponse
- `call.go` - RejectCallRequest, CallActionResponse
- `newsletter.go` - CreateNewsletterRequest, NewsletterResponse, etc.
- `privacy.go` - PrivacySettingsResponse, SetPrivacyRequest, etc.
- `profile.go` - SetStatusMessageRequest, BusinessProfileResponse, etc.
- `community.go` - LinkGroupRequest, SubGroupResponse
- `response.go` - Response, ActionResponse, Success(), Error(), Created()

### `internal/api/handlers/`
Todos handlers dependem de `domain.Provider` interface.

**Arquivos:**
- `session.go` - CRUD de sessoes, QR, connect/disconnect/logout
- `message.go` - Envio de mensagens (text, image, video, audio, document, etc)
- `contact.go` - CheckPhone, GetInfo, GetAvatar (rotas /contacts/*)
- `group.go` - CRUD de grupos, participantes, settings, community
- `chat.go` - MarkRead, Presence, Disappearing
- `blocklist.go` - GetBlocklist, Block, Unblock
- `call.go` - RejectCall
- `newsletter.go` - CRUD de canais (not implemented)
- `privacy.go` - Privacy settings (not implemented)
- `profile.go` - Profile e business profile (not implemented)

**Nota:** Rotas de contato usam `/contacts/` (diferente do WuzAPI que usa `/users/`).

### `internal/api/auth/`
**`auth.go`** - Middleware de autenticacao:
- `Global()` - Requer GLOBAL_API_TOKEN
- `Session()` - Aceita GLOBAL_API_TOKEN ou Session Token

### `internal/repository/`
Camada de persistencia com Repository Pattern.

**Arquivos:**
- `repository.go` - Agregador `Repositories` (cria todos os repos)
- `models.go` - `SessionModel` + helpers (`GetJID()`, `NullString()`)
- `session.go` - `SessionRepository` interface + implementacao

### `internal/providers/wameow/`
Implementacao do `domain.Provider` usando whatsmeow.

**Arquivos:**
- `manager.go` - Manager struct, CreateSession, Connect, Disconnect (usa SessionRepository)
- `session.go` - Session struct implementando domain.Session
- `message.go` - SendText, SendImage, SendVideo, etc.
- `group.go` - CreateGroup, GetGroups, SetGroupName, Participants, etc.
- `chat.go` - MarkRead, SendTyping, SetDisappearing, SendPresence
- `user.go` - CheckPhone, GetUserInfo, GetProfilePicture, Blocklist

**Manager depende de:**
- `sqlstore.Container` - Para devices do whatsmeow
- `SessionRepository` - Para persistir sessoes no banco

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
- Handlers dependem de interfaces (`domain.Provider`)
- Logs com zerolog
- Erros wrappados com `fmt.Errorf("context: %w", err)`

### Adding a New Provider
1. Criar pasta em `internal/providers/novoProvider/`
2. Implementar `domain.Provider` interface
3. Criar struct que implementa `domain.Session` interface
4. Injetar `SessionRepository` se precisar persistir sessoes
5. Registrar no `cmd/server/main.go`

### Adding a New Repository
1. Criar model em `internal/repository/models.go`
2. Criar arquivo `internal/repository/novaentidade.go` com interface + impl
3. Adicionar ao agregador `Repositories` em `repository.go`
4. Criar migration SQL em `internal/database/upgrades/`

### Database
- PostgreSQL com pgx driver
- Migrations em `internal/database/upgrades/`
- Colunas em camelCase com aspas duplas
- whatsmeow usa suas proprias tabelas (prefixo `whatsmeow_`)
- Sessoes persistidas via `SessionRepository` (nao apenas em memoria)

### Repository Pattern
- Interfaces em `internal/repository/`
- Models separados dos DTOs (models = banco, DTOs = API)
- Agregador `Repositories` para facilitar injecao de dependencias
- Helpers como `NullString()` para conversao sql.NullString

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
- QR code no terminal usa `qrterminal.GenerateHalfBlock()` (half-blocks Unicode)
- Sessoes persistidas no PostgreSQL via `SessionRepository` + whatsmeow sqlstore
- Sempre usar aspas duplas para colunas camelCase no SQL
- Todos os campos JSON usam PascalCase (compatibilidade WuzAPI)
- Handlers dependem de `domain.Provider` interface, nao de implementacoes concretas
- Provider `wameow` depende de `SessionRepository` para persistencia
- Varios endpoints de Newsletter, Privacy e Profile ainda nao estao implementados
</coding_guidelines>
