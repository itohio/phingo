package main

import (
	"os"

	"github.com/itohio/phingo/pkg/repository"
)

var (
	flagRepositoryUrl *string
	flagAccount       *string
	flagProject       *string

	globalRepository repository.Repository
)

func main() {
	rootCmd := newRootCmd(
		newInitCmd(),
		newProjectCmd(),
		newClockCmd(),
		newInvoiceCmd(),
		newExportCmd(),
		newServeCmd(),
	)

	flagRepositoryUrl = rootCmd.PersistentFlags().StringP("repository", "r", ".phingo", "Data Repository path.")
	flagAccount = rootCmd.PersistentFlags().StringP("account", "a", "", "Account id")
	flagProject = rootCmd.PersistentFlags().StringP("project", "p", "", "Project id")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
