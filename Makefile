.DEFAULT_GOAL := build

.PHONY:fmt vet build run
fmt:
	go fmt ./...
vet: fmt
	go vet ./...
build: vet
	go build
install: build
	mv -f ./terminal-doodle /Users/matejtopolovac/go/bin/td

.PHONY:clean
clean:
	go clean