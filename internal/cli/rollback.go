package cli

import (
	"github.com/FeLvi-zzz/tentez"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback switch, switch old:new = 100:0",
	Long: `
# rollback
$ tentez -f ./examples/example.yaml rollback`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := cmd.Flag("filepath").Value.String()
		t, err := tentez.NewFromYaml(filepath)
		if err != nil {
			return err
		}

		return t.Rollback(true)
	},
}

func init() {
	rollbackCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "config file for tentez")
	if err := rollbackCmd.MarkFlagRequired("filepath"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(rollbackCmd)
}
