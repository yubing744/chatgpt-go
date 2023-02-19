.PHONY: deps clean build test run

deps:
	go install github.com/joho/godotenv/cmd/godotenv@latest

clean:
	rm -rf build/*

build: clean
	CGO_ENABLED=0 go build -o ./build/bbgo ./cmd/cmd.go

test: deps
	godotenv -f .env.local go test ./...

run:
	godotenv -f .env.local go run ./cmd/cmd.go
