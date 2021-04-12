all: server

server:
	go build -o excel2config cmd/main.go
