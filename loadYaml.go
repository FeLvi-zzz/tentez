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

	config, err := newConfig()
	if err != nil {
		return nil, err
	}

	return tentez{
		Steps: yamlStruct.Steps,
		Targets: map[string]Targets{
			"aws_listener_rules": yamlStruct.AwsListenerRules,
			"aws_listeners":      yamlStruct.AwsListeners,
		},
		config: config,
	}, nil
}
