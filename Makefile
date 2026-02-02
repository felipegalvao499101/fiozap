.PHONY: build run test lint swagger clean docker-up docker-down help restart stop

# Variaveis
BINARY_NAME=fiozap
BINARY_PATH=./bin
MAIN_PATH=./cmd/server
DOCS_PATH=./docs

# Build
build:
	@mkdir -p $(BINARY_PATH)
	go build -o $(BINARY_PATH)/$(BINARY_NAME) $(MAIN_PATH)

# Run
run:
	go run $(MAIN_PATH)

# Test
test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

# Swagger
swagger:
	swag init -g cmd/server/main.go -o $(DOCS_PATH)

swagger-fmt:
	swag fmt

# Clean
clean:
	rm -rf $(BINARY_PATH)
	rm -f coverage.out coverage.html

# Docker
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-build:
	docker build -t $(BINARY_NAME) .

# Database
db-up:
	docker compose up -d postgres

db-down:
	docker compose down postgres

# Development
dev: swagger run

# Stop server
stop:
	@echo "Parando servidor..."
	@-pkill -f "$(BINARY_NAME)" 2>/dev/null || true
	@-pkill -f "go run $(MAIN_PATH)" 2>/dev/null || true
	@-lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@echo "Servidor parado"

# Restart server
restart: stop
	@sleep 1
	@echo "Iniciando servidor..."
	@$(MAKE) run

# Install tools
tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Help
help:
	@echo "Comandos disponiveis:"
	@echo "  make build         - Compila o binario"
	@echo "  make run           - Executa a aplicacao"
	@echo "  make test          - Roda os testes"
	@echo "  make test-verbose  - Roda os testes com output detalhado"
	@echo "  make test-coverage - Gera relatorio de cobertura"
	@echo "  make lint          - Roda o linter"
	@echo "  make lint-fix      - Roda o linter e corrige automaticamente"
	@echo "  make swagger       - Gera documentacao Swagger"
	@echo "  make swagger-fmt   - Formata anotacoes Swagger"
	@echo "  make clean         - Remove arquivos gerados"
	@echo "  make docker-up     - Sobe containers Docker"
	@echo "  make docker-down   - Para containers Docker"
	@echo "  make docker-logs   - Mostra logs dos containers"
	@echo "  make docker-build  - Build da imagem Docker"
	@echo "  make db-up         - Sobe apenas o PostgreSQL"
	@echo "  make dev           - Gera swagger e executa"
	@echo "  make tools         - Instala ferramentas de desenvolvimento"
	@echo "  make stop          - Para o servidor"
	@echo "  make restart       - Reinicia o servidor"
	@echo "  make help          - Mostra esta ajuda"
