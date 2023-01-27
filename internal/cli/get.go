package cli

import (
	"fmt"

	"github.com/FeLvi-zzz/tentez"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Show current state of targets",
	Long: `
# Show current state of targets
$ tentez -f ./examples/example.yaml get`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := cmd.Flag("filename").Value.String()
		t, err := tentez.NewFromYaml(filename)
		if err != nil {
			return err
		}

		dataMap, err := t.Get()
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
