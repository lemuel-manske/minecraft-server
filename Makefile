.PHONY: init build install test

init:
	terraform -chdir=tf init

test:
	cd cli && go test ./...

build:
	cd cli && go build -o ../mc .

install: build
	sudo cp mc /usr/local/bin/mc
	rm mc
