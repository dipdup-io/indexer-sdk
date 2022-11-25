package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/dipdup-net/indexer-sdk/pkg/contract"
	js "github.com/dipdup-net/indexer-sdk/pkg/jsonschema"
	"github.com/spf13/cobra"
)

var (
	abiCmd = &cobra.Command{
		Use:                        "abi [path to file or directory]",
		Short:                      "generates application scaffolding by passed ABI",
		Args:                       cobra.MinimumNArgs(1),
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 1,
		RunE:                       handleAbiCmd,
	}
)

func handleAbiCmd(cmd *cobra.Command, args []string) error {
	fi, err := os.Lstat(args[0])
	if err != nil {
		return err
	}

	types := make([]map[string]js.Type, 0)
	if fi.IsDir() {
		entries, err := os.ReadDir(args[0])
		if err != nil {
			return err
		}
		for i := range entries {
			if entries[i].IsDir() {
				continue
			}
			if filepath.Ext(entries[i].Name()) != ".json" {
				continue
			}
			schema, err := schemaFromFile(entries[i].Name())
			if err != nil {
				return err
			}
			types = append(types, schema)
		}
	} else {
		schema, err := schemaFromFile(fi.Name())
		if err != nil {
			return err
		}
		types = append(types, schema)
	}

	return generate(cmdLineArgs, types)
}

func schemaFromFile(filename string) (map[string]js.Type, error) {
	log.Printf("found abi file: %s", filename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	abiData, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	data, err := contract.JSONSchema(contract.TypeEvm, abiData)
	if err != nil {
		return nil, err
	}

	var schema map[string]js.Type
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}

	return schema, nil
}
