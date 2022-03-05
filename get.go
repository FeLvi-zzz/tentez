package tentez

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func outputData(targets Targets, client Client) error {
	targetsData, err := targets.fetchData(client)
	if err != nil {
		return err
	}
	if targetsData == nil {
		return nil
	}

	output, err := yaml.Marshal(&targetsData)
	if err != nil {
		return err
	}

	fmt.Print(string(output))

	return nil
}
