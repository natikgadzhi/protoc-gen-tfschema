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
include example.mk