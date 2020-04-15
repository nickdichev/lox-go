build:
	go build -ldflags "-X 'github.com/ziyoung/lox-go/interpreter.envFlag=repl'" -o lox ./cmd/lox/main.go

test:
	go test ./...
