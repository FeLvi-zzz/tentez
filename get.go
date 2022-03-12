package tentez

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func outputData(targets Targets, cfg Config) error {
	targetsData, err := targets.fetchData(cfg)
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

	fmt.Fprint(cfg.io.out, string(output))

	return nil
}
