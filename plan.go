package tentez

import (
	"fmt"
	"strings"
)

func GetTargetNames(yamlData *YamlStruct) ([]string, error) {
	targetNames := []string{}

	for _, rule := range yamlData.AwsListenerRules {
		errAttributes := []string{}
		if rule.Name == "" {
			errAttributes = append(errAttributes, "name")
		}
		if rule.Target == "" {
			errAttributes = append(errAttributes, "target")
		}
		if rule.Switch.New == "" {
			errAttributes = append(errAttributes, "switch.new")
		}
		if rule.Switch.Old == "" {
			errAttributes = append(errAttributes, "switch.old")
		}

		if len(errAttributes) > 0 {
			return nil, fmt.Errorf("Error: %s must be set.", strings.Join(errAttributes, ", "))
		}

		targetNames = append(targetNames, rule.Name)
	}

	for _, rule := range yamlData.AwsListeners {
		errAttributes := []string{}
		if rule.Name == "" {
			errAttributes = append(errAttributes, "name")
		}
		if rule.Target == "" {
			errAttributes = append(errAttributes, "target")
		}
		if rule.Switch.New == "" {
			errAttributes = append(errAttributes, "switch.new")
		}
		if rule.Switch.Old == "" {
			errAttributes = append(errAttributes, "switch.old")
		}

		if len(errAttributes) > 0 {
			return nil, fmt.Errorf("Error: %s must be set.", strings.Join(errAttributes, ", "))
		}

		targetNames = append(targetNames, rule.Name)
	}

	return targetNames, nil
}

func Plan(yamlData *YamlStruct) error {
	fmt.Println("Plan:")
	targetNames, err := GetTargetNames(yamlData)
	if err != nil {
		return err
	}

	for i, step := range yamlData.Steps {
		fmt.Printf("%d. ", i+1)

		switch step.Type {
		case "pause":
			fmt.Println("pause")

		case "switch":
			weight := step.Weight
			fmt.Printf("switch old:new = %d:%d\n", weight.Old, weight.New)
			for _, name := range targetNames {
				fmt.Printf("  - %s\n", name)
			}

		case "sleep":
			sleepSec := step.SleepSeconds
			fmt.Printf("sleep %ds\n", sleepSec)

		default:
			return fmt.Errorf(`Error: unknown type "%s"`, step.Type)
		}
	}

	return nil
}
