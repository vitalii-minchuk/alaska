build:
	@go build -o bin/alaska cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/alaska