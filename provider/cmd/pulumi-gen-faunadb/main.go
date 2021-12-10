package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	dotnetgen "github.com/pulumi/pulumi/pkg/v3/codegen/dotnet"
	gogen "github.com/pulumi/pulumi/pkg/v3/codegen/go"
	nodejsgen "github.com/pulumi/pulumi/pkg/v3/codegen/nodejs"
	pygen "github.com/pulumi/pulumi/pkg/v3/codegen/python"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	providerVersion "github.com/rawkode/pulumi-faunadb/provider/pkg/version"
)

// TemplateDir is the path to the base directory for code generator templates.
var TemplateDir string

// BaseDir is the path to the base pulumi-kubernetes directory.
var BaseDir string

func main() {
    flag.Usage = func() {
		const usageFormat = "Usage: %s <language> <swagger-or-schema-file> <root-pulumi-kubernetes-dir>"
		_, err := fmt.Fprintf(flag.CommandLine.Output(), usageFormat, os.Args[0])
		contract.IgnoreError(err)
		flag.PrintDefaults()
	}

    var version string
    flag.StringVar(&version, "version", providerVersion.Version, "the provider version to record in the generated code")
    flag.Parse()
    args := flag.Args()

    language, schemaPath, outDir := args[0], args[1], args[2]

    err := emitSDK(language, outDir, schemaPath, version)
	if err != nil {
		fmt.Printf("Failed: %s", err.Error())
	}
}

func emitSDK(language, outdir, schemaPath string, version string) error {
	pkg, err := readSchema(schemaPath, version)
	if err != nil {
		return err
	}

	tool := "Pulumi SDK Generator"
	extraFiles := map[string][]byte{}

	var generator func() (map[string][]byte, error)
	switch language {
	case "dotnet":
		generator = func() (map[string][]byte, error) { return dotnetgen.GeneratePackage(tool, pkg, extraFiles) }
	case "go":
		generator = func() (map[string][]byte, error) { return gogen.GeneratePackage(tool, pkg) }
	case "nodejs":
		generator = func() (map[string][]byte, error) { return nodejsgen.GeneratePackage(tool, pkg, extraFiles) }
	case "python":
		generator = func() (map[string][]byte, error) { return pygen.GeneratePackage(tool, pkg, extraFiles) }
	default:
		return errors.Errorf("Unrecognized language %q", language)
	}

	files, err := generator()
	if err != nil {
		return errors.Wrapf(err, "generating %s package", language)
	}

	for f, contents := range files {
		if err := emitFile(outdir, f, contents); err != nil {
			return errors.Wrapf(err, "emitting file %v", f)
		}
	}

	return nil
}

func readSchema(schemaPath string, version string) (*schema.Package, error) {
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading schema")
	}

	if strings.HasSuffix(schemaPath, ".yaml") {
		schemaBytes, err = yaml.YAMLToJSON(schemaBytes)
		if err != nil {
			return nil, errors.Wrap(err, "reading YAML schema")
		}
	}

	var spec schema.PackageSpec
	if err = json.Unmarshal(schemaBytes, &spec); err != nil {
		return nil, errors.Wrap(err, "unmarshalling schema")
	}
	spec.Version = version

	pkg, err := schema.ImportSpec(spec, nil)
	if err != nil {
		return nil, errors.Wrap(err, "importing schema")
	}
	return pkg, nil
}

func emitFile(rootDir, filename string, contents []byte) error {
	outPath := filepath.Join(rootDir, filename)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(outPath, contents, 0600); err != nil {
		return err
	}
	return nil
}
