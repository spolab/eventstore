.PHONY: all docker start stop ps build

BIN := bin
DEPS := $(BIN) $(wildcard cmd/**/*) $(wildcard pkg/**/*) $(wildcard ../schema/*)

all: $(BIN)/eventstore build start

$(BIN):
	mkdir -p $@

$(BIN)/eventstore: $(DEPS)
	protoc -I../schema --go_out=paths=source_relative:gen --go-grpc_out=paths=source_relative:gen ../schema/eventstore.proto
	go build -o bin/eventstore cmd/main.go

start:
	docker compose -f docker/docker-compose.yml up -d

stop:
	docker compose -f docker/docker-compose.yml down

build:
	docker compose -f docker/docker-compose.yml build

ps:
	docker compose -f docker/docker-compose.yml ps