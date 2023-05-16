# Define targets

all: tests lintcheck build

run:
	docker-compose up -d

rebuild:
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