EXECUTABLE := vendopunkto-server
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
MIGRATIONDIR := store/postgres/migrations
MIGRATIONS :=  $(wildcard ${MIGRATIONDIR}/*.sql)
TOOLS := ${GOPATH}/bin/mockery \
    ${GOPATH}/bin/wire

.PHONY: default
default: ${EXECUTABLE}

${GOPATH}/bin/mockery:
	go get github.com/vektra/mockery/internal/cmd/mockery

${GOPATH}/bin/wire:
	go get github.com/google/wire
	go get github.com/google/wire/internal/cmd/wire

tools: ${TOOLS}

internal/cmd/wire_gen.go: internal/cmd/wire.go
	wire ./internal/cmd/...

.PHONY: mocks
mocks: tools
	mockery -dir ./gorestapi -name ThingStore

.PHONY: ${EXECUTABLE}
${EXECUTABLE}: tools internal/cmd/wire_gen.go 
	# Compiling...
	go build -ldflags "-X ${PACKAGENAME}/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}

.PHONY: test
test: 
	go test -cover ./...

.PHONY: deps
deps:
	# Fetching dependancies...
	go get -d -v # Adding -u here will break CI

dbclean:
	docker stop vendopunktopostgres; 
	docker rm vendopunktopostgres;
	docker volume rm vendopunkto_db;
	docker-compose up -d vendopunktopostgres

run:
	STORAGE_HOST=localhost \
	PLUGINS_ENABLED="http://localhost:4000" \
	PLUGINS_DEFAULT_EXCHANGE_RATES=fake-exchange-rates \
	LOGGER_LEVEL=debug \
	${GOPATH}/src/${PACKAGENAME}/vendopunkto-server api


build-cli:
	go build -o ./vendopunkto-cli ./cli/main.go


plugins/monero/internal/wire_gen.go: plugins/monero/internal/wire.go
	wire ./plugins/monero/internal/...

build-monero: tools plugins/monero/internal/wire_gen.go
	go build -o ./vendopunkto-monero ./plugins/monero/main.go

run-monero:
	MONERO_WALLET_RPC_URL=http://localhost:18082 \
	${GOPATH}/src/${PACKAGENAME}/vendopunkto-monero


plugins/exchange-rates/internal/wire_gen.go: plugins/exchange-rates/internal/wire.go
	wire ./plugins/exchange-rates/internal/...

build-rates: tools plugins/exchange-rates/internal/wire_gen.go
	go build -o ./vendopunkto-rates ./plugins/exchange-rates/main.go

run-rates:
	${GOPATH}/src/${PACKAGENAME}/vendopunkto-rates


plugins/development/internal/wire_gen.go: plugins/development/internal/wire.go
	wire ./plugins/development/internal/...

dev-plugin: tools plugins/development/internal/wire_gen.go
	go build -o ./vendopunkto-dev ./plugins/development/main.go
	${GOPATH}/src/${PACKAGENAME}/vendopunkto-dev