run:
	go run .

build:
	@go build -o bin/gateway main.go

test:
	@go test -v ./... -count=1

build-local-graph:
	docker compose build

start-local-graph:
	docker compose up


