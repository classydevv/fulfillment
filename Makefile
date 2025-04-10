# note: call scripts from /scripts
# Путь до protobuf файлов
PROTO_PATH := $(CURDIR)/api

# Путь до сгенеренных .pb.go файлов
PKG_PROTO_PATH := $(CURDIR)/pkg

# Имя сервиса
SERVICE_NAME := providers

# устанавливаем необходимые плагины
.bin-deps:
	$(info Installing binary dependencies...)

	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# генерация .go файлов с помощью protoc
.protoc-generate:
	mkdir -p $(PKG_PROTO_PATH)
	protoc --proto_path=$(CURDIR) \
	--go_out=$(PKG_PROTO_PATH) --go_opt=paths=source_relative \
	--go-grpc_out=$(PKG_PROTO_PATH) --go-grpc_opt=paths=source_relative \
	$(PROTO_PATH)/$(SERVICE_NAME)/service.proto \
	$(PROTO_PATH)/$(SERVICE_NAME)/messages.proto

# установка нужных зависимостей
.tidy:
	go mod tidy

# тестовые запросы с помощью grpcurl
grpc-provider-save:
	grpcurl -plaintext -d '{"id": "kuper", "name": "Купер"}' \
	localhost:8080 github.com.classydevv.metro_fulfillment.providers.ProvidersService.SaveProvider
grpc-provider-list-all:
	grpcurl -plaintext -d '' \
	localhost:8080 github.com.classydevv.metro_fulfillment.providers.ProvidersService.ListProviders


generate: .bin-deps .protoc-generate .tidy

.PHONY: \
	.bin-deps \
	.protoc-generate \
	.tidy \
	generate