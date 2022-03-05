package tentez

import (
	"flag"
	"fmt"
)

func flagParse() (cmd string, filepath string, err error) {
	flag.StringVar(&filepath, "f", "", "filepath")

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

	config, err := newConfig()
	if err != nil {
		return err
	}

	t := tentez{
		Steps:   steps,
		Targets: targets,
		config:  config,
	}

	return Exec(t, cmd)
}
