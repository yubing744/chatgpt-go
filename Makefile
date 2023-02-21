.PHONY: deps clean build *test run

NAME="chatgpt"

deps:
	go install github.com/joho/godotenv/cmd/godotenv@latest

unit-test:
	go test ./pkg/...

int-test: deps
	godotenv -f .env.local go test ./test/...

test: unit-test int-test

run-%:
	go run ./examples/$*/cmd.go
