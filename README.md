# Card Deck

## Implementing REST API over 52-card deck game 

### Requirements
* Go 1.16 or higher
* [Golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) to execute migrations (also check [getting started](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md) to understand how it works)
* [Golint](https://github.com/golang/lint) for linting
* [Gofumpt](https://github.com/mvdan/gofumpt) for smarter formatting

### Getting start:
* Clone the project `git clone https://github.com/dmi-kov/card-deck.git`
* Chdir to project directory
* Run `go version` to ensure that Golang version is 1.16 or higher
* Run postgres db `make start-db`
* Before first run init database `make init-db`
* Run `go mod tidy` to ensure that deps installed
* Run `make run` to run application with `config.hcl`

### Developing
* Run `make run` to start app in current terminal session
* Run `make test && make fmt && make lint` before committing
* Run `make help` to see all available commands
* To extend API functionality add new handler into `internal/api` and init it

### Database
* Run `make start-db` to start Docker container with postgres
* Run `make stop-db` to stop Docker container with postgres
* Run `make init-db` to init database
* Run `make attach-db` to attach Docker container with postgres using psql

### Database migrations
* Migration is executed automatically on app start, or can be executed manually by running `make migrate-up`
* Run `make new-migration MIGRATION_NAME=add_new_table` to create new migration
* Run `make migrate-drop` to drop everything inside DB **(BE CAREFUL!)**
* Run `make migrate-version` to print current version
* Run `make migrate-set-version VERSION=N` to set version N but don't run migration (ignores dirty state)
* Run `make migrate-up` to apply all up migrations
* Run `make migrate-down [STEP=N]` to apply all or N down migrations by passing optional STEP arg

### API description 
* Create deck (default), set `is_shuffled: true` if deck should be shuffled
```
curl --request POST 'http://localhost:8083/v1/deck' \
--header 'Content-Type: application/json' \
--data-raw '{
    "is_shuffled": true
}'
```

* Create deck with provided cards, set `is_shuffled: true` if deck should be shuffled
```
curl --request POST 'http://localhost:8083/v1/deck' \
--header 'Content-Type: application/json' \
--data-raw '{
    "is_shuffled": true,
    "cards": ["10S", "QS", "AD"]
}'
```
* Open Deck by provided deckID
```
curl http://localhost:8083/v1/deck/{deckID}
```
* Draw cards by provided deckID and count of cards
```
curl --request PATCH 'http://localhost:8083/v1/deck/{deckID}/cards' \
--header 'Content-Type: application/json' \
--data-raw '{"count": 10}'
```

### What else?
* Add Dockerfile to build image for running in Docker
* Add e2e tests
* Add validation by [ozzo-validation](github.com/go-ozzo/ozzo-validation/v4)
* Add service layer to move logic from handler??
* Add swagger
