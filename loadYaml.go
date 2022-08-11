package tentez

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func NewFromYaml(filepath string) (t Tentez, err error) {
	if filepath == "" {
		return nil, fmt.Errorf("filepath(-f option) must be set")
	}

	data, err := ioutil.ReadFile(filepath)
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
