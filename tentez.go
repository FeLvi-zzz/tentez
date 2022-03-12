package tentez

import (
	"fmt"
)

type Tentez interface {
	Plan() error
	Apply() error
	Get() error
}

type tentez struct {
	Targets map[string]Targets
	Steps   []Step
	config  Config
}

func (t tentez) Apply() (err error) {
	for i, step := range t.Steps {
		fmt.Fprintf(t.config.io.out, "\n%d / %d steps\n", i+1, len(t.Steps))

		switch step.Type {
		case "pause":
			pause(t.config)
		case "sleep":
			sleep(step.SleepSeconds, t.config)
		case "switch":
			err = execSwitch(t.Targets, step.Weight, t.config)
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
	fmt.Fprintln(t.config.io.out, "Plan:")
	targetNames := getTargetNames(t.Targets)

	for i, step := range t.Steps {
		fmt.Fprintf(t.config.io.out, "%d. ", i+1)

		switch step.Type {
		case "pause":
			fmt.Fprintln(t.config.io.out, "pause")

		case "switch":
			weight := step.Weight
			fmt.Fprintf(t.config.io.out, "switch old:new = %d:%d\n", weight.Old, weight.New)
			for _, name := range targetNames {
				fmt.Fprintf(t.config.io.out, "  - %s\n", name)
			}

		case "sleep":
			fmt.Fprintf(t.config.io.out, "sleep %ds\n", step.SleepSeconds)

		default:
			return fmt.Errorf(`unknown step type "%s"`, step.Type)
		}
	}

	return nil
}

func (t tentez) Get() (err error) {
	for _, targetResouces := range t.Targets {
		if err = outputData(targetResouces, t.config); err != nil {
			return err
		}
	}

	return
}
