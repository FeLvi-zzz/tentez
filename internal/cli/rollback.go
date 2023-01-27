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
		filename := cmd.Flag("filename").Value.String()
		t, err := tentez.NewFromYaml(filename)
		if err != nil {
			return err
		}

		return t.Rollback(true)
	},
}

func init() {
	rollbackCmd.Flags().StringVarP(&filename, "filename", "f", "", "config file for tentez")
	if err := rollbackCmd.MarkFlagRequired("filename"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(rollbackCmd)
}
