package tentez

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/goccy/go-yaml"
)

func TestAwsListenerRule_execSwitch(t *testing.T) {
	describeRulesResult := elbv2MockDescribeRulesResult{
		Value: elbv2.DescribeRulesOutput{
			Rules: []elbv2Types.Rule{
				{
					Actions: []elbv2Types.Action{
						NewDummyAuthOidcAction(),
						NewDummyForwardAction(),
					},
				},
			},
		},
	}
	noForwardDescribeRulesResult := elbv2MockDescribeRulesResult{
		Value: elbv2.DescribeRulesOutput{
			Rules: []elbv2Types.Rule{
				{
					Actions: []elbv2Types.Action{
						NewDummyAuthOidcAction(),
					},
				},
			},
		},
	}

	cases := []struct {
		isError         bool
		isForce         bool
		awsListenerRule AwsListenerRule
		weight          Weight
		elbv2Mock       elbv2Mock
	}{
		{
			isError: false,
			isForce: false,
			awsListenerRule: AwsListenerRule{
				Name:   "success",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			elbv2Mock: elbv2Mock{
				DescribeRulesResult: describeRulesResult,
			},
		},
		{
			isError: true,
			isForce: false,
			awsListenerRule: AwsListenerRule{
				Name:   "success",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			elbv2Mock: elbv2Mock{
				DescribeRulesResult: noForwardDescribeRulesResult,
			},
		},
		{
			isError: true,
			isForce: false,
			awsListenerRule: AwsListenerRule{
				Name:   "api_error",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			weight: Weight{
				Old: 30,
				New: 70,
			},
			elbv2Mock: elbv2Mock{
				DescribeRulesResult: describeRulesResult,
				ModifyRuleResult: elbv2MockModifyRuleResult{
					Error: fmt.Errorf("error"),
				},
			},
		},
		{
			isError: true,
			isForce: false,
			awsListenerRule: AwsListenerRule{
				Name:   "skip_switch_error",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			weight: Weight{
				Old: 100,
				New: 0,
			},
			elbv2Mock: elbv2Mock{
				DescribeRulesResult: describeRulesResult,
			},
		},
		{
			isError: false,
			isForce: true,
			awsListenerRule: AwsListenerRule{
				Name:   "success_force_switch",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			weight: Weight{
				Old: 100,
				New: 0,
			},
			elbv2Mock: elbv2Mock{
				DescribeRulesResult: describeRulesResult,
			},
		},
	}

	for _, c := range cases {
		err := c.awsListenerRule.execSwitch(context.TODO(), c.weight, c.isForce, Config{
			client: Client{
				elbv2: c.elbv2Mock,
			},
			clock: clockMock{},
		})

		if c.isError != (err != nil) {
			t.Errorf("%s: expect isError == %t, but got %v", c.awsListenerRule.Name, c.isError, err)
		}
	}
}

func TestAwsListenerRules_fetchData(t *testing.T) {
	cases := []struct {
		isError          bool
		expect           interface{}
		awsListenerRules AwsListenerRules
		elbv2Mock        elbv2Mock
	}{
		{
			isError: false,
			expect: []AwsListenerRuleData{
				{
					Name:            "success",
					ListenerRuleArn: "validTarget",
					Weights: []AwsTargetGroupTuple{
						{
							TargetGroupArn: "oldTarget",
							Weight:         50,
							Type:           "old",
						},
						{
							TargetGroupArn: "newTarget",
							Weight:         50,
							Type:           "new",
						},
					},
				},
			},
			awsListenerRules: AwsListenerRules{
				AwsListenerRule{
					Name:   "success",
					Target: "validTarget",
					Switch: Switch{
						Old: "oldTarget",
						New: "newTarget",
					},
				},
			},
			elbv2Mock: elbv2Mock{
				DescribeRulesResult: elbv2MockDescribeRulesResult{
					Value: elbv2.DescribeRulesOutput{
						Rules: []elbv2Types.Rule{
							{
								RuleArn: aws.String("validTarget"),
								Actions: []elbv2Types.Action{
									NewDummyForwardAction(),
								},
							},
						},
					},
				},
			},
		},
		{
			isError: true,
			expect:  nil,
			awsListenerRules: AwsListenerRules{
				AwsListenerRule{
					Name:   "success",
					Target: "validTarget",
					Switch: Switch{
						Old: "oldTarget",
						New: "newTarget",
					},
				},
			},
			elbv2Mock: elbv2Mock{
				DescribeRulesResult: elbv2MockDescribeRulesResult{
					Error: fmt.Errorf("error"),
				},
			},
		},
	}

	for _, c := range cases {
		got, gotErr := c.awsListenerRules.fetchData(context.TODO(), Config{
			client: Client{
				elbv2: c.elbv2Mock,
			},
		})

		expectedYaml, err := yaml.Marshal(c.expect)
		if err != nil {
			t.Errorf("cannot yaml.Marshal: %s", err)
		}
		gotYaml, err := yaml.Marshal(got)
		if err != nil {
			t.Errorf("cannot yaml.Marshal: %s", err)
		}

		if !reflect.DeepEqual(expectedYaml, gotYaml) {
			t.Errorf("expect %+v, but got %+v", c.expect, got)
		}
		if c.isError != (gotErr != nil) {
			t.Errorf("expect isError == %t, but got %v", c.isError, gotErr)
		}
	}
}
