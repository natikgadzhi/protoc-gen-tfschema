# Clones Teleport and builds tfschema. protoc and protobuf must be installed.
gopath = $(shell go env GOPATH)
srcpath = $(gopath)/src
teleport_url = github.com/gravitational/teleport
teleport_repo = https://$(teleport_url)
teleport_dir = $(srcpath)/$(teleport_url)
out_dir = "./_out"

.PHONY: example
example: build
ifeq ("$(wildcard $(teleport_dir))", "")
	@echo "Teleport source code is required to build this example!"
	@echo "Cloning teleport repo $(teleport_repo) to $(teleport_dir)"
	@git clone $(teleport_repo) $(teleport_dir)
endif
	@echo "Teleport has been downloaded."
	@protoc \
		-I$(teleport_dir)/api/types \
		-I$(teleport_dir)/vendor/github.com/gogo/protobuf \
		-I$(teleport_dir) \
		-I$(srcpath) \
		--tfschema_out=$(out_dir) \
		--go_out=$(out_dir) \
		--plugin=./build/protoc-gen-tfschema \
		--tfschema_opt="pkgname=tfschema,types=UserSpecV2:Metadata:UserV2" \
		types.proto
