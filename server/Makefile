run:
	@echo "Run Server Script"
	go run cmd/api/main.go --config="./config/development"

dev:
	@echo "Starting docker development enviroment"
	docker-compose -f docker-compose.dev.yml up --build

# Migrations

driver = "postgres"
dsn = "host=localhost port=5432 user=postgres password=postgres dbname=gopix_dev sslmode=disable"
migrationPath = "./migrations"
type = "sql"

.PHONY: goose
goose_status:
	@GOOSE_DRIVER=$(driver) GOOSE_DBSTRING=$(dsn) goose -dir=$(migrationPath) status

goose_up:
	@GOOSE_DRIVER=$(driver) GOOSE_DBSTRING=$(dsn) goose -dir=$(migrationPath) up

goose_reset:
	@GOOSE_DRIVER=$(driver) GOOSE_DBSTRING=$(dsn) goose -dir=$(migrationPath) reset

goose_create:
	@GOOSE_DRIVER=$(driver) GOOSE_DBSTRING=$(dsn) goose -dir=$(migrationPath) create $(name) $(type)

.PHONY: pg_migrate
pg_status: goose_status

pg_migrate: goose_up

pg_reset: goose_reset