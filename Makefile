BINARY_NAME=octopus-server

.PHONY: build
build:
	go build -o ${BINARY_NAME} ./cmd/server

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: dev
dev:
	air

.PHONY: clean
clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -rf ./tmp

.PHONY: test
test:
	go test ./...

.PHONY: deps
deps:
	go mod tidy
	go mod download
