# FioZap

API REST multi-session para WhatsApp usando whatsmeow, com integracao Chatwoot.

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
│   │   ├── handlers/    → HTTP handlers (session, message)
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

## Database Schema

### Tabela `sessions`

| Coluna      | Tipo         | Descricao                     |
|-------------|--------------|-------------------------------|
| `id`        | VARCHAR(255) | Primary key (UUID)            |
| `name`      | VARCHAR(255) | Identificador unico da sessao |
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

**Ao criar sessao**, o token e retornado na resposta (guarde-o!):
```json
{"success":true,"data":{"name":"session1","token":"abc123...","connected":false}}
```

## API Endpoints

### Session (requer GLOBAL_API_TOKEN)
- `POST /sessions` - Criar sessao `{"name": "session1"}` → retorna token
- `GET /sessions` - Listar sessoes

### Session (aceita GLOBAL ou Session Token)
- `GET /sessions/:name` - Status da sessao
- `GET /sessions/:name/qr` - QR code (base64 ou image)
- `POST /sessions/:name/connect` - Conectar
- `POST /sessions/:name/disconnect` - Desconectar
- `DELETE /sessions/:name` - Remover sessao

### Messages (aceita GLOBAL ou Session Token)
- `POST /sessions/:name/messages/text` - Enviar texto
- `POST /sessions/:name/messages/image` - Enviar imagem
- `POST /sessions/:name/messages/document` - Enviar documento
- `POST /sessions/:name/messages/audio` - Enviar audio
- `POST /sessions/:name/messages/location` - Enviar localizacao

### Users (aceita GLOBAL ou Session Token)
- `POST /sessions/:name/users/check` - Verificar WhatsApp
- `GET /sessions/:name/users/:phone` - Info do usuario
- `GET /sessions/:name/users/:phone/avatar` - Foto de perfil

## Modulos e Responsabilidades

### `cmd/server/`
**Responsabilidade:** Ponto de entrada da aplicacao.
- Carregar configuracoes
- Inicializar dependencias (database, logger, manager)
- Iniciar servidor HTTP
- Graceful shutdown

**NAO deve:**
- Conter logica de negocio
- Importar handlers diretamente
- Manipular requests/responses

---

### `internal/config/`
**Responsabilidade:** Gerenciamento de configuracoes.
- Carregar variaveis de ambiente
- Validar configuracoes obrigatorias
- Prover struct `Config` imutavel

**NAO deve:**
- Acessar banco de dados
- Fazer logging complexo
- Depender de outros modulos internos

---

### `internal/database/`
**Responsabilidade:** Conexao e migrations do banco de dados.
- Gerenciar conexao PostgreSQL
- Executar migrations SQL
- Prover container sqlstore para whatsmeow

**NAO deve:**
- Conter queries de negocio (usar repositories)
- Conhecer entidades de dominio
- Fazer cache

**Arquivos:**
- `postgres.go` - Conexao e migration runner
- `upgrades/*.sql` - Arquivos de migration

---

### `internal/logger/`
**Responsabilidade:** Configuracao de logging.
- Criar logger zerolog configurado
- Adapter para whatsmeow (WALogger)

**NAO deve:**
- Conter logica de negocio
- Persistir logs (apenas stdout)
- Depender de outros modulos

---

### `internal/api/dto/`
**Responsabilidade:** Data Transfer Objects para API.
- Structs de request (input)
- Structs de response (output)
- Funcoes helper (SuccessResponse, ErrorResponse)

**NAO deve:**
- Conter validacao complexa
- Acessar banco de dados
- Conhecer whatsmeow

**Convencoes:**
- Usar tags `json:"fieldName"`
- Request structs terminam com `Request`
- Response structs terminam com `Response`

---

### `internal/api/handlers/`
**Responsabilidade:** Handlers HTTP (controllers).
- Receber requests HTTP
- Validar input basico
- Chamar metodos do Manager
- Retornar responses padronizadas

**NAO deve:**
- Conter logica de negocio complexa
- Acessar banco diretamente
- Conhecer detalhes do whatsmeow

**Arquivos:**
- `session.go` - CRUD de sessoes, QR, connect/disconnect
- `message.go` - Envio de mensagens, check phone, user info

**Convencoes:**
- Handlers sao `func(w http.ResponseWriter, r *http.Request)`
- Usar `chi.URLParam(r, "name")` para extrair params
- Usar `dto.Success()`, `dto.Error()`, `dto.Created()` para responses

---

### `internal/api/router/`
**Responsabilidade:** Configuracao de rotas HTTP.
- Criar router Chi
- Registrar middlewares (Recoverer, RequestID, Logger)
- Mapear rotas para handlers

**NAO deve:**
- Conter logica de handler
- Acessar banco de dados
- Instanciar servicos

---

### `internal/zap/`
**Responsabilidade:** Gerenciamento de sessoes WhatsApp.
- CRUD de sessoes em memoria
- Conectar/Desconectar do WhatsApp
- Enviar mensagens e medias
- Tratar eventos do WhatsApp

**NAO deve:**
- Conhecer HTTP/Echo
- Acessar DTOs da API
- Fazer logging excessivo

**Arquivos:**
- `session.go` - Struct Session com estado da conexao
- `manager.go` - Gerenciador de sessoes (CRUD, Connect, Disconnect)
- `send.go` - Metodos de envio (SendText, SendImage, etc)
- `events.go` - Handler de eventos (Connected, Disconnected, Message)

**Convencoes:**
- Manager e thread-safe (sync.RWMutex)
- Session armazena estado em memoria
- Handlers chamam metodos do Manager diretamente

---

### `internal/webhook/`
**Responsabilidade:** Disparar webhooks para URLs externas (futuro).
- Enviar eventos via HTTP POST
- Retry em caso de falha
- Serializar payload JSON

**NAO deve:**
- Conhecer whatsmeow diretamente
- Persistir webhooks
- Bloquear execucao principal

---

### `internal/chatwoot/`
**Responsabilidade:** Integracao com Chatwoot (futuro).
- Sincronizar mensagens WhatsApp <-> Chatwoot
- Gerenciar inbox e contatos
- Traduzir eventos entre sistemas

**NAO deve:**
- Conhecer HTTP handlers
- Acessar whatsmeow diretamente (usar zap/)
- Gerenciar sessoes

---

## Fluxo de Dependencias

```
cmd/server/
    ├── config/
    ├── database/
    ├── logger/
    └── api/router/
            └── api/handlers/
                    ├── api/dto/
                    └── zap/ (Manager)

zap/
    └── database/ (sqlstore container)

chatwoot/ (futuro)
    ├── zap/ (eventos)
    └── webhook/ (notificacoes)
```

**Regra:** Dependencias fluem de cima para baixo. Modulos inferiores NAO importam superiores.

---

## Development Patterns

### Coding Style
- Go 1.22+
- Usar `internal/` para codigo privado
- DTOs em `internal/api/dto/`
- Handlers retornam `dto.Response` padronizado
- Logs com zerolog
- Erros wrappados com `fmt.Errorf("context: %w", err)`

### Database
- PostgreSQL com pgx driver
- Migrations em `internal/database/upgrades/`
- Colunas em camelCase com aspas duplas
- whatsmeow usa suas proprias tabelas (prefixo `whatsmeow_`)

### WhatsApp (zap/)
- `Manager` gerencia multiplas sessoes e expoe todos os metodos
- `Session` contem estado da conexao (Client, Device, QRCode)
- Eventos tratados em `events.go` via funcao package-level

## Environment Variables

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
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
