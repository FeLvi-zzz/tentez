package cli

import (
	"github.com/FeLvi-zzz/tentez"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Switch targets weights",
	Long: `
# show plan and apply
$ tentez -f ./examples/example.yaml apply`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		t, err := tentez.NewFromYaml(ctx, filename)
		if err != nil {
			return err
		}

		if err := t.Plan(); err != nil {
			return err
		}
		return t.Apply(ctx, false)
	},
}

func init() {
	applyCmd.Flags().StringVarP(&filename, "filename", "f", "", "config file for tentez")
	if err := applyCmd.MarkFlagRequired("filename"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(applyCmd)
}
