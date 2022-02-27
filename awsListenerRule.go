package tentez

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type AwsListenerRule struct {
	Name   string `yaml:"name"`
	Target string `yaml:"target"`
	Switch Switch `yaml:"switch"`
}

type AwsListenerRules []AwsListenerRule

type AwsListenerRuleData struct {
	ListnerRuleArn string                `yaml:"target"`
	Weights        []AwsTargetGroupTuple `yaml:"weights"`
}

func (r AwsListenerRule) execSwitch(weight Weight) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	elbv2svc := elbv2.NewFromConfig(cfg)

	if _, err := elbv2svc.ModifyRule(context.TODO(), &elbv2.ModifyRuleInput{
		RuleArn: aws.String(r.Target),
		Actions: []elbv2Types.Action{
			{
				Type: "forward",
				ForwardConfig: &elbv2Types.ForwardActionConfig{
					TargetGroups: []elbv2Types.TargetGroupTuple{
						{
							TargetGroupArn: aws.String(r.Switch.Old),
							Weight:         aws.Int32(weight.Old),
						},
						{
							TargetGroupArn: aws.String(r.Switch.New),
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

func (r AwsListenerRule) getName() string {
	return r.Name
}

func (rs AwsListenerRules) fetchData() (interface{}, error) {
	if len(rs) == 0 {
		return nil, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	ruleArns := []string{}
	for _, rule := range rs {
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

func (rs AwsListenerRules) targetsSlice() (targets []Target) {
	for _, target := range rs {
		targets = append(targets, target)
	}
	return targets
}
