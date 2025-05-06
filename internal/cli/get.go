package cli

import (
	"fmt"

	"github.com/FeLvi-zzz/tentez"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Show current state of targets",
	Long: `
# Show current state of targets
$ tentez -f ./examples/example.yaml get`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		t, err := tentez.NewFromYaml(ctx, filename)
		if err != nil {
			return err
		}

		dataMap, err := t.Get(ctx)
		if err != nil {
			return err
		}
		output, err := yaml.Marshal(&dataMap)
		if err != nil {
			return err
		}
		fmt.Print(string(output))

		return nil
	},
}

func init() {
	getCmd.Flags().StringVarP(&filename, "filename", "f", "", "config file for tentez")
	if err := getCmd.MarkFlagRequired("filename"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(getCmd)
}
