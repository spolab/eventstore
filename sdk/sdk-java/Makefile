BIN := target
SCHEMA := ../../schema

all: $(BIN)
	protoc -I$(SCHEMA) --java_out=$(BIN) $(SCHEMA)/*.proto

$(BIN):
	mkdir -p $@