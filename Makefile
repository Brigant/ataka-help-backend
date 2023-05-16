# Define targets

all: tests lintcheck build

dev-run: dbrun run

dbrun:
	docker-compose up -d db

migration:
	docker-compose up -d db migrate-up

migration-down:
	docker-compose -f migrate-down.yml up -d

run:
	go run ./cmd/main.go

build:
	docker-compose up -d --build

stop:
	docker-compose down

lintcheck:
	golangci-lint run

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html	&& xdg-open ./coverage.html