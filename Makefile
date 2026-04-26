APP_NAME ?= hcm-go
BINARY ?= bin/$(APP_NAME)

.PHONY: run build test fmt vet tidy docker-build compose-up compose-down

run:
	go run ./cmd/api

build:
	mkdir -p bin
	go build -o $(BINARY) ./cmd/api

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

docker-build:
	docker build -t $(APP_NAME):local .

compose-up:
	docker compose up --build -d

compose-down:
	docker compose down
