package tentez

import (
	"context"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func NewFromYaml(ctx context.Context, filename string) (t Tentez, err error) {
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
		ctx,
		map[TargetType]Targets{
			TargetTypeAwsListenerRule: yamlStruct.AwsListenerRules,
			TargetTypeAwsListener:     yamlStruct.AwsListeners,
		},
		yamlStruct.Steps,
	)
}
