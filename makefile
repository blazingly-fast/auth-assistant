build:
	@go build -o bin/network

run: build
	@./bin/network

test:
	@go test -v -cover ./...
