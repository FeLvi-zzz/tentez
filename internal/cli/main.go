package cli

import (
	"flag"
	"fmt"

	"github.com/FeLvi-zzz/tentez"
)

func flagParse() (cmd string, filepath string) {
	flag.Usage = help

	flag.StringVar(&filepath, "f", "", "filepath")

	flag.Parse()

	cmd = flag.Arg(0)

	return
}

func help() {
	helpText := `
Usage:
  tentez [-f <filename>] <subcommand>

Commands:
  plan     Show steps how to apply
  apply    Switch targets weights
  get      Show current state of targets
  rollback Rollback switch, switch old:new = 100:0
  help     Show this help
	version  Show version

Flags:
  -f <filename>  Specify YAML file
  -h             Show this help
`
	fmt.Println(helpText)
}

func Run() error {
	cmd, filepath := flagParse()

	var (
		t   tentez.Tentez
		err error
	)

	switch cmd {
	case "plan", "apply", "get", "rollback":
		t, err = tentez.NewFromYaml(filepath)
		if err != nil {
			return err
		}
	}

	switch cmd {
	case "plan":
		return t.Plan()
	case "apply":
		if err := t.Plan(); err != nil {
			return err
		}
		return t.Apply()
	case "rollback":
		return t.Rollback()
	case "get":
		return t.Get()
	case "version":
		fmt.Printf("tentez version: %s (rev: %s)\n", tentez.Version, tentez.Revision)
		return nil
	case "help", "":
		help()
		return nil
	default:
		help()
		return fmt.Errorf(`unknown command "%s"`, cmd)
	}
}
