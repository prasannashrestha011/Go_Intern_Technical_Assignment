include .env
export
migrations_dir=migrations

.PHONY: migrate-up migrate-down migrate-status migrate-create

migrate-up:
	@goose -dir $(migrations_dir) postgres "$(DB_URL)" up

migrate-down:
	@goose -dir $(migrations_dir) postgres "$(DB_URL)" down

migrate-status:
	@goose -dir $(migrations_dir) postgres "$(DB_URL)" status

migrate-create:
	@goose -dir $(migrations_dir) create $(name) sql
