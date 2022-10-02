generate:
	go generate ./...

run: generate
	go run .

.PHONY: generate run