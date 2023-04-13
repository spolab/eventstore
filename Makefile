.PHONY: image clean build

PROTO_NAMESPACE=toremo.com/eventstore/v1
GEN_DIR := gen
BIN_DIR := bin
PKG_DIR := pkg
SCHEMA_DIR := schema
PROTOC_OUT_OPTS := M=$(PROTO_NAMESPACE),paths=source_relative:$(GEN_DIR)
BIN_DEPS := cmd/main.go gen/eventstore.pb.go $(wildcard $(PKG_DIR)/**/*.go)

build: $(BIN_DIR)/eventstore
	docker build -f docker/Dockerfile --tag toremo/eventstore:latest .

image: 
	docker build -f docker/Dockerfile --tag toremo/eventstore:latest .

clean:
	rm -Rf $(GEN_DIR) $(BIN_DIR)

$(BIN_DIR)/eventstore: $(BIN_DEPS)
	go mod download
	go test ./...
	gosec ./...
	CGO_ENABLED=0 go build -o $@ $<

$(GEN_DIR):
	mkdir -p $@

$(GEN_DIR)/%.pb.go: $(SCHEMA_DIR)/%.proto $(GEN_DIR)
	protoc -I$(SCHEMA_DIR) --go_out=$(PROTOC_OUT_OPTS) --go-grpc_out=$(PROTOC_OUT_OPTS) $<
