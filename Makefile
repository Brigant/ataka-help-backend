# Define targets

all: tests lintcheck build

dev-run: dbrun run

dbrun:
	docker-compose up -d db

run:
	go run ./cmd/main.go

build:
	docker-compose up -d --build

stop:
	docker-compose down

lintckeck-all:
	golangci-lint --enable-all --no-config

lintcheck:
	golangci-lint

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html	&& xdg-open ./coverage.html