run:
	go run ./cmd/node/ -f ./config.toml

build:
	go build -o bin/node ./cmd/node

tests:
	go test ./... -count=1 -timeout 30s

build-local-graph:
	docker compose build

start-local-graph:
	docker compose up
