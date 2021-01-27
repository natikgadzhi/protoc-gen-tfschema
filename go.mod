module github.com/nategadzhi/protoc-gen-tfschema

go 1.15

require (
	github.com/gravitational/trace v1.1.13
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/sanity-io/litter v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/sanity-io/litter => github.com/gzigzigzeo/litter v1.3.1-0.20210126033604-e45922096d71
