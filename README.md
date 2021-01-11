# protoc-gen-schema

Generates Terraform Resource Schemas from proto. Heads up, this is an
implementation note, not a public-release-ready readme.

## Usage

```bash

go install github.com/nategadzhi/protoc-gen-schema

protoc --proto_path ./proto/teleport --proto_path=. --proto_path=/Users/xnutsive/go/src/github.com/gravitational/teleport/vendor/github.com/gogo/protobuf --proto_path=/Users/xnutsive/go/src types.proto --tfschema_out=./out/teleport --go_out=./out/teleport --tfschema_opt="types=Metadata UserSpecV2 UserV2"

```

## Working with complex field types

When rendering a schema, fields can be represented in several different ways:

1. A field can be of a basic type (Bool, String, Int). In that case, the plugin
   should convert the field to it's basic type counterpart in terraform. Lookup
   available Terraform types to do that.
2. A field can be a list or set of a primitive type. Then the schema will look a
   bit different, but the basic premise is the same, convert primitive type to
   Terraform notation, and wrap it with the Terraform notation of a list.
3. A field can be of a complex custom type. Then we should have a map of said
   types in memory and we should generate an in-place Terraform resource for
   that type.

### Plugin Workflow

Assuming plugin is initialized and we have all of the options provided, here's
how the plugin should run and generate type schemas.

1. protoc API provides all the files in a loop. The generator reads the files
   and saves all the messages into a set. The messages are in a standard protoc
   format.

2. For each message, generator checks if that message is required to be
   generated.

3. For required messages, generator tries to generate the message schema
   function without all the wraparound code at this point.

4. The template goes over message fields. Pritimitive fields can be simply
   rendered, whereas a message_type fields need to be generated as a partial or
   a function that calls another template.

5. That message_type partial takes in the message name and message spec from the
   map of messages, and resolves it's fields recursively.

   The only difference is that it generates a schema.Resource with a
   schema.Schema inside istead of just the schema.Schema. So the internal logic
   can be extracted into a helper function.

6. If some message_type can't be resolved (is not in all of the provided proto
   files) — error out.

Some intermidiate steps that can be implemented:

1. Generate just the primitive types
2. Add lists, sets, and maps of primitive types

3. Implement recursive message_type
4. Implement lists, sets, and maps of message_types

## Additional work on fields

- [ ] Is there a way to set fields as computed, or leave it out for now?
- [ ] Is there a way to setup default values?
