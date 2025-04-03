VERSION=$(shell git describe --tags --dirty --always)

.PHONY: build
build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-activemq-classic.version=${VERSION}'" -o activemq-classic cmd/connector/main.go

.PHONY: test
test: up
	go test -v -count=1 -race .; ret=$$?; \
		docker compose -f test/docker-compose.yml down && \
		exit $$ret

.PHONY: generate
generate:
	go generate ./...
	conn-sdk-cli readmegen -w

.PHONY: install-tools
install-tools:
	@echo Installing tools from tools/go.mod
	@go list -modfile=tools/go.mod tool | xargs -I % go list -modfile=tools/go.mod -f "%@{{.Module.Version}}" % | xargs -tI % go install %
	@go mod tidy

.PHONY: up
up:
	docker compose -f test/docker-compose.yml up --quiet-pull -d --wait

.PHONY: down
down:
	docker compose -f test/docker-compose.yml down -v --remove-orphans

.PHONY: lint
lint:
	golangci-lint run
