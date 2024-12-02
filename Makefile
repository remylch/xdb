run:
	go run .

build:
	@go build -o bin/node *.go

tests:
	go test ./... -count=1 -timeout 30s

build-local-graph:
	docker compose build

start-local-graph:
	docker compose up
