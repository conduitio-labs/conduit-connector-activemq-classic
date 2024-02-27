.PHONY: build test test-integration generate install-paramgen install-tools golangci-lint-install

VERSION=$(shell git describe --tags --dirty --always)


build:
	go build -ldflags "-X 'github.com/alarbada/conduit-connector-activemq.version=${VERSION}'" -o activemq cmd/connector/main.go

generate:
	go generate ./...

install-paramgen:
	go install github.com/conduitio/conduit-connector-sdk/cmd/paramgen@latest

up:
	docker compose -f test/docker-compose.yml up activemq-dev --quiet-pull -d --wait 

down:
	docker compose -f test/docker-compose.yml down -v --remove-orphans

test-acceptance:
	make up
	go test -v -run TestAcceptance .
	make down
