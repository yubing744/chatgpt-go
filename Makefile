.PHONY: clean build test run


clean:
	rm -rf build/*

build: clean
	CGO_ENABLED=0 go build -o ./build/bbgo ./cmd/cmd.go

test:
	go test ./...

run: build
	go run ./cmd/cmd.go
