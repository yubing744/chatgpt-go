.PHONY: deps clean build *test run

NAME="chatgpt"

deps:
	go install github.com/joho/godotenv/cmd/godotenv@latest

clean:
	rm -rf bin/*

unit-test: deps
	godotenv -f .env.local go test ./pkg/...

int-test: deps
	godotenv -f .env.local go test ./test/...

run-%: deps
	godotenv -f .env.local go run ./examples/$*/cmd.go
