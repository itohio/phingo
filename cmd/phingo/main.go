package main

import (
	"os"

	"github.com/itohio/phingo/pkg/repository"
)

var (
	globalRepository repository.Repository
)

func main() {
	rootCmd := newRootCmd(
		newInitCmd(),
		newAccountCmd(),
		newClientCmd(),
		newProjectCmd(),
		newInvoiceCmd(),
		newExportCmd(),
		newServeCmd(),
	)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
