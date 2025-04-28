include .env
export

# Путь до protobuf файлов
PROTO_PATH := $(CURDIR)/api

# Путь до сгенеренных .pb.go файлов
PKG_PROTO_PATH := $(CURDIR)/pkg

# Имя сервиса
SERVICE_NAME := providers

# Стэк сервисов без тестов
BASE_STACK = docker compose -f docker-compose.yml


# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-build: ### Build docker compose
	$(BASE_STACK) build

.PHONY: compose-build

compose-up: ### Up docker compose
	$(BASE_STACK) up -d
.PHONY: compose-up

compose-down: ### Down docker compose
	$(BASE_STACK) down --remove-orphans
.PHONY: compose-down

compose-reload: compose-down compose-build compose-up ### Quick reload for development use 
.PHONY: compose-reload

linter-golangci: ### Check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

swag:	
	swag init -g internal/providers/controller/http/router.go
.PHONY: swag

.protoc-deps: ### Install necessary protoc plugins
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: .bin-deps

# генерация .go файлов с помощью protoc
.protoc-generate:
	mkdir -p $(PKG_PROTO_PATH)
	protoc --proto_path=$(CURDIR) \
	--go_out=$(PKG_PROTO_PATH) --go_opt=paths=source_relative \
	--go-grpc_out=$(PKG_PROTO_PATH) --go-grpc_opt=paths=source_relative \
	$(PROTO_PATH)/$(SERVICE_NAME)/service.proto \
	$(PROTO_PATH)/$(SERVICE_NAME)/messages.proto
.PHONY: .protoc-generate

.bin-deps: ### Install necessary tools
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
.PHONY: .bin-deps

.tidy:
	go mod tidy && go mod verify
.PHONY: .tidy

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### migration down
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' down
.PHONY: migrate-down

generate: .protoc-deps .protoc-generate .tidy
.PHONY: generate

run: .bin-deps .tidy
	CGO_ENABLED=0 go run -tags migrate ./cmd/providers
.PHONY: run

# тестовые запросы с помощью grpcurl
grpc-provider-save:
	grpcurl -plaintext -d '{"id": "kuper", "name": "Купер"}' \
	localhost:8082 github.com.classydevv.fulfillment.providers.ProvidersService.SaveProvider
grpc-provider-list-all:
	grpcurl -plaintext -d '' \
	localhost:8082 github.com.classydevv.fulfillment.providers.ProvidersService.ListProviders