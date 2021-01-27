package config

import (
	"flag"
	"strings"
)

var (
	flags flag.FlagSet

	// PkgName Generated Go package name
	PkgName = flags.String("pkgname", "tfschema", "Generated Go package name")
	types   = flags.String("types", "", "Explicit list of types to export")

	// Types Explicit list of types to export
	Types []string

	// Set exports FlagSet set function
	Set = flags.Set
)

// Finalize finalizes configuration
func Finalize() {
	// Split list of types by ":"
	if (types != nil) && (*types != "") {
		Types = strings.Split(*types, ":")
	}
}
