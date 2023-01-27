package tentez

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func NewFromYaml(filename string) (t Tentez, err error) {
	if filename == "" {
		return nil, fmt.Errorf("filename must be specified")
	}

	data, err := os.ReadFile(filename)
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
