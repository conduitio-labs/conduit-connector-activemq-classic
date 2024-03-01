.PHONY: build test clean generate install-paramgen install-tools golangci-lint-install

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/alarbada/conduit-connector-activemq-classic.version=${VERSION}'" -o activemq cmd/connector/main.go

generate:
	go generate ./...

install-paramgen:
	go install github.com/conduitio/conduit-connector-sdk/cmd/paramgen@latest

up:
	docker compose -f test/docker-compose.yml up activemq activemq-tls --quiet-pull -d --wait 

up-dev:
	docker compose -f test/docker-compose.yml up activemq-dev --quiet-pull -d --wait 

down:
	docker compose -f test/docker-compose.yml down -v --remove-orphans

lint:
	golangci-lint run

test: up
	go test -v -count=1 -race .; ret=$$?; \
		docker compose -f test/docker-compose.yml down && \
		exit $$ret

