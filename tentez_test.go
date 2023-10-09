package tentez_test

import (
	"context"
	"fmt"

	"github.com/FeLvi-zzz/tentez"
	"gopkg.in/yaml.v3"
)

func Example() {
	t, err := tentez.New(
		context.TODO(),
		map[tentez.TargetType]tentez.Targets{
			tentez.TargetTypeAwsListenerRule: tentez.AwsListenerRules([]tentez.AwsListenerRule{}),
			tentez.TargetTypeAwsListener:     tentez.AwsListeners([]tentez.AwsListener{}),
		},
		[]tentez.Step{
			{
				Type: tentez.StepTypeSleep,
			},
		},
	)
	if err != nil {
		return
	}

	targetsData, err := t.Get(context.TODO())
	if err != nil {
		return
	}

	output, err := yaml.Marshal(&targetsData)
	if err != nil {
		return
	}
	fmt.Print(string(output))

	data, ok := targetsData[tentez.TargetTypeAwsListener].([]tentez.AwsListenerData)
	if ok && len(data) > 0 {
		name := data[0].Name
		fmt.Println(name)
	}
}
