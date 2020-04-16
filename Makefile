PKG_NAME = github.com/ziyoung/lox-go
GO_LDFLAGS = -X ${PKG_NAME}/interpreter.evalEnv=repl

build:
	go build -ldflags "${GO_LDFLAGS}" -o lox ./cmd/lox/main.go

test:
	go test ./...
