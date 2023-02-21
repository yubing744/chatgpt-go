.PHONY: deps clean build *test run

NAME="chatgpt"

deps:
	go install github.com/joho/godotenv/cmd/godotenv@latest

unit-test: deps
	go test ./pkg/...

int-test: deps
	godotenv -f .env.local go test ./test/...

test: unit-test int-test

run-%: unit-test
	godotenv -f .env.local go run ./examples/$*/cmd.go
