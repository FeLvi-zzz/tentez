package tentez

import (
	"errors"
	"os"
	"sync"
	"time"
)

func (t tentez) pause() {
	t.ui.Outputln("Pause")
	t.ui.Outputln(`enter "yes", continue steps.`)
	t.ui.Outputln(`If you'd like to interrupt steps, enter "quit".`)

	for {
		input := t.ui.Ask("> ")

		if input == "yes" {
			t.ui.Outputln("continue step")
			break
		} else if input == "quit" {
			t.ui.Outputln("Bye")
			os.Exit(0)
		}
	}
}

func (t tentez) sleep(sec int) {
	seconds := time.Duration(sec) * time.Second
	finishAt := time.Now().Add(seconds)

	t.ui.Outputf("Sleep %ds\n", sec)
	t.ui.Outputf("Resume at %s\n", finishAt.Format("2006-01-02 15:04:05"))

	var wg sync.WaitGroup

	ticker := time.NewTicker(1 * time.Second)
	tickerStop := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case tk := <-ticker.C:
				t.ui.Outputf("\rRemain: %ds ", int(finishAt.Sub(tk).Seconds()))

			case <-tickerStop:
				t.ui.Outputln("\a")
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		t.config.clock.Sleep(seconds)
		ticker.Stop()
		tickerStop <- true
	}()

	wg.Wait()

	t.ui.Outputln("Resume")
}

func (t tentez) execSwitch(weight Weight, isForce bool) error {
	t.ui.Outputf("Switch old:new = %d:%d\n", weight.Old, weight.New)

	i := 0
	for _, targetResouces := range t.Targets {
		for _, target := range targetResouces.targetsSlice() {
			i++

			t.ui.Outputf("%d. %s ", i, target.getName())
			if err := target.execSwitch(weight, isForce, t.config); err != nil {
				if !errors.As(err, &SkipSwitchError{}) {
					return err
				}
				t.ui.Outputln(err.Error())
			} else {
				t.ui.Outputln("switched!")
			}
		}
	}

	t.ui.Outputf("Switched at %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}
