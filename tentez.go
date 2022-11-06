package tentez

import (
	"fmt"
	"strings"
)

type Tentez interface {
	Plan() error
	Apply(isForce bool) error
	Get() (map[TargetType]TargetsData, error)
	Rollback(hasPause bool) error
}

type tentez struct {
	Targets map[TargetType]Targets
	Steps   []Step
	config  Config
}

func New(targets map[TargetType]Targets, steps []Step) (tentez, error) {
	config, err := newConfig()
	if err != nil {
		return tentez{}, err
	}

	return tentez{
		Targets: targets,
		Steps:   steps,
		config:  config,
	}, nil
}

func (t tentez) Apply(isForce bool) (err error) {
	for i, step := range t.Steps {
		fmt.Fprintf(t.config.io.out, "\n%d / %d steps\n", i+1, len(t.Steps))

		switch step.Type {
		case StepTypePause:
			pause(t.config)
		case StepTypeSleep:
			sleep(step.SleepSeconds, t.config)
		case StepTypeSwitch:
			err = execSwitch(t.Targets, step.Weight, isForce, t.config)
		default:
			return fmt.Errorf(`unknown step type "%s"`, step.Type)
		}

		if err != nil {
			return err
		}

		fmt.Fprintln(t.config.io.out, "")
	}

	fmt.Fprintln(t.config.io.out, "Apply complete!")

	return nil
}

func (t tentez) Plan() error {
	var output strings.Builder

	fmt.Fprintln(&output, "Plan:")
	targetNames := getTargetNames(t.Targets)

	for i, step := range t.Steps {
		fmt.Fprintf(&output, "%d. ", i+1)

		switch step.Type {
		case StepTypePause:
			fmt.Fprintln(&output, "pause")

		case StepTypeSwitch:
			weight := step.Weight
			fmt.Fprintf(&output, "switch old:new = %d:%d\n", weight.Old, weight.New)
			for _, name := range targetNames {
				fmt.Fprintf(&output, "  - %s\n", name)
			}

		case StepTypeSleep:
			fmt.Fprintf(&output, "sleep %ds\n", step.SleepSeconds)

		default:
			return fmt.Errorf(`unknown step type "%s"`, step.Type)
		}
	}

	fmt.Fprint(t.config.io.out, output.String())

	return nil
}

func (t tentez) Get() (targetsMap map[TargetType]TargetsData, err error) {
	mapData := map[TargetType]TargetsData{}
	for targetType, targetResources := range t.Targets {
		data, err := targetResources.fetchData(t.config)
		if err != nil {
			return nil, err
		}
		if data == nil {
			continue
		}
		mapData[targetType] = data
	}

	return mapData, nil
}

func (t tentez) Rollback(hasPause bool) (err error) {
	t.Steps = []Step{}

	if hasPause {
		t.Steps = append(t.Steps, Step{
			Type: StepTypePause,
		})
	}

	t.Steps = append(t.Steps, Step{
		Type: StepTypeSwitch,
		Weight: Weight{
			Old: 100,
			New: 0,
		},
	})

	if err = t.Plan(); err != nil {
		return err
	}
	return t.Apply(true)
}
