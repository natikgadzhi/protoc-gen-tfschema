package main

import (
	"fmt"
	"strings"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	// Package name used by default for generated Terraform schema
	defaultPackageName = "tfschema"
)

// Generator is the Terraform Schema generator
type Generator struct {
	Plugin        *protogen.Plugin // Plugin
	packageName   string           // Package name
	types         []string         // Types to export
	protocVersion string           // protoc version
}

// NewGenerator returns new generator instance.
func NewGenerator(plugin *protogen.Plugin) *Generator {
	g := &Generator{Plugin: plugin}

	g.parseArgs()
	g.fillProtocVersion()

	log.Infof("Command line options: %s", g.Plugin.Request.GetParameter())
	log.Infof("Package name: %s", g.packageName)
	log.Infof("Types: %s", g.types)
	log.Infof("Protoc version: %s", g.protocVersion)

	return g
}

// fillVersion saves protoc version number as string
// NOTE: questionable!
func (g *Generator) fillProtocVersion() {
	if v := g.Plugin.Request.GetCompilerVersion(); v != nil {
		g.protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
	}
}

// parseOptions fills configuration from command line args
func (g *Generator) parseArgs() {
	var dict map[string]string = make(map[string]string)

	fields := strings.Fields(*g.Plugin.Request.Parameter)
	for _, value := range fields {
		parts := strings.Split(value, "=")
		if len(parts) > 1 {
			key := strings.ToLower(parts[0])
			value := parts[1]
			dict[key] = strings.Trim(value, " ")
		}
	}

	g.packageName = dict["pkgname"]
	if g.packageName == "" {
		g.packageName = defaultPackageName
	}

	if dict["types"] != "" {
		g.types = strings.Split(dict["types"], ",")
	}
}

// Generate is an entry point for generator
func (g *Generator) Generate() ([]*protogen.GeneratedFile, error) {
	var files []*protogen.GeneratedFile

	for _, file := range g.Plugin.Files {
		if !file.Generate {
			continue
		}

		// Target file name
		filename := file.GeneratedFilenamePrefix + ".tfschema.go"
		out := g.Plugin.NewGeneratedFile(filename, ".")

		//result = reader.ReadFile(file, g.types)
		// parsed = Parser.file
		// result = Renderer.file

		//_, err := out.Write(result.Bytes())
		_, err := out.Write([]byte("ok"))
		if err != nil {
			log.Errorf("Couldn't write %s, error: %v", filename, err)
			return files, trace.Wrap(err)
		}

		files = append(files, out)
		log.Infof("Will emit %s", filename)
	}

	return files, nil
}
