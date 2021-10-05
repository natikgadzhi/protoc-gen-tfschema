## This repo is a spike, and is not maintained

Check out [Teleport Plugins](https://github.com/gravitational/teleport-plugins) if you need Teleport's Terraform plugin.  

# protoc-gen-tfschema

protoc-gen-tfschema is a protoc plugin that generates Terraform-compatible
schemas that you can use to build your own Terraform providers for a service
that already has a GRPC API and defines it's with Protobuf.

protoc-gen-schema is in active development. Your contributions are very welcome.
Maybe don't use it in production yet, as the schema output format may change.

## Usage

You'd need several components to build a Terraform provider:

- your API client and settings to authenticate it
- resource schema definitions
- funcitons to create, update, delete, and read the data from your API
- as a part of those funcitons, marshalling and unmarshalling data between
  terraform and your API formats.

If your target app provides `.proto` description of it's API data types, you can
use `protoc-gen-tfschema` to generate the resource definitions and
marshalling/unmarshalling components automatically.

You'll still need to provide an API client initialization code, and call the API
to execute CRUD features.

We assume that protobuf & protobuf-gen-go are installed on your system.

Here's how to generate resource schemas:

```bash

go install github.com/nategadzhi/protoc-gen-tfschema

protoc --proto_path ./proto/teleport --proto_path=. --proto_path=/Users/xnutsive/go/src/github.com/gravitational/teleport/vendor/github.com/gogo/protobuf --proto_path=/Users/xnutsive/go/src types.proto --tfschema_out=./out/teleport --go_out=./out/teleport --tfschema_opt="types=Metadata,UserSpecV2,UserV2"
```

You also could clone this repo and run:

```
make example # Builds example schema for Terraform
```

### Command line options

`--tfschema_opt=` accepts several options:

- `types` defines what Protobuf Messages (types) to generate resource schemas
  for.
- `pkgname` defines what package name to use for the resulting files. Defaults
  to `tfschema`.

## Features

- **Resource Schemas**: generate schemas for your resources with a single CLI
  command.
- **Unmarshalling from Terraform data**: generates functions that unmarshall
  data form Terraform to Go types generated from the Protobuf messages.

### Things that are NOT supported yet

- The plugin can't generate `Computed` flag for schema fields. There's no way to
  know that the field is computed from the proto.

### Next steps

- [*] An example project to test the generated schemas
- [ ] Provide unmarshalling generators
- [ ] Provide terraform data source schema generators
- [ ] Provide terraform schema importers, versions, and updaters
