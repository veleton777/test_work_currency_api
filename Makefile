ifneq ("$(wildcard .env)","")
  include .env
endif

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_BUILD:=$(CURDIR)/build
PROJECT_PATH := $(shell pwd)
GO = $(shell which go)

LOCAL_MIGRATION_DIR=./migrations
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE) user=$(PG_USER) password=$(PG_PASSWORD)"

first-init: install-go-deps docker-postgres-up wait-db local-migration-up

init-env:
	cp $(CURDIR)/.env.example $(CURDIR)/.env

docker-postgres-up:
	docker-compose up postgres -d

run:
	docker-compose up app

install-go-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@v2.43.2
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@v1.16.3

local-migration-status:
	./bin/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	./bin/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	./bin/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

test-migration-up:
	./bin/goose -dir ${LOCAL_MIGRATION_DIR} postgres "host=localhost port=$(PG_PORT) dbname=test user=$(PG_USER) password=$(PG_PASSWORD)" up -v

wait-db:
	sleep 5

build-app:
	cd cmd/app && go build -o ${LOCAL_BUILD}/app

local-start:
	go run cmd/app/main.go

prepare-tests: ## подготовить окружение к тестированию
	@printf $(COLOR) "Generate ..."
	$(GO) generate ./...

test: prepare-tests ## Запустить тесты
	@printf $(COLOR) "Run tests ..."
	$(GO) test -race -cover -short -v -tags mock -count=2 \
				-coverprofile profile.cov.tmp -p 1 \
				./...
	cat profile.cov.tmp | grep -Ev "_gen.go|mock_|mocks" > profile.cov
	$(MAKE) cover

local-test: prepare-tests create-test-db test-migration-up ## Запустить все тесты, включая те, что не используют моки
	@printf $(COLOR) "Run tests ..."
	export PG_DATABASE=test PG_HOST=localhost PG_PORT=${PG_PORT} PG_USER=${PG_USER} PG_PASSWORD=${PG_PASSWORD} && \
	$(GO) test -race -cover -short -v -tags mock,integration -count=1 \
				-coverprofile profile.cov.tmp -p 1 \
				./...
	cat profile.cov.tmp | grep -Ev "_gen.go|mock_|mocks" > profile.cov
	$(MAKE) cover

cover: ## Посчитать coverage проекта
	@printf $(COLOR) "Code coverage ..."
	$(GO) tool cover -func profile.cov

create-test-db:
	docker-compose exec postgres bash -c "echo \"SELECT 'CREATE DATABASE test' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'test')\gexec\" | psql \"postgresql://${PG_USER}:${PG_PASSWORD}@localhost/postgres\""

swag:
	./bin/swag init -d internal/server,internal -g router.go -o ./api/swagger --ot yaml

lint:
	./bin/golangci-lint run
