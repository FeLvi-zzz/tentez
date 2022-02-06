package tentez

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func loadYaml(filepath *string) (yamlStruct *YamlStruct, err error) {
	data, err := ioutil.ReadFile(*filepath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &yamlStruct); err != nil {
		return nil, err
	}

	return yamlStruct, nil
}
