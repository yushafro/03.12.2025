all:
	go run cmd/status/*.go

build:
	go build -o status cmd/status/*.go