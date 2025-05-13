# Step 1: Modules caching
FROM golang:1.24.3-alpine as modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

# Step 2: Builder
FROM golang:1.24.3-alpine as builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/providers/main.go 
    # go build -o /bin/app ./cmd/providers/main.go 

# Step 3: Final
FROM scratch

# COPY --from=builder /app/configs/providers /config
COPY --from=builder /bin/app /app
COPY --from=builder /app/migrations /migrations

# EXPOSE 8080

ENTRYPOINT ["/app"]