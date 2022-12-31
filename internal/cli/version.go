package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long: `
# show version
$ tentez version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tentez version %s\n", cmd.Root().Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
