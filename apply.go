package tentez

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func pause() {
	fmt.Println("Pause")
	fmt.Println(`enter "yes", continue steps.`)
	fmt.Println(`If you'd like to interrupt steps, enter "quit".`)

	for {
		input := ""
		fmt.Print("> ")
		fmt.Scan(&input)

		if input == "yes" {
			fmt.Println("continue step")
			break
		} else if input == "quit" {
			fmt.Println("Bye")
			os.Exit(0)
		}
	}
}

func sleep(sec int) {
	seconds := time.Duration(sec) * time.Second
	finishAt := time.Now().Add(seconds)

	fmt.Printf("Sleep %ds\n", sec)
	fmt.Printf("Resume at %s\n", finishAt.Format("2006-01-02 15:04:05"))

	var wg sync.WaitGroup

	ticker := time.NewTicker(1 * time.Second)
	tickerStop := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case t := <-ticker.C:
				fmt.Printf("\rRemain: %ds ", int(finishAt.Sub(t).Seconds()))

			case <-tickerStop:
				fmt.Println("\a")
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

	fmt.Println("Resume")
}

func execSwitch(targets map[string]Targets, weight Weight) error {
	fmt.Printf("Switch old:new = %d:%d\n", weight.Old, weight.New)

	i := 0
	for _, targetResouces := range targets {
		for _, target := range targetResouces.(interface{}).([]Target) {
			i++

			fmt.Printf("%d. %s ", i, target.getName())
			if err := target.execSwitch(weight); err != nil {
				return err
			}
			fmt.Println("switched!")
		}
	}

	fmt.Printf("Switched at %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}
