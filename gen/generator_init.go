package gen

import (
	"fmt"
	"strings"
)

const (
	defaultPackageName = "tfschema"
)

// TODO: Move to generator after splitting up file generator and base generator

// parseOptions fills configuration from command line args
func (g *Generator) parseArgs() {
	var dict map[string]string = make(map[string]string)

	fields := strings.Fields(*g.Plugin.Request.Parameter)
	for _, value := range fields {
		parts := strings.Split(value, "=")
		if len(parts) > 1 {
			key := strings.ToLower(parts[0])
			value := parts[1]
			dict[key] = value
		}
	}

	g.packageName = dict["pkgname"]
	g.types = strings.Split(dict["types"], ",")

	if g.packageName == "" {
		g.packageName = defaultPackageName
	}
}

// fillVersion saves protoc version number as string
func (g *Generator) fillVersion() {
	if v := g.Plugin.Request.GetCompilerVersion(); v != nil {
		g.protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
	}
}
