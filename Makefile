ifneq (,$(wildcard .env))
include .env
export
endif


BIN_NAME ?= mc
BIN_HOME ?= /usr/local/bin


.PHONY: init build install test


login:
	aws login --remote

init:
	terraform -chdir=tf init


fmt:
	cd cli && go fmt ./...

test:
	cd cli && go test ./...

check: fmt test


build:
	cd cli && go build -o ../$(BIN_NAME) .

install: build
	sudo cp $(BIN_NAME) $(BIN_HOME)/$(BIN_NAME)
	rm $(BIN_NAME)
