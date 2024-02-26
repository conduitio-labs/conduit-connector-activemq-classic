.PHONY: build test test-integration generate install-paramgen install-tools golangci-lint-install

VERSION=$(shell git describe --tags --dirty --always)


build:
	go build -ldflags "-X 'github.com/alarbada/conduit-connector-activemq.version=${VERSION}'" -o activemq cmd/connector/main.go

up:
	docker compose -f test/docker-compose.yml up activemq-dev --quiet-pull -d --wait 

down:
	docker compose -f test/docker-compose.yml down -v --remove-orphans
