package tentez

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"gopkg.in/yaml.v2"
)

func (rules AwsListenerRules) FetchData() (interface{}, error) {
	if len(rules) == 0 {
		return nil, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	ruleArns := []string{}
	for _, rule := range rules {
		ruleArns = append(ruleArns, rule.Target)
	}

	elbv2svc := elbv2.NewFromConfig(cfg)

	rulesData, err := elbv2svc.DescribeRules(context.TODO(), &elbv2.DescribeRulesInput{
		RuleArns: ruleArns,
	})
	if err != nil {
		return nil, err
	}

	res := struct {
		AwsListenerRules []AwsListenerRuleData `yaml:"aws_listener_rules"`
	}{
		AwsListenerRules: []AwsListenerRuleData{},
	}

	for _, ruleData := range rulesData.Rules {
		for _, action := range ruleData.Actions {
			targetGroupTuples := []AwsTargetGroupTuple{}
			for _, tgTuple := range action.ForwardConfig.TargetGroups {
				targetGroupTuples = append(targetGroupTuples, AwsTargetGroupTuple{
					TargetGroupArn: aws.ToString(tgTuple.TargetGroupArn),
					Weight:         aws.ToInt32(tgTuple.Weight),
				})
			}

			res.AwsListenerRules = append(res.AwsListenerRules, AwsListenerRuleData{
				ListnerRuleArn: *ruleData.RuleArn,
				Weights:        targetGroupTuples,
			})
		}
	}

	return res, nil
}

func (listeners AwsListeners) FetchData() (interface{}, error) {
	if len(listeners) == 0 {
		return nil, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	listenerArns := []string{}
	for _, listener := range listeners {
		listenerArns = append(listenerArns, listener.Target)
	}

	elbv2svc := elbv2.NewFromConfig(cfg)

	listenersData, err := elbv2svc.DescribeListeners(context.TODO(), &elbv2.DescribeListenersInput{
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
					TargetGroupArn: aws.ToString(tgTuple.TargetGroupArn),
					Weight:         aws.ToInt32(tgTuple.Weight),
				})
			}

			res.AwsListeners = append(res.AwsListeners, AwsListenerData{
				ListnerArn: *listenerData.ListenerArn,
				Weights:    targetGroupTuples,
			})
		}
	}

	return res, nil
}

func outputData(targets GetTargets) error {
	targetsData, err := targets.FetchData()
	if err != nil {
		return err
	}
	if targetsData == nil {
		return nil
	}

	output, err := yaml.Marshal(&targetsData)
	if err != nil {
		return err
	}

	fmt.Print(string(output))

	return nil
}

func Get(yamlData *YamlStruct) (err error) {
	if err = outputData(yamlData.AwsListenerRules); err != nil {
		return err
	}

	if err = outputData(yamlData.AwsListeners); err != nil {
		return err
	}

	return
}
