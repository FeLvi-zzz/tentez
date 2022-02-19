package tentez

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
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

func (rule *AwsListenerRule) ExecSwitch(weight Weight) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	elbv2svc := elbv2.NewFromConfig(cfg)

	if _, err := elbv2svc.ModifyRule(context.TODO(), &elbv2.ModifyRuleInput{
		RuleArn: aws.String(rule.Target),
		Actions: []elbv2Types.Action{
			{
				Type: "forward",
				ForwardConfig: &elbv2Types.ForwardActionConfig{
					TargetGroups: []elbv2Types.TargetGroupTuple{
						{
							TargetGroupArn: aws.String(rule.Switch.Old),
							Weight:         aws.Int32(weight.Old),
						},
						{
							TargetGroupArn: aws.String(rule.Switch.New),
							Weight:         aws.Int32(weight.New),
						},
					},
				},
			},
		},
	}); err != nil {
		return err
	}

	return nil
}
func (listener *AwsListener) ExecSwitch(weight Weight) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	elbv2svc := elbv2.NewFromConfig(cfg)

	// avoid rate limit
	time.Sleep(1 * time.Second)

	if _, err := elbv2svc.ModifyListener(context.TODO(), &elbv2.ModifyListenerInput{
		ListenerArn: aws.String(listener.Target),
		DefaultActions: []elbv2Types.Action{
			{
				Type: "forward",
				ForwardConfig: &elbv2Types.ForwardActionConfig{
					TargetGroups: []elbv2Types.TargetGroupTuple{
						{
							TargetGroupArn: aws.String(listener.Switch.Old),
							Weight:         aws.Int32(weight.Old),
						},
						{
							TargetGroupArn: aws.String(listener.Switch.New),
							Weight:         aws.Int32(weight.New),
						},
					},
				},
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func execSwitch(yamlData *YamlStruct, weight Weight) error {
	fmt.Printf("Switch old:new = %d:%d\n", weight.Old, weight.New)

	i := 0
	for _, rule := range yamlData.AwsListenerRules {
		i++

		fmt.Printf("%d. %s ", i, rule.Name)
		if err := rule.ExecSwitch(weight); err != nil {
			return err
		}
		fmt.Println("switched!")
	}

	for _, rule := range yamlData.AwsListeners {
		i++

		fmt.Printf("%d. %s ", i, rule.Name)
		if err := rule.ExecSwitch(weight); err != nil {
			return err
		}
		fmt.Println("switched!")
	}

	fmt.Printf("Switched at %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

func Apply(yamlData *YamlStruct) (err error) {
	for i, step := range yamlData.Steps {
		fmt.Printf("\n%d / %d steps\n", i+1, len(yamlData.Steps))

		switch step.Type {
		case "pause":
			pause()
		case "sleep":
			sleep(step.SleepSeconds)
		case "switch":
			err = execSwitch(yamlData, step.Weight)
		default:
			return fmt.Errorf(`Error: unknown step type "%s"`, step.Type)
		}

		if err != nil {
			return err
		}

		fmt.Println("")
	}

	fmt.Println("Apply complete!")

	return nil
}
