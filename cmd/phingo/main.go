package main

import (
	"os"
	"strings"

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

func parseKeyValue(out map[string]string, vals []string) {
	for _, c := range vals {
		kv := strings.SplitN(c, "=", 2)
		if len(kv) == 2 && kv[1] != "" {
			out[kv[0]] = kv[1]
		} else {
			delete(out, kv[0])
		}
	}
}
