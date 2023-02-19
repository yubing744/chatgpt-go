.PHONY: deps clean build test run

NAME="chatgpt"

deps:
	go install github.com/joho/godotenv/cmd/godotenv@latest

clean:
	rm -rf bin/*

build: clean
	go build -o ./bin/${NAME} ./cmd/cmd.go

test: deps
	godotenv -f .env.local go test ./...

run: deps
	godotenv -f .env.local go run ./cmd/cmd.go
