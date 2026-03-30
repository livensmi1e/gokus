APP_NAME=gokus

.PHONY: all build test run clean tidy fmt vet

all: test build

build:
	go build -o bin/$(APP_NAME) ./cmd/gokus

run: build
	./bin/$(APP_NAME)

test:
	go test ./internal/... -v

test-short:
	go test ./internal/... -v -short

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin