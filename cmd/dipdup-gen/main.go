package main

import (
	"embed"
	"log"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	//go:embed templates/*/*.tmpl
	templates embed.FS
)

var (
	rootCmd = &cobra.Command{
		Use:   "dipdup-gen",
		Short: "DipDup generator for EVM rollup",
	}
)

type cmdLine struct {
	appName     string
	packageName string
	destination string
	address     string
}

var cmdLineArgs cmdLine

func main() {
	rootCmd.PersistentFlags().StringVarP(&cmdLineArgs.appName, "cmd", "c", "app", "application name (default: app)")
	rootCmd.PersistentFlags().StringVarP(&cmdLineArgs.packageName, "package", "p", "", "package name (example: github.com/organization/repository-name)")
	rootCmd.PersistentFlags().StringVarP(&cmdLineArgs.destination, "output", "o", ".", "output directory (default: current directory)")
	rootCmd.AddCommand(abiCmd, addressCmd)

	abiCmd.Flags().StringVarP(&cmdLineArgs.address, "address", "a", "", "contract address (example: 0x002b1ee9B1CF77233F9f96Fc9ee6191D2b857Be2)")
	if err := abiCmd.MarkFlagRequired("address"); err != nil {
		log.Panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}
