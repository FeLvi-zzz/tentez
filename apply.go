package tentez

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func pause(cfg Config) {
	fmt.Fprintln(cfg.io.out, "Pause")
	fmt.Fprintln(cfg.io.out, `enter "yes", continue steps.`)
	fmt.Fprintln(cfg.io.out, `If you'd like to interrupt steps, enter "quit".`)

	for {
		input := ""
		fmt.Fprint(cfg.io.out, "> ")
		fmt.Fscan(cfg.io.in, &input)

		if input == "yes" {
			fmt.Fprintln(cfg.io.out, "continue step")
			break
		} else if input == "quit" {
			fmt.Fprintln(cfg.io.out, "Bye")
			os.Exit(0)
		}
	}
}

func sleep(sec int, cfg Config) {
	seconds := time.Duration(sec) * time.Second
	finishAt := time.Now().Add(seconds)

	fmt.Fprintf(cfg.io.out, "Sleep %ds\n", sec)
	fmt.Fprintf(cfg.io.out, "Resume at %s\n", finishAt.Format("2006-01-02 15:04:05"))

	var wg sync.WaitGroup

	ticker := time.NewTicker(1 * time.Second)
	tickerStop := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case t := <-ticker.C:
				fmt.Fprintf(cfg.io.out, "\rRemain: %ds ", int(finishAt.Sub(t).Seconds()))

			case <-tickerStop:
				fmt.Fprintln(cfg.io.out, "\a")
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(seconds)
		ticker.Stop()
		tickerStop <- true
	}()

	wg.Wait()

	fmt.Fprintln(cfg.io.out, "Resume")
}

func execSwitch(targets map[string]Targets, weight Weight, isForce bool, cfg Config) error {
	fmt.Fprintf(cfg.io.out, "Switch old:new = %d:%d\n", weight.Old, weight.New)

	i := 0
	for _, targetResouces := range targets {
		for _, target := range targetResouces.targetsSlice() {
			i++

			fmt.Fprintf(cfg.io.out, "%d. %s ", i, target.getName())
			if err := target.execSwitch(weight, false, cfg); err != nil {
				_, ok := err.(SkipSwitchError)
				if !ok {
					return err
				}
				fmt.Fprintln(cfg.io.out, err.Error())
			} else {
				fmt.Fprintln(cfg.io.out, "switched!")
			}
		}
	}

	fmt.Fprintf(cfg.io.out, "Switched at %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}
