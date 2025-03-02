.PHONY: all build run clean test

all: clean build run

build:
	go build -o bin/weatherapp cmd/server/main.go

run:
	./bin/weatherapp

test:
	go test ./...

clean:
	go mod tidy
	rm -rf bin/*

dev:
	go run cmd/server/main.go



