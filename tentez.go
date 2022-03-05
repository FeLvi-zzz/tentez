package tentez

import (
	"flag"
	"fmt"
)

type Tentez interface {
	plan() error
	apply() error
	get() error
	help()
}

type tentez struct {
	Targets map[string]Targets
	Steps   []Step
	config  Config
}

func Exec(t Tentez, cmd string) error {
	switch cmd {
	case "plan":
		return t.plan()
	case "apply":
		if err := t.plan(); err != nil {
			return err
		}
		return t.apply()
	case "get":
		return t.get()
	case "help", "":
		t.help()
		return nil
	default:
		t.help()
		return fmt.Errorf(`unknown command "%s"`, cmd)
	}
}

func (t tentez) apply() (err error) {
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

func (t tentez) plan() error {
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

func (t tentez) get() (err error) {
	for _, targetResouces := range t.Targets {
		if err = outputData(targetResouces, t.config); err != nil {
			return err
		}
	}

	return
}

func (t tentez) help() {
	flag.Usage = t.help
	helpText := `
Usage:
  tentez -f <filename> <subcommand>

Commands:
  plan   Show steps how to apply
  apply  Switch targets weights
  get    Show current state of targets
  help   Show this help

Flags:
  -f <filename>  Specify YAML file
  -h             Show this help
`
	fmt.Fprintln(t.config.io.out, helpText)
}
