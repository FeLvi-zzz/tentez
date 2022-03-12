package main

import (
	"fmt"
	"os"

	"github.com/FeLvi-zzz/tentez/internal/cli"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	err := cli.Run()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}
