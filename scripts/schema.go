package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/shopware/shopware-cli/extension"
	"github.com/shopware/shopware-cli/shop"
)

var genericTypeRegex = regexp.MustCompile(`^(.+?)\[(.+)]$`)

func getSimpleTypeName(name string) string {
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}

func getGenericName(name string) string {
	if matches := genericTypeRegex.FindStringSubmatch(name); matches != nil {
		parent := getSimpleTypeName(matches[1])
		child := getGenericName(matches[2])
		return fmt.Sprintf("%s[%s]", parent, child)
	}
	return getSimpleTypeName(name)
}

func nameGenerics(r reflect.Type) string {
	return getGenericName(r.Name())
}

func generateProjectSchema() error {
	r := new(jsonschema.Reflector)
	r.FieldNameTag = "yaml"
	r.RequiredFromJSONSchemaTags = true

	if err := r.AddGoComments("github.com/shopware/shopware-cli", "./shop"); err != nil {
		return err
	}

	schema := r.Reflect(&shop.Config{})

	bytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile("shop/shopware-project-schema.json", bytes, 0o644); err != nil {
		return err
	}

	return nil
}

func generateExtensionSchema() error {
	r := jsonschema.Reflector{Namer: nameGenerics}
	r.FieldNameTag = "yaml"
	r.RequiredFromJSONSchemaTags = true

	if err := r.AddGoComments("github.com/shopware/shopware-cli", "./extension"); err != nil {
		return err
	}

	schema := r.Reflect(&extension.Config{})

	bytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile("extension/shopware-extension-schema.json", bytes, 0o644); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := generateProjectSchema(); err != nil {
		panic(err)
	}

	if err := generateExtensionSchema(); err != nil {
		panic(err)
	}
}
