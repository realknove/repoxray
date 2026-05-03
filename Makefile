.PHONY: fmt lint test run build

fmt:
	go fmt ./...

lint:
	go vet ./...

test:
	go test ./...

run:
	go run ./cmd/repoxray scan .

build:
	go build -o repoxray ./cmd/repoxray
