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
	Name            string                `yaml:"name"`
	ListenerRuleArn string                `yaml:"target"`
	Weights         []AwsTargetGroupTuple `yaml:"weights"`
}

func (r AwsListenerRule) execSwitch(targetWeight Weight, isForce bool, cfg Config) error {
	// avoid rate limit
	time.Sleep(1 * time.Second)

	ruleData, err := cfg.client.elbv2.DescribeRules(context.TODO(), &elbv2.DescribeRulesInput{
		RuleArns: []string{r.Target},
	})
	if err != nil {
		return err
	}

	rule := ruleData.Rules[0]
	if rule.IsDefault {
		return fmt.Errorf("this is a default listener rule. Use `aws_listeners`")
	}

	tgWeight := Weight{}
	for _, action := range rule.Actions {
		if action.Type != elbv2Types.ActionTypeEnumForward {
			return fmt.Errorf("invalid action type: %s", action.Type)
		}

		for _, tgTuple := range action.ForwardConfig.TargetGroups {
			switch *tgTuple.TargetGroupArn {
			case r.Switch.Old:
				tgWeight.Old = *tgTuple.Weight
			case r.Switch.New:
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

	_, err = cfg.client.elbv2.ModifyRule(context.TODO(), &elbv2.ModifyRuleInput{
		RuleArn: aws.String(r.Target),
		Actions: compactActions(r.Switch, targetWeight),
	})

	return err
}

func (r AwsListenerRule) getName() string {
	return r.Name
}

func (rs AwsListenerRules) fetchData(cfg Config) (TargetsData, error) {
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

	for _, tgArnsBatch := range chunk(tgArns, maxDescribeTargetGroupsItems) {
		if _, err := cfg.client.elbv2.DescribeTargetGroups(context.TODO(), &elbv2.DescribeTargetGroupsInput{
			TargetGroupArns: tgArnsBatch,
		}); err != nil {
			fmt.Fprintln(cfg.io.err, err.Error())
		}
	}

	rules := []elbv2Types.Rule{}
	for _, ruleArnsBatch := range chunk(ruleArns, maxDescribeRulesItems) {
		rulesOutput, err := cfg.client.elbv2.DescribeRules(context.TODO(), &elbv2.DescribeRulesInput{
			RuleArns: ruleArnsBatch,
		})
		if err != nil {
			return nil, err
		}
		rules = append(rules, rulesOutput.Rules...)
	}

	res := []AwsListenerRuleData{}

	for _, rule := range rules {
		if rule.IsDefault {
			return nil, fmt.Errorf("%s is a default listener rule. Use `aws_listeners`", aws.ToString(rule.RuleArn))
		}

		for _, action := range rule.Actions {
			targetGroupTuples := []AwsTargetGroupTuple{}

			if action.Type == elbv2Types.ActionTypeEnumForward {
				for _, tgTuple := range action.ForwardConfig.TargetGroups {
					targetGroupTuples = append(targetGroupTuples, AwsTargetGroupTuple{
						Type:           ruleMap[aws.ToString(rule.RuleArn)].Switch.getType(aws.ToString(tgTuple.TargetGroupArn)),
						TargetGroupArn: aws.ToString(tgTuple.TargetGroupArn),
						Weight:         aws.ToInt32(tgTuple.Weight),
					})
				}
			}

			res = append(res, AwsListenerRuleData{
				Name:            ruleMap[aws.ToString(rule.RuleArn)].Name,
				ListenerRuleArn: aws.ToString(rule.RuleArn),
				Weights:         targetGroupTuples,
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
