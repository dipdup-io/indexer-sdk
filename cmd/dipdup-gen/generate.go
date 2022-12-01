package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	js "github.com/dipdup-net/indexer-sdk/pkg/jsonschema"
	"github.com/iancoleman/strcase"
)

type generateResult struct {
	methods map[string]string
	events  map[string]string
	models  map[string]goType
}

func newResult() *generateResult {
	return &generateResult{
		methods: make(map[string]string),
		events:  make(map[string]string),
		models:  make(map[string]goType),
	}
}

func (r *generateResult) merge(other *generateResult) {
	if other == nil {
		return
	}
	for key, value := range other.methods {
		r.methods[key] = value
	}
	for key, value := range other.events {
		r.events[key] = value
	}
	for key, value := range other.models {
		r.models[key] = value
	}
}

// String -
func (r *generateResult) String() string {
	var builder strings.Builder
	builder.WriteString("Methods:\r\n\t")
	for key := range r.methods {
		builder.WriteString(key)
		builder.WriteByte(',')
	}
	builder.WriteString("\r\nEvents:\r\n\t")
	for key := range r.events {
		builder.WriteString(key)
		builder.WriteByte(',')
	}
	builder.WriteString("\r\nModels:\r\n\t")
	for key := range r.models {
		builder.WriteString(key)
		builder.WriteByte(',')
	}
	builder.WriteString("\r\n")
	return builder.String()
}

func generate(args cmdLine, schema []map[string]js.Type) error {
	dirs, err := createProjectDirs(args.destination, args.appName)
	if err != nil {
		return err
	}

	types := newResult()
	for _, s := range schema {
		result, err := generateModels(args, dirs, s)
		if err != nil {
			return err
		}
		types.merge(result)
	}

	return applyTemplates(args, dirs, types)
}

func applyTemplates(args cmdLine, dirs *projectDirs, result *generateResult) error {
	if err := generateFromTemplate("core.go", "postgres_core.tmpl", dirs.postgres, map[string]any{
		"Models":      result.models,
		"PackageName": args.packageName,
	}, true); err != nil {
		return err
	}

	if err := generateFromTemplate("postgres.go", "postgres_module.tmpl", dirs.cmd, map[string]any{
		"Models":      result.models,
		"PackageName": args.packageName,
	}, true); err != nil {
		return err
	}

	if err := generateFromTemplate("decoder.go", "decoder.tmpl", dirs.cmd, map[string]any{
		"Events":      result.events,
		"Methods":     result.methods,
		"PackageName": args.packageName,
		"Models":      result.models,
	}, true); err != nil {
		return err
	}

	if err := generateFromTemplate("main.go", "main.tmpl", dirs.cmd, map[string]any{
		"Address":     args.address,
		"Models":      result.models,
		"PackageName": args.packageName,
		"App":         args.appName,
	}, true); err != nil {
		return err
	}

	if err := generateFromTemplate("dipdup.yml", "config.tmpl", dirs.build, map[string]any{
		"Address": args.address,
		"App":     args.appName,
	}, false); err != nil {
		return err
	}

	if err := generateFromTemplate("Dockerfile", "Dockerfile.tmpl", dirs.build, map[string]any{
		"App":         args.appName,
		"PackageName": args.packageName,
	}, false); err != nil {
		return err
	}

	if err := generateFromTemplate("docker-compose.yml", "compose.tmpl", dirs.root, map[string]any{
		"App": args.appName,
	}, false); err != nil {
		return err
	}

	if err := generateFromTemplate("Makefile", "makefile.tmpl", dirs.root, map[string]any{
		"App": args.appName,
	}, false); err != nil {
		return err
	}

	return nil
}

func generateModels(args cmdLine, dirs *projectDirs, schema map[string]js.Type) (*generateResult, error) {
	result := newResult()

	for name, entity := range schema {
		types := make(map[string]goType)
		inputName := buildName(name, entity.Type)

		entityType := generateTypes(name, entity.Type, entity.Inputs, types)
		delete(types, entityType.Name)

		if err := generateStorageBySchema(entityType, inputName, dirs.storage, "model.tmpl", args.packageName, types); err != nil {
			return nil, err
		}
		if err := generateStorageBySchema(entityType, inputName, dirs.postgres, "postgres.tmpl", args.packageName, nil); err != nil {
			return nil, err
		}

		result.models[inputName] = entityType

		switch entity.Type {
		case "method":
			result.methods[inputName] = entity.Signature
		case "event":
			result.events[inputName] = entity.Signature
		}
	}

	return result, nil
}

func generateStorageBySchema(typ goType, name, dir, templateName, packageName string, types map[string]goType) error {
	return generateFromTemplate(strcase.ToSnake(name)+".go", templateName, dir, map[string]any{
		"GoType":      typ,
		"PackageName": packageName,
		"Nested":      types,
	}, true)
}

func generateFromTemplate(outputName, templateFileName, dest string, ctx interface{}, needFormat bool) error {
	tmpl, err := template.New("storage").Funcs(template.FuncMap{
		"LowerCamelCase": strcase.ToLowerCamel,
	}).ParseFS(templates, "templates/*/*")
	if err != nil {
		return err
	}
	targetFile := filepath.Join(dest, outputName)
	templateFile, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer templateFile.Close()

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, templateFileName, ctx); err != nil {
		return err
	}

	var result []byte
	if needFormat {
		result, err = format.Source(buf.Bytes())
		if err != nil {
			return err
		}
	} else {
		result = buf.Bytes()
	}

	_, err = templateFile.Write(result)
	return err
}

func buildName(name, postfix string) string {
	return strcase.ToCamel(fmt.Sprintf("%s_%s", name, postfix))
}
