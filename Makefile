.PHONY: clean

clean:
	@mkdir -p build
	rm -rf ./build/*
	go clean

.PHONY: build
build: clean
	go build -o build/protoc-gen-tfschema

.PHONY: install
install: build
	go install .

# make example
# Clones Teleport repo to go/src, builds Teleport schema output with protoc
include example.mk