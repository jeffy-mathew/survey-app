all: prepare test build run

prepare:
	go mod download

test:
	go test -race -cover ./...

build:
	go build -o survey-platform cmd/main.go

run:
	./survey-platform