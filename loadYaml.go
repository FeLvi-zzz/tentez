package tentez

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func NewFromYaml(filepath string) (t Tentez, err error) {
	if filepath == "" {
		return nil, fmt.Errorf("filepath must be specified")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return
	}

	yamlStruct := YamlStruct{}
	if err = yaml.Unmarshal(data, &yamlStruct); err != nil {
		return
	}

	return New(
		map[TargetType]Targets{
			TargetTypeAwsListenerRule: yamlStruct.AwsListenerRules,
			TargetTypeAwsListener:     yamlStruct.AwsListeners,
		},
		yamlStruct.Steps,
	)
}
