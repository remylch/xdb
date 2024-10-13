run:
	go run .

build:
	@go build -o bin/node main.go

test:
	HASH_KEY=your-32-byte-secret-key-here!!!! go test ./... -count=1 -timeout 30s

build-local-graph:
	docker compose build

start-local-graph:
	docker compose up
