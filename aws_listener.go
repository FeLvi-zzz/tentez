package tentez

import (
	"context"
	"fmt"

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
	Name              string                      `yaml:"name"`
	ListnerArn        string                      `yaml:"target"`
	Weights           []AwsTargetGroupTuple       `yaml:"weights"`
	AdditionalActions []elbv2Types.ActionTypeEnum `yaml:"additional_actions,omitempty"`
}

func (l AwsListener) execSwitch(ctx context.Context, targetWeight Weight, isForce bool, cfg Config) error {
	listenerData, err := cfg.client.elbv2.DescribeListeners(ctx, &elbv2.DescribeListenersInput{
		ListenerArns: []string{l.Target},
	})
	if err != nil {
		return err
	}

	tgWeight := Weight{}
	for _, action := range listenerData.Listeners[0].DefaultActions {
		if action.Type != elbv2Types.ActionTypeEnumForward {
			continue
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
	if tgWeight.Old == 0 && tgWeight.New == 0 {
		return fmt.Errorf("%s does not have forward action", l.Target)
	}

	if !isForce {
		if tgWeight.CalcOldRatio() < targetWeight.CalcOldRatio() {
			return SkipSwitchError{"the old weight target is larger than current one."}
		}

		if tgWeight.CalcNewRatio() > targetWeight.CalcNewRatio() {
			return SkipSwitchError{"the new weight target is smaller than current one."}
		}
	}

	_, err = cfg.client.elbv2.ModifyListener(ctx, &elbv2.ModifyListenerInput{
		ListenerArn:    aws.String(l.Target),
		DefaultActions: makeNewActions(listenerData.Listeners[0].DefaultActions, l.Switch, targetWeight),
	})

	return err
}

func (l AwsListener) getName() string {
	return l.Name
}

func (ls AwsListeners) fetchData(ctx context.Context, cfg Config) (TargetsData, error) {
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

	errTgs, err := checkTargetGroupsExistense(ctx, cfg.client.elbv2, tgArns)
	if err != nil {
		return nil, err
	}

	listeners := []elbv2Types.Listener{}
	for _, listenerArnsBatch := range chunk(listenerArns, maxDescribeListenersItems) {
		listenersOutput, err := cfg.client.elbv2.DescribeListeners(ctx, &elbv2.DescribeListenersInput{
			ListenerArns: listenerArnsBatch,
		})
		if err != nil {
			return nil, err
		}
		listeners = append(listeners, listenersOutput.Listeners...)
	}

	res := []AwsListenerData{}

	for _, listener := range listeners {
		targetGroupTuples := []AwsTargetGroupTuple{}
		additionalActions := make([]elbv2Types.ActionTypeEnum, 0, len(listener.DefaultActions))

		for _, action := range listener.DefaultActions {
			if action.Type != elbv2Types.ActionTypeEnumForward {
				additionalActions = append(additionalActions, action.Type)
				continue
			}

			for _, tgTuple := range action.ForwardConfig.TargetGroups {
				targetGroupTuples = append(targetGroupTuples, AwsTargetGroupTuple{
					Type:           listenerMap[aws.ToString(listener.ListenerArn)].Switch.getType(aws.ToString(tgTuple.TargetGroupArn)),
					TargetGroupArn: aws.ToString(tgTuple.TargetGroupArn),
					Weight:         aws.ToInt32(tgTuple.Weight),
				})
			}
		}

		res = append(res, AwsListenerData{
			Name:              listenerMap[aws.ToString(listener.ListenerArn)].Name,
			ListnerArn:        aws.ToString(listener.ListenerArn),
			Weights:           targetGroupTuples,
			AdditionalActions: additionalActions,
		})
	}

	return res, NewFailedFetchTargetGroupsError(errTgs)
}

func (ls AwsListeners) targetsSlice() (targets []Target) {
	for _, target := range ls {
		targets = append(targets, target)
	}
	return targets
}
