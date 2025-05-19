include .env
export

# Путь до protobuf файлов
PROTO_PATH := $(CURDIR)/api

# Путь до buf yaml файла
BUF_GEN_FILE := $(CURDIR)/configs/buf/buf.gen.yaml

# Путь до сгенеренных .pb.go файлов
PKG_PROTO_PATH := $(CURDIR)/pkg

# Путь до завендореных protobuf файлов
VENDOR_PROTO_PATH := $(CURDIR)/vendor.protobuf

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

bin-deps: ### Install necessary tools
	go install tool
.PHONY: bin-deps

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

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### migration down
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' down
.PHONY: migrate-down

linter-golangci: ### Check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

format: ### Format code
	gofumpt -l -w .
.PHONY: format

tidy:
	go mod tidy && go mod verify
.PHONY: tidy

security: ### Security check
	govulncheck ./...
	gosec ./...
.PHONY: security

swag:	
	swag init -g internal/providers/controller/http/router.go
.PHONY: swag

vendor-proto: .vendor-proto-reset .vendor-proto-googleapis .vendor-proto-google-protobuf .vendor-proto-protovalidate .vendor-proto-protoc-gen-openapiv2 .vendor-proto-tidy

.vendor-proto-reset:
	rm -rf $(VENDOR_PROTO_PATH)
	mkdir -p $(VENDOR_PROTO_PATH)
.PHONY: .vendor-proto-reset

.vendor-proto-tidy:
	find $(VENDOR_PROTO_PATH) -type f ! -name "*.proto" -delete
	find $(VENDOR_PROTO_PATH) -empty -type d -delete
.PHONY: .vendor-proto-tidy

.vendor-proto-google-protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf $(VENDOR_PROTO_PATH)/protobuf &&\
	cd $(VENDOR_PROTO_PATH)/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/google
	mv $(VENDOR_PROTO_PATH)/protobuf/src/google/protobuf $(VENDOR_PROTO_PATH)/google
	rm -rf $(VENDOR_PROTO_PATH)/protobuf
.PHONY: .vendor-proto-google-protobuf

.vendor-proto-protovalidate:
	git clone -b main --single-branch --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protovalidate $(VENDOR_PROTO_PATH)/protovalidate && \
	cd $(VENDOR_PROTO_PATH)/protovalidate
	git checkout
	mv $(VENDOR_PROTO_PATH)/protovalidate/proto/protovalidate/buf $(VENDOR_PROTO_PATH)
	rm -rf $(VENDOR_PROTO_PATH)/protovalidate
.PHONY: .vendor-proto-protovalidate

.vendor-proto-googleapis:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis $(VENDOR_PROTO_PATH)/googleapis &&\
	cd $(VENDOR_PROTO_PATH)/googleapis &&\
	git checkout
	mv $(VENDOR_PROTO_PATH)/googleapis/google $(VENDOR_PROTO_PATH)
	rm -rf $(VENDOR_PROTO_PATH)/googleapis
.PHONY: .vendor-proto-googleapis

.vendor-proto-protoc-gen-openapiv2:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway $(VENDOR_PROTO_PATH)/grpc-gateway && \
 	cd $(VENDOR_PROTO_PATH)/grpc-gateway && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	mv $(VENDOR_PROTO_PATH)/grpc-gateway/protoc-gen-openapiv2/options $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	rm -rf $(VENDOR_PROTO_PATH)/grpc-gateway
.PHONY: .vendor-proto-protoc-gen-openapiv2

protoc-generate: ### Generate code from .proto using protoc
	mkdir -p $(PKG_PROTO_PATH)
	protoc -I $(VENDOR_PROTO_PATH) --proto_path=$(CURDIR) \
	--go_out=$(PKG_PROTO_PATH) --go_opt paths=source_relative \
	--go-grpc_out=$(PKG_PROTO_PATH) --go-grpc_opt paths=source_relative \
	--grpc-gateway_out=$(PKG_PROTO_PATH) --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	$(PROTO_PATH)/$(SERVICE_NAME)/service.proto \
	$(PROTO_PATH)/$(SERVICE_NAME)/messages.proto

	protoc -I $(VENDOR_PROTO_PATH) --proto_path=$(CURDIR) \
	--openapiv2_out=. --openapiv2_opt logtostderr=true --openapiv2_opt generate_unbound_methods=true \
	$(PROTO_PATH)/$(SERVICE_NAME)/service.proto

buf-generate: ### Generate code from .proto using buf
	buf --template=$(BUF_GEN_FILE) generate
.PHONY: buf-generate

proto-generate: protoc-generate tidy
.PHONY: proto-generate

mock:
	go generate -run mockgen -n ./internal/...
.PHONY: mock

run: bin-deps tidy
	CGO_ENABLED=0 go run -tags migrate ./cmd/providers
.PHONY: run

test:
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test

pre-commit: swag proto-generate mock format linter-golangci test
.PHONY: pre-commit

# тестовые запросы с помощью grpcurl
grpc-provider-create:
	grpcurl -plaintext -d '{"provider_id": "kuper", "name": "Купер"}' \
	localhost:8082 github.com.classydevv.fulfillment.providers.v1.ProvidersService.ProviderCreate
grpc-provider-list-all:
	grpcurl -plaintext -d '' \
	localhost:8082 github.com.classydevv.fulfillment.providers.v1.ProvidersService.ProviderListAll
grpc-provider-update:
	grpcurl -plaintext -d '{"provider_id": "kuper"}' \
	localhost:8082 github.com.classydevv.fulfillment.providers.v1.ProvidersService.ProviderUpdate
grpc-provider-delete:
	grpcurl -plaintext -d '{"provider_id": "kuper"}' \
	localhost:8082 github.com.classydevv.fulfillment.providers.v1.ProvidersService.ProviderDelete