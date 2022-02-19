package tentez

import (
	"flag"
	"fmt"
)

func Run() error {
	flag.Usage = func() {
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
		fmt.Println(helpText)
	}

	filepath := flag.String("f", "", "filepath")

	flag.Parse()

	cmd := flag.Arg(0)

	if *filepath == "" {
		return fmt.Errorf("filepath(-f option) must be set.")
	}

	yamlData, err := loadYaml(filepath)
	if err != nil {
		return err
	}

	switch cmd {
	case "plan":
		return Plan(yamlData)
	case "apply":
		if err := Plan(yamlData); err != nil {
			return err
		}
		return Apply(yamlData)
	case "get":
		return Get(yamlData)
	case "help", "":
		flag.Usage()
		return nil
	default:
		flag.Usage()
		return fmt.Errorf(`Error: unknown command "%s"`, cmd)
	}
}
