package tentez_test

import (
	"github.com/FeLvi-zzz/tentez"
)

func Example() {
	t, err := tentez.New(
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

	t.Get()
}
