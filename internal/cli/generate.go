package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/FeLvi-zzz/tentez"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var filepaths = []string{}
var output = ""

var generateConfigCmd = &cobra.Command{
	Use:   "generate-config",
	Short: "generate configs",
}

var generateConfigTerraformPlanJsonCmd = &cobra.Command{
	Use:   "tfplanjson",
	Short: "generate from terraform plan json",
	Long: `
# generate config from terraform plan json
$ terraform plan -out tfplan && terraform show -json tfplan > tfplan.json
$ tentez generate-config tfplanjson -f ./tfplan.json -o tentez.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tfplanjsons := []tentez.TerraformPlanJson{}
		for _, filepath := range filepaths {
			data, err := os.ReadFile(filepath)
			if err != nil {
				return fmt.Errorf("cannot read file: %w", err)
			}

			tfplanjson := tentez.TerraformPlanJson{}
			if err := json.Unmarshal(data, &tfplanjson); err != nil {
				return fmt.Errorf("cannot parse json: %w", err)
			}

			tfplanjsons = append(tfplanjsons, tfplanjson)
		}

		configYaml, err := tentez.GenerateConfigFromTerraformPlanJsons(tfplanjsons)
		if err != nil {
			return fmt.Errorf("cannot generate config: %w", err)
		}

		configYamlBytes, err := yaml.Marshal(configYaml)
		if err != nil {
			return fmt.Errorf("cannot marshal config yaml: %w", err)
		}

		if output == "" {
			output = "tentez.yaml"
		}

		file, err := os.Create(output)
		if err != nil {
			return err
		}

		if _, err := file.Write(configYamlBytes); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	generateConfigCmd.PersistentFlags().StringArrayVarP(&filepaths, "filepath", "f", []string{}, "terraform plan json file")
	if err := generateConfigCmd.MarkPersistentFlagRequired("filepath"); err != nil {
		panic(err)
	}

	generateConfigCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")

	generateConfigCmd.AddCommand(generateConfigTerraformPlanJsonCmd)

	rootCmd.AddCommand(generateConfigCmd)
}
