APP=friends
SHELL := /bin/bash
DOCKER_COMPOSE_FILE=docker-compose.yaml
# GO_SERVICE=friends1 friends2

.PHONY: help
help: Makefile ## Show this help
	@echo
	@echo "Choose a command run in "$(APP)":"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | sed -e 's/\(\:.*\#\#\)/\:\ /' | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

build-u-srv: ## Build users-srv
	go build -o bin/usrv cmd/users-srv/main.go

build-l-srv: ## Build users-srv
	go build -o bin/lsrv cmd/links-srv/main.go

lint:
	go mod tidy
	golangci-lint run ./...

test: ## Test app
	go test -failfast -count=1 -v ./... -coverpkg=./... -coverprofile=coverpkg.out

migrate-create: ## Migrate DB Create New: use NAME=next_migration_name make migrate-create
	migrate create -ext sql -dir migrations -seq $(NAME)

migrate-up: ## Migrate DB UP
	migrate -path migrations -database postgres://postgres:postgres@localhost:5434/final?sslmode=disable up

migrate-down: ## Migrate DB DOWN
	migrate -path migrations -database postgres://postgres:postgres@localhost:5434/final?sslmode=disable down

start: ## Start all deployed services
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

stop: ## Stop all deployed services
	docker-compose -f $(DOCKER_COMPOSE_FILE) stop

# redeploy: ## Redeploy go services	
# 	docker-compose -f $(DOCKER_COMPOSE_FILE) stop $(GO_SERVICE); \
# 	docker-compose -f $(DOCKER_COMPOSE_FILE) rm -f $(GO_SERVICE); \
# 	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d $(GO_SERVICE); \
# 	docker-compose -f $(DOCKER_COMPOSE_FILE) restart proxy; \