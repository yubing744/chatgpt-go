.PHONY: deps clean build test run

NAME="chatgpt"

deps:
	go install github.com/joho/godotenv/cmd/godotenv@latest

clean:
	rm -rf bin/*

test: deps
	godotenv -f .env.local go test ./...

run-%: deps
	godotenv -f .env.local go run ./examples/$*/cmd.go
