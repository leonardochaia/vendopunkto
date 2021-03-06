EXECUTABLE := vendopunkto-server
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
MIGRATIONDIR := store/postgres/migrations
MIGRATIONS :=  $(wildcard ${MIGRATIONDIR}/*.sql)
TOOLS := ${GOPATH}/bin/wire

.PHONY: default
default: ${EXECUTABLE}

${GOPATH}/bin/wire:
	go get github.com/google/wire/cmd/wire

tools: ${TOOLS}

internal/cmd/wire_gen.go: internal/cmd/wire.go
	wire ./internal/cmd/...

.PHONY: ${EXECUTABLE}
${EXECUTABLE}: tools internal/cmd/wire_gen.go 
	# Compiling...
	go build -ldflags "-X ${PACKAGENAME}/internal/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/internal/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}

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
	PLUGIN_HOSTS="http://localhost:4000" \
	EXCHANGE_RATES_PLUGIN=fake-exchange-rates \
	CURRENCY_METADATA_PLUGIN=fake-currency-metadata \
	LOGGER_LEVEL=debug \
	${GOPATH}/src/${PACKAGENAME}/vendopunkto-server api

run-vp-monero:
	STORAGE_HOST=localhost \
	PLUGIN_HOSTS="http://localhost:4200 http://localhost:4201" \
	EXCHANGE_RATES_PLUGIN=gecko-exchange-rates \
	LOGGER_LEVEL=debug \
	${GOPATH}/src/${PACKAGENAME}/vendopunkto-server api


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