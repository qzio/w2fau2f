SHELL := /bin/bash

TARGET := $(shell echo $${PWD\#\#*/})
#TARGET := $(shell echo $(basename  $(pwd)))

.DEFAULT_GOAL := $(TARGET)

#VERSION := 0.1.0
#BUILD := $(shell git rev-parse HEAD)
#LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
SRC = $(shell find . -type f -name '*.go')

.PHONY: all build clean install uninstall fmt simplify check test run

all: clean check test build install

$(TARGET): $(SRC)
	$(info building)
	@go build

run:
	$(info build and run)
	@go build
	./$(TARGET)

clean:
	@go clean
