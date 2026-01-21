.DEFAULT_GOAL := build
.PHONY: fmt vet build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o bin/statement-parser cmd/main.go

test:
	go test ./...

clean:
	rm -rf bin/
