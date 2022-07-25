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

	if !isForce {
		listenerData, err := cfg.client.elbv2.DescribeListeners(context.TODO(), &elbv2.DescribeListenersInput{
			ListenerArns: []string{l.Target},
		})
		if err != nil {
			return err
		}

		targetGroupWeightMap := map[string]float64{}
		for _, action := range listenerData.Listeners[0].DefaultActions {
			for _, tgTuple := range action.ForwardConfig.TargetGroups {
				targetGroupWeightMap[*tgTuple.TargetGroupArn] = float64(*tgTuple.Weight)
			}
		}
		weightSum := 0.0
		for _, v := range targetGroupWeightMap {
			weightSum += v
		}

		oldWeight, ok := targetGroupWeightMap[l.Switch.Old]
		if !ok || oldWeight/weightSum < float64(targetWeight.Old)/float64(targetWeight.New+targetWeight.Old) {
			return SkipSwitchError{"the old weight target is larger than current one."}
		}

		newWeight, ok := targetGroupWeightMap[l.Switch.New]
		if ok && newWeight/weightSum > float64(targetWeight.New)/float64(targetWeight.New+targetWeight.Old) {
			return SkipSwitchError{"the new weight target is smaller than current one."}
		}
	}

	_, err := cfg.client.elbv2.ModifyListener(context.TODO(), &elbv2.ModifyListenerInput{
		ListenerArn: aws.String(l.Target),
		DefaultActions: []elbv2Types.Action{
			{
				Type: "forward",
				ForwardConfig: &elbv2Types.ForwardActionConfig{
					TargetGroups: []elbv2Types.TargetGroupTuple{
						{
							TargetGroupArn: aws.String(l.Switch.Old),
							Weight:         aws.Int32(targetWeight.Old),
						},
						{
							TargetGroupArn: aws.String(l.Switch.New),
							Weight:         aws.Int32(targetWeight.New),
						},
					},
				},
			},
		},
	})

	return err
}

func (l AwsListener) getName() string {
	return l.Name
}

func (ls AwsListeners) fetchData(cfg Config) (interface{}, error) {
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

	if _, err := cfg.client.elbv2.DescribeTargetGroups(context.TODO(), &elbv2.DescribeTargetGroupsInput{
		TargetGroupArns: tgArns,
	}); err != nil {
		fmt.Fprintln(cfg.io.err, err.Error())
	}

	listenersData, err := cfg.client.elbv2.DescribeListeners(context.TODO(), &elbv2.DescribeListenersInput{
		ListenerArns: listenerArns,
	})
	if err != nil {
		return nil, err
	}

	res := struct {
		AwsListeners []AwsListenerData `yaml:"aws_listeners"`
	}{
		AwsListeners: []AwsListenerData{},
	}

	for _, listenerData := range listenersData.Listeners {
		for _, action := range listenerData.DefaultActions {
			targetGroupTuples := []AwsTargetGroupTuple{}
			for _, tgTuple := range action.ForwardConfig.TargetGroups {
				targetGroupTuples = append(targetGroupTuples, AwsTargetGroupTuple{
					Type:           listenerMap[*listenerData.ListenerArn].Switch.getType(*tgTuple.TargetGroupArn),
					TargetGroupArn: aws.ToString(tgTuple.TargetGroupArn),
					Weight:         aws.ToInt32(tgTuple.Weight),
				})
			}

			res.AwsListeners = append(res.AwsListeners, AwsListenerData{
				Name:       listenerMap[*listenerData.ListenerArn].Name,
				ListnerArn: *listenerData.ListenerArn,
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
