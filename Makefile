include .env
export

.PHONY: migrate
migrate:
	tern migrate --migrations ./internal/migrations --config ./internal/migrations/tern.conf

.PHONY: migrate-down
migrate-down:
	tern migrate --destination -1 --migrations ./internal/migrations --config ./internal/migrations/tern.conf