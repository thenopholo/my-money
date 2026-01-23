# My Money - Makefile
# Carrega variáveis do .env
include .env
export

# ==================== DATABASE ====================

.PHONY: db-up
db-up: ## Inicia o PostgreSQL
	docker-compose up -d db

.PHONY: db-down
db-down: ## Para o PostgreSQL
	docker-compose down

.PHONY: db-logs
db-logs: ## Mostra logs do PostgreSQL
	docker-compose logs -f db

.PHONY: db-shell
db-shell: ## Abre psql no container
	docker-compose exec db psql -U $(DB_USER) -d $(DB_NAME)

# ==================== MIGRATIONS ====================

.PHONY: migrate
migrate: ## Roda todas as migrations pendentes
	tern migrate \
		--migrations ./internal/migrations \
		--config ./internal/migrations/tern.conf

.PHONY: migrate-down
migrate-down: ## Reverte a última migration
	tern migrate \
		--destination -1 \
		--migrations ./internal/migrations \
		--config ./internal/migrations/tern.conf

.PHONY: migrate-status
migrate-status: ## Mostra status das migrations
	tern status \
		--migrations ./internal/migrations \
		--config ./internal/migrations/tern.conf

.PHONY: migrate-new
migrate-new: ## Cria uma nova migration (uso: make migrate-new name=create_users)
	@if [ -z "$(name)" ]; then \
		echo "Erro: especifique o nome da migration com name=xxx"; \
		exit 1; \
	fi
	@count=$$(ls -1 ./internal/migrations/*.sql 2>/dev/null | wc -l); \
	num=$$(printf "%03d" $$((count + 1))); \
	touch "./internal/migrations/$${num}_$(name).sql"; \
	echo "Criada: ./internal/migrations/$${num}_$(name).sql"

# ==================== SQLC ====================

.PHONY: sqlc
sqlc: ## Gera código Go a partir das queries SQL
	sqlc generate

.PHONY: sqlc-check
sqlc-check: ## Verifica se o código gerado está atualizado
	sqlc diff

# ==================== DEVELOPMENT ====================

.PHONY: run
run: ## Roda a aplicação
	go run ./cmd/api

.PHONY: dev
dev: ## Roda com hot-reload (air)
	air

.PHONY: build
build: ## Compila a aplicação
	go build -o ./bin/api ./cmd/api

.PHONY: test
test: ## Roda os testes
	go test -v ./...

.PHONY: test-cover
test-cover: ## Roda testes com coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Abra coverage.html no navegador"

.PHONY: lint
lint: ## Roda o linter
	golangci-lint run

# ==================== SETUP ====================

.PHONY: setup
setup: ## Setup inicial do projeto
	@echo "Instalando dependências..."
	go mod download
	@echo "Instalando ferramentas..."
	go install github.com/air-verse/air@latest
	go install github.com/jackc/tern/v2@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Copiando .env.example para .env..."
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@echo "Setup completo!"

# ==================== DOCKER ====================

.PHONY: docker-build
docker-build: ## Build da imagem Docker
	docker build -t my-money:latest .

.PHONY: docker-run
docker-run: ## Roda a aplicação no Docker
	docker-compose up -d

.PHONY: docker-stop
docker-stop: ## Para todos os containers
	docker-compose down

# ==================== HELP ====================

.PHONY: help
help: ## Mostra esta mensagem de ajuda
	@echo "My Money - Comandos disponíveis:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help