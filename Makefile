PKG_MAIN=./cmd/main.go
BINARY=card-deck
RUN_ARGS=-config=config.hcl
DB_USER?=games
DB_PASSWORD?=games
DB_PORT?=5437
DB_CONTAINER_NAME?=games_postgres
DB_NAME?=games
MIGRATIONS_DIR?=./db/migrations
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@127.0.0.1:$(DB_PORT)/$(DB_NAME)?sslmode=disable&search_path=public
GOMIGRATE_CMD=migrate -source=file://$(MIGRATIONS_DIR) -database='$(DB_URL)'

.PHONY:run
run: ## Run application with default config
	@go run $(PKG_MAIN) $(RUN_ARGS)

.PHONY:build
build: ## Build application
	@echo Building app... && \
	go build -o $(BINARY) $(PKG_MAIN)

.PHONY:clean
clean: ## Remove binary and run go clean
	@echo Cleaning... && \
	rm -rf $(BINARY) \
	go clean

.PHONY:fmt
fmt: ## Run formatting by gofumpt
	@echo Formatting... && \
	gofumpt -extra -l -w .

.PHONY:lint
lint: ## Run golint
	@echo Running golint... && \
	golint ./...

.PHONY:test
test: ## Run tests
	@echo Running tests... && \
	go test ./...

.PHONY:start-db
start-db: ## Start docker container with postgres
	@echo Starting docker container"$(DB_CONTAINER_NAME)"... && \
	docker-compose up -d db

.PHONY:stop-db
stop-db: ## Stop docker container with postgres
	@echo Stopping docker container "$(DB_CONTAINER_NAME)"... && \
	docker stop $(DB_CONTAINER_NAME)

.PHONY:init-db
init-db: ## Create new database
	@echo Creating database "$(DB_NAME)"... && \
	docker exec $(DB_CONTAINER_NAME) psql -U $(DB_USER) -c 'CREATE DATABASE $(DB_NAME)'

.PHONY:attach-db
attach-db: ## Attach to DB running in Docker
	@echo Attaching database "$(DB_NAME)"... && \
	docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME)

.PHONY:new-migration
new-migration: ## Create new migration with specified name by passed MIGRATION_NAME param
	@if [ -z $(MIGRATION_NAME) ]; then echo "usage: 'make new-migration MIGRATION_NAME=migration-name'" && exit 1; fi; \
	$ migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(MIGRATION_NAME) && \
	echo "Migration '$(MIGRATION_NAME)' has been created"

.PHONY:migrate-version
migrate-version: ## Check current version of migration
	@echo Version is ... && \
	$(GOMIGRATE_CMD) version

.PHONY:migrate-drop
migrate-drop: ## Drop everything inside DB
	@echo Dropping ... && \
	$(GOMIGRATE_CMD) drop

.PHONY:migrate-set-version
migrate-set-version: ## Set version N but don't run migration (ignores dirty state)
	@if [ -z $(VERSION) ]; then echo "usage: 'make migrate-set-version VERSION=N'" && exit 1; fi; \
	$(GOMIGRATE_CMD) force $(VERSION);

.PHONY:migrate-up
migrate-up: ## Apply all up migrations
	$(GOMIGRATE_CMD) up

.PHONY:migrate-down
migrate-down: ## Apply all or N down migrations by passing STEP arg
	$(GOMIGRATE_CMD) down $(STEP)

.PHONY: help
help: ## Show all commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'