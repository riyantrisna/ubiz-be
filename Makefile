generate:
	go generate ./...

run: generate
	go run .

.PHONY: generate wire run