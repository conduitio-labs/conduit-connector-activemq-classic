.PHONY: build test clean generate install-paramgen install-tools golangci-lint-install

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-activemq-classic.version=${VERSION}'" -o activemq-classic cmd/connector/main.go

generate:
	go generate ./...
	conn-sdk-cli readmegen -w

.PHONY: install-tools
install-tools:
	@echo Installing tools from tools.go
	@go list -e -f '{{ join .Imports "\n" }}' tools.go | xargs -I % go list -f "%@{{.Module.Version}}" % | xargs -tI % go install %
	@go mod tidy

up:
	docker compose -f test/docker-compose.yml up --quiet-pull -d --wait

down:
	docker compose -f test/docker-compose.yml down -v --remove-orphans

lint:
	golangci-lint run

test: up
	go test -v -count=1 -race .; ret=$$?; \
		docker compose -f test/docker-compose.yml down && \
		exit $$ret

