package tentez

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type AwsListener struct {
	Name   string `yaml:"name"`
	Target string `yaml:"target"`
	Switch Switch `yaml:"switch"`
}

type AwsListeners []AwsListener

type AwsListenerData struct {
	Name       string                `yaml:"name"`
	ListnerArn string                `yaml:"target"`
	Weights    []AwsTargetGroupTuple `yaml:"weights"`
}

func (l AwsListener) execSwitch(targetWeight Weight, isForce bool, cfg Config) error {
	// avoid rate limit
	time.Sleep(1 * time.Second)

	listenerData, err := cfg.client.elbv2.DescribeListeners(context.TODO(), &elbv2.DescribeListenersInput{
		ListenerArns: []string{l.Target},
	})
	if err != nil {
		return err
	}

	tgWeight := Weight{}
	for _, action := range listenerData.Listeners[0].DefaultActions {
		if action.Type != elbv2Types.ActionTypeEnumForward {
			return fmt.Errorf("invalid action type: %s", action.Type)
		}

		for _, tgTuple := range action.ForwardConfig.TargetGroups {
			switch *tgTuple.TargetGroupArn {
			case l.Switch.Old:
				tgWeight.Old = *tgTuple.Weight
			case l.Switch.New:
				tgWeight.New = *tgTuple.Weight
			}
		}
	}

	if !isForce {
		if tgWeight.CalcOldRatio() < targetWeight.CalcOldRatio() {
			return SkipSwitchError{"the old weight target is larger than current one."}
		}

		if tgWeight.CalcNewRatio() > targetWeight.CalcNewRatio() {
			return SkipSwitchError{"the new weight target is smaller than current one."}
		}
	}

	_, err = cfg.client.elbv2.ModifyListener(context.TODO(), &elbv2.ModifyListenerInput{
		ListenerArn:    aws.String(l.Target),
		DefaultActions: compactActions(l.Switch, targetWeight),
	})

	return err
}

func (l AwsListener) getName() string {
	return l.Name
}

func (ls AwsListeners) fetchData(cfg Config) (TargetsData, error) {
	if len(ls) == 0 {
		return nil, nil
	}

	listenerArns := []string{}
	tgArns := []string{}
	listenerMap := map[string]AwsListener{}
	for _, listener := range ls {
		listenerArns = append(listenerArns, listener.Target)
		tgArns = append(tgArns, listener.Switch.New, listener.Switch.Old)
		listenerMap[listener.Target] = listener
	}

	for _, tgArnsBatch := range chunk(tgArns, maxDescribeTargetGroupsItems) {
		if _, err := cfg.client.elbv2.DescribeTargetGroups(context.TODO(), &elbv2.DescribeTargetGroupsInput{
			TargetGroupArns: tgArnsBatch,
		}); err != nil {
			fmt.Fprintln(cfg.io.err, err.Error())
		}
	}

	listeners := []elbv2Types.Listener{}
	for _, listenerArnsBatch := range chunk(listenerArns, maxDescribeListenersItems) {
		listenersOutput, err := cfg.client.elbv2.DescribeListeners(context.TODO(), &elbv2.DescribeListenersInput{
			ListenerArns: listenerArnsBatch,
		})
		if err != nil {
			return nil, err
		}
		listeners = append(listeners, listenersOutput.Listeners...)
	}

	res := []AwsListenerData{}

	for _, listener := range listeners {
		for _, action := range listener.DefaultActions {
			targetGroupTuples := []AwsTargetGroupTuple{}

			if action.Type == elbv2Types.ActionTypeEnumForward {
				for _, tgTuple := range action.ForwardConfig.TargetGroups {
					targetGroupTuples = append(targetGroupTuples, AwsTargetGroupTuple{
						Type:           listenerMap[aws.ToString(listener.ListenerArn)].Switch.getType(aws.ToString(tgTuple.TargetGroupArn)),
						TargetGroupArn: aws.ToString(tgTuple.TargetGroupArn),
						Weight:         aws.ToInt32(tgTuple.Weight),
					})
				}
			}

			res = append(res, AwsListenerData{
				Name:       listenerMap[aws.ToString(listener.ListenerArn)].Name,
				ListnerArn: aws.ToString(listener.ListenerArn),
				Weights:    targetGroupTuples,
			})
		}
	}

	return res, nil
}

func (ls AwsListeners) targetsSlice() (targets []Target) {
	for _, target := range ls {
		targets = append(targets, target)
	}
	return targets
}

func (l AwsListener) checkHealth(weight Weight, cfg Config) error {
	return checkSwitchTargetGroupHealth(l.Switch, weight, cfg)
}
