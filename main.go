package tentez

import (
	"flag"
	"fmt"
)

func help() {
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

func flagParse() (cmd string, filepath string, err error) {
	flag.Usage = help

	filepath = *flag.String("f", "", "filepath")

	flag.Parse()

	cmd = flag.Arg(0)

	if filepath == "" {
		err = fmt.Errorf("filepath(-f option) must be set")
	}

	return
}

func Run() error {
	cmd, filepath, err := flagParse()
	if err != nil {
		return err
	}

	steps, targets, err := loadYaml(filepath)
	if err != nil {
		return err
	}

	t := tentez{
		Steps:   steps,
		Targets: targets,
	}

	return t.Exec(cmd)
}
