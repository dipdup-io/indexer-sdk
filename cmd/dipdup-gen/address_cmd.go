package main

import (
	"github.com/spf13/cobra"
)

var (
	addressCmd = &cobra.Command{
		Use:                        "address [0x...]",
		Short:                      "generates application scaffolding by passed address",
		Args:                       cobra.MinimumNArgs(1),
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 1,
		Run:                        handleAddressCmd,
	}
)

func handleAddressCmd(cmd *cobra.Command, args []string) {
	cmdLineArgs.address = args[0]
}
