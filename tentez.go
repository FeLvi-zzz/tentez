package tentez

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Tentez interface {
	Plan() error
	Apply(isForce bool) error
	Get() (map[TargetType]TargetsData, error)
	Rollback(hasPause bool) error
	Switch(weights []int, hasPause bool) error
}

type tentez struct {
	Targets map[TargetType]Targets
	Steps   []Step
	config  Config
	ui      Ui
}

var (
	defaultSteps = []Step{
		{
			Type: StepTypePause,
		},
		{
			Type: StepTypeSwitch,
			Weight: Weight{
				Old: 70,
				New: 30,
			},
		},
		{
			Type:         StepTypeSleep,
			SleepSeconds: 600,
		},
		{
			Type: StepTypePause,
		},
		{
			Type: StepTypeSwitch,
			Weight: Weight{
				Old: 30,
				New: 70,
			},
		},
		{
			Type:         StepTypeSleep,
			SleepSeconds: 600,
		},
		{
			Type: StepTypePause,
		},
		{
			Type: StepTypeSwitch,
			Weight: Weight{
				Old: 0,
				New: 100,
			},
		},
		{
			Type:         StepTypeSleep,
			SleepSeconds: 600,
		},
	}
)

func New(targets map[TargetType]Targets, steps []Step) (tentez, error) {
	config, err := NewConfig()
	if err != nil {
		return tentez{}, err
	}

	return tentez{
		Targets: targets,
		Steps:   steps,
		config:  config,
		ui: &cui{
			in:  os.Stdin,
			out: os.Stdout,
			err: os.Stderr,
		},
	}, nil
}

func (t tentez) Apply(isForce bool) (err error) {
	for i, step := range t.Steps {
		t.ui.Outputf("\n%d / %d steps\n", i+1, len(t.Steps))

		switch step.Type {
		case StepTypePause:
			t.pause()
		case StepTypeSleep:
			t.sleep(step.SleepSeconds)
		case StepTypeSwitch:
			err = t.execSwitch(step.Weight, isForce)
		default:
			return fmt.Errorf(`unknown step type "%s"`, step.Type)
		}

		if err != nil {
			return err
		}

		t.ui.Outputln("")
	}

	t.ui.Outputln("Apply complete!")

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

	t.ui.Outputf(output.String())

	return nil
}

func (t tentez) Get() (targetsMap map[TargetType]TargetsData, err error) {
	mapData := map[TargetType]TargetsData{}
	for targetType, targetResources := range t.Targets {
		data, err := targetResources.fetchData(t.config)
		if err != nil {
			if errors.Is(err, &FailedFetchTargetGroupsError{}) {
				t.ui.OutputErrln(err.Error())
			} else {
				return nil, err
			}
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

func (t tentez) Switch(weights []int, hasPause bool) (err error) {
	t.Steps = []Step{}

	if hasPause {
		t.Steps = append(t.Steps, Step{
			Type: StepTypePause,
		})
	}

	t.Steps = append(t.Steps, Step{
		Type: StepTypeSwitch,
		Weight: Weight{
			Old: int32(weights[0]),
			New: int32(weights[1]),
		},
	})

	if err = t.Plan(); err != nil {
		return err
	}
	return t.Apply(false)
}
