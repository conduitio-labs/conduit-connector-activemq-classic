.PHONY: build test clean generate install-paramgen install-tools golangci-lint-install

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/alarbada/conduit-connector-activemq.version=${VERSION}'" -o activemq cmd/connector/main.go

generate:
	go generate ./...

install-paramgen:
	go install github.com/conduitio/conduit-connector-sdk/cmd/paramgen@latest

up:
	docker compose -f test/docker-compose.yml up activemq activemq-tls --quiet-pull -d --wait 

down:
	docker compose -f test/docker-compose.yml down -v --remove-orphans

up-tls: clean-tls setup-tls
	docker compose -f test/docker-compose.yml up activemq-tls --quiet-pull -d --wait 

setup-tls:
	cd test && ./setup-tls.sh

clean-tls:
	rm -rf test/certs

lint:
	golangci-lint run

test:
	rm -rf test/certs
	cd test && ./setup-tls.sh
	docker compose -f test/docker-compose.yml up activemq \
		--quiet-pull -d --wait
	go test -v -count=1 -race .; ret=$$?; \
		docker compose -f test/docker-compose.yml down && \
		rm -rf test/certs && \
		exit $$ret



