GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

vuln:
	govulncheck -show verbose ./...

cli:
	go build -mod vendor -o bin/broadcast cmd/broadcast/main.go

hello:
	go run cmd/broadcast/main.go -broadcaster log:// -broadcaster null:// -body "hello world"
