generate:
	go generate ./...

wire: 
	wire

run: wire generate
	go run .

.PHONY: generate wire run