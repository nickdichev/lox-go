build:
	go build -o lox ./cmd/lox/main.go

test:
	go test ./...
