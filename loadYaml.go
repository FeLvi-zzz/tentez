package tentez

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func loadYaml(filepath string) (steps []Step, targets map[string]Targets, err error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}

	yamlStruct := YamlStruct{}

	if err = yaml.Unmarshal(data, &yamlStruct); err != nil {
		return
	}

	steps = yamlStruct.Steps

	targets = map[string]Targets{}

	targets["aws_listener_rules"] = yamlStruct.AwsListenerRules
	targets["aws_listeners"] = yamlStruct.AwsListeners

	return
}
