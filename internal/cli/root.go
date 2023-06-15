package cli

import (
	"fmt"
	"os"

	"github.com/FeLvi-zzz/tentez"
	"github.com/spf13/cobra"
)

var filename = ""
var noPause = false

var rootCmd = &cobra.Command{
	Use:     "tentez",
	Short:   "Tentez helps you switching traffic.",
	Version: fmt.Sprintf("%s (rev: %s)\n", tentez.Version, tentez.Revision),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
