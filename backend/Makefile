DB_NAME=practice_IstraNet
DB_PASSWORD=1
DB_PORT=5436
DB_HOST=localhost
APP_PORT=8000
MIGRATIONS_PATH=./migrations
DATABASE_URL=postgres://postgres:$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: db-start db-stop db-clean migrate-up run build clean

db-start:
	docker run --name=todo-db \
	 -e POSTGRES_PASSWORD='$(DB_PASSWORD)' \
	 -e POSTGRES_DB='$(DB_NAME)' \
	 -p $(DB_PORT):5432 -d --rm postgres

db-stop:
	docker rm -f todo-db || true

db-clean: db-stop
	@echo "Database container cleaned"

db-restart: db-stop db-start

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

run:
	go run cmd/main.go

dev: db-start migrate-up run


deps:
	go mod download

tidy:
	go mod tidy

db-check:
	@pg_isready -h $(DB_HOST) -p $(DB_PORT) -d $(DB_NAME) -U postgres || echo "Database is not ready"

db-status:
	@docker ps -f name=todo-db

reset: db-clean db-start migrate-up
	@echo "Reset complete!"
