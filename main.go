package tentez

import (
	"flag"
	"fmt"
	"log"
)

func Run() error {
	filepath := flag.String("f", "", "filepath")

	flag.Parse()

	cmd := flag.Arg(0)

	if *filepath == "" {
		log.Fatalf("filepath(-f option) must be set.")
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
	default:
		return fmt.Errorf(`Error: unknown command "%s"`, cmd)
	}
}
