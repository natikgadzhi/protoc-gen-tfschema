package config

import (
	"flag"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	// Package name used by default for generated Terraform schema
	defaultPackageName = "tfschema"
)

var (
	flags flag.FlagSet

	// PackageName Generated Go package name
	PackageName = flags.String("pkgname", defaultPackageName, "Generated Go package name")
	types       = flags.String("types", "", "Explicit list of types to export")

	// Types Explicit list of types to export
	Types []string

	// ProtocVersion is protoc version as a string
	ProtocVersion string

	// Set exports FlagSet set function for protobuf option parser
	Set = flags.Set
)

// Finalize finalizes configuration
func Finalize() {
	// Split list of types by ":"
	if (types != nil) && (*types != "") {
		Types = strings.Split(*types, ":")
	}

	dump()
}

func dump() {
	log.Infof("Package name: %s", *PackageName)
	log.Infof("Types to export: %s", Types)
	log.Infof("Protoc version: %s", ProtocVersion)
}
