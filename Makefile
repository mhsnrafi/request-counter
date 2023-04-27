BUILD_DIR := $(shell pwd)/output

.PHONY: all server build test clean

all: build

build: server

server:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/server ./cmd

test:
	go test -race ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -f counts.json
