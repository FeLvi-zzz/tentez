package cli

import (
	"fmt"
	"os"

	"github.com/FeLvi-zzz/tentez"
	"github.com/spf13/cobra"
)

var filepath = ""

var rootCmd = &cobra.Command{
	Use:   "tentez",
	Short: "Tentez helps you switching traffic.",
	Long: `Tentez helps you switching traffic.
# show plan
$ tentez -f ./examples/example.yaml plan

# show plan and apply
$ tentez -f ./examples/example.yaml apply

# get target resources' current states.
$ tentez -f ./examples/example.yaml get

# rollback
$ tentez -f ./examples/example.yaml rollback

# show version
$ tentez version
`,
	Version: fmt.Sprintf("%s (rev: %s)\n", tentez.Version, tentez.Revision),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
