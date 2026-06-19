.PHONY: build run test lint clean up down

APP_NAME = mentorship-backend
BUILD_DIR = ./bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)

up:
	docker compose up -d

down:
	docker compose down

restart: down up