build:
	@env GOOS=linux GOARCH=amd64 go build -o dnsproxy.bin ./cmd

run:
	@./build/dnsproxy

test:
	go test -v ./...
