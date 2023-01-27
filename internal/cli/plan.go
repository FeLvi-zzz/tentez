package cli

import (
	"github.com/FeLvi-zzz/tentez"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show steps how to apply",
	Long: `
# show plan
$ tentez -f ./examples/example.yaml plan`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := cmd.Flag("filename").Value.String()
		t, err := tentez.NewFromYaml(filename)
		if err != nil {
			return err
		}

		return t.Plan()
	},
}

func init() {
	planCmd.Flags().StringVarP(&filename, "filename", "f", "", "config file for tentez")
	if err := planCmd.MarkFlagRequired("filename"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(planCmd)
}
