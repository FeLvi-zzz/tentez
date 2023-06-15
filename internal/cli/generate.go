package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/FeLvi-zzz/tentez"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var filenames = []string{}
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
		extendedFilenames := []string{}

		for _, filename := range filenames {
			fileInfo, err := os.Stat(filename)
			if err != nil {
				return fmt.Errorf("cannot read file: %w", err)
			}

			if !fileInfo.IsDir() {
				extendedFilenames = append(extendedFilenames, filename)
				continue
			}

			dirEntries, err := os.ReadDir(filename)
			if err != nil {
				return fmt.Errorf("cannot read dir: %w", err)
			}
			for _, e := range dirEntries {
				extendedFilenames = append(extendedFilenames, filepath.Join(filename, e.Name()))
			}
		}

		for _, filename := range extendedFilenames {
			data, err := os.ReadFile(filename)
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

var generateConfigResourceTagCmd = &cobra.Command{
	Use:   "resource-tag",
	Short: "generate from resource tags",
	Long: `
# generate config from resource tags
# refer https://github.com/FeLvi-zzz/tentez/blob/master/examples/tentez.ResourceTag.v1beta1.yaml.
$ tentez generate-config resource-tag -f examples/tentez.ResourceTag.v1beta1.yaml -o tentez.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("cannot read file: %w", err)
		}

		generateConfigResourceTag := tentez.GenerateConfigResourceTag{}
		if err := yaml.Unmarshal(data, &generateConfigResourceTag); err != nil {
			return fmt.Errorf("cannot parse yaml: %w", err)
		}

		var configYaml tentez.YamlStruct

		switch generateConfigResourceTag.Version {
		case "tentez.ResourceTag.v1beta1":
			config := tentez.GenerateConfigResourceTagV1beta1{}
			if err := yaml.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("cannot parse yaml: %w", err)
			}
			cfg, err := tentez.NewConfig()
			if err != nil {
				return fmt.Errorf("cannot create config: %w", err)
			}
			spec := config.Spec
			configYaml, err = tentez.GenerateConfigFromResourceTags(
				spec.FilterTags,
				spec.MatchingTagKeys,
				spec.SwitchTag.Key,
				spec.SwitchTag.Value.Old,
				spec.SwitchTag.Value.New,
				cfg,
			)
			if err != nil {
				return fmt.Errorf("cannot generate config: %w", err)
			}
		default:
			return fmt.Errorf("unknown version: %s", generateConfigResourceTag.Version)
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
	generateConfigCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output path")

	generateConfigTerraformPlanJsonCmd.Flags().StringArrayVarP(&filenames, "filename", "f", []string{}, "terraform plan json file")
	if err := generateConfigTerraformPlanJsonCmd.MarkFlagRequired("filename"); err != nil {
		panic(err)
	}
	generateConfigCmd.AddCommand(generateConfigTerraformPlanJsonCmd)

	generateConfigResourceTagCmd.Flags().StringVarP(&filename, "filename", "f", "", "resource tags yaml file")
	if err := generateConfigResourceTagCmd.MarkFlagRequired("filename"); err != nil {
		panic(err)
	}
	generateConfigCmd.AddCommand(generateConfigResourceTagCmd)

	rootCmd.AddCommand(generateConfigCmd)
}
