package tentez

import "fmt"

type tentez struct {
	Targets map[string]Targets
	Steps   []Step
}

func (t tentez) Exec(cmd string) error {
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
		help()
		return nil
	default:
		help()
		return fmt.Errorf(`unknown command "%s"`, cmd)
	}
}

func (t tentez) apply() (err error) {
	for i, step := range t.Steps {
		fmt.Printf("\n%d / %d steps\n", i+1, len(t.Steps))

		switch step.Type {
		case "pause":
			pause()
		case "sleep":
			sleep(step.SleepSeconds)
		case "switch":
			err = execSwitch(t.Targets, step.Weight)
		default:
			return fmt.Errorf(`unknown step type "%s"`, step.Type)
		}

		if err != nil {
			return err
		}

		fmt.Println("")
	}

	fmt.Println("Apply complete!")

	return nil
}

func (t tentez) plan() error {
	fmt.Println("Plan:")
	targetNames := getTargetNames(t.Targets)

	for i, step := range t.Steps {
		fmt.Printf("%d. ", i+1)

		switch step.Type {
		case "pause":
			fmt.Println("pause")

		case "switch":
			weight := step.Weight
			fmt.Printf("switch old:new = %d:%d\n", weight.Old, weight.New)
			for _, name := range targetNames {
				fmt.Printf("  - %s\n", name)
			}

		case "sleep":
			fmt.Printf("sleep %ds\n", step.SleepSeconds)

		default:
			return fmt.Errorf(`unknown step type "%s"`, step.Type)
		}
	}

	return nil
}

func (t tentez) get() (err error) {
	for _, targetResouces := range t.Targets {
		if err = outputData(targetResouces); err != nil {
			return err
		}
	}

	return
}
