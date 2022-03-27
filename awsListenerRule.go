package tentez

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	Name           string                `yaml:"name"`
	ListnerRuleArn string                `yaml:"target"`
	Weights        []AwsTargetGroupTuple `yaml:"weights"`
}

func (r AwsListenerRule) execSwitch(weight Weight, cfg Config) error {
	// avoid rate limit
	time.Sleep(1 * time.Second)

	_, err := cfg.client.elbv2.ModifyRule(context.TODO(), &elbv2.ModifyRuleInput{
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
	})

	return err
}

func (r AwsListenerRule) getName() string {
	return r.Name
}

func (rs AwsListenerRules) fetchData(cfg Config) (interface{}, error) {
	if len(rs) == 0 {
		return nil, nil
	}

	ruleArns := []string{}
	tgArns := []string{}
	ruleMap := map[string]AwsListenerRule{}
	for _, rule := range rs {
		ruleArns = append(ruleArns, rule.Target)
		tgArns = append(tgArns, rule.Switch.New, rule.Switch.Old)
		ruleMap[rule.Target] = rule
	}

	if _, err := cfg.client.elbv2.DescribeTargetGroups(context.TODO(), &elbv2.DescribeTargetGroupsInput{
		TargetGroupArns: tgArns,
	}); err != nil {
		fmt.Fprintln(cfg.io.err, err.Error())
	}

	rulesData, err := cfg.client.elbv2.DescribeRules(context.TODO(), &elbv2.DescribeRulesInput{
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
					Type:           ruleMap[*ruleData.RuleArn].Switch.getType(*tgTuple.TargetGroupArn),
					TargetGroupArn: aws.ToString(tgTuple.TargetGroupArn),
					Weight:         aws.ToInt32(tgTuple.Weight),
				})
			}

			res.AwsListenerRules = append(res.AwsListenerRules, AwsListenerRuleData{
				Name:           ruleMap[*ruleData.RuleArn].Name,
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
