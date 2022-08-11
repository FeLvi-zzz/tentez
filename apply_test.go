package tentez

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type targetsMock []targetMock
type targetMock struct {
	Name          string
	CurrentWeight Weight
}

func (t targetMock) getName() string {
	return t.Name
}

func (t targetMock) execSwitch(targetWeight Weight, isForce bool, cfg Config) error {
	if isForce {
		return nil
	}

	if t.CurrentWeight.CalcOldRatio() < targetWeight.CalcOldRatio() {
		return SkipSwitchError{"the old weight target is larger than current one."}
	}
	if t.CurrentWeight.CalcNewRatio() > targetWeight.CalcNewRatio() {
		return SkipSwitchError{"the new weight target is smaller than current one."}
	}

	return nil
}

func (t targetsMock) targetsSlice() []Target {
	res := []Target{}
	for _, v := range t {
		res = append(res, v)
	}
	return res
}

func (t targetsMock) fetchData(cfg Config) (interface{}, error) {
	return []Target{}, nil
}

type elbv2Mock struct {
	ModifyListenerError    error
	ModifyRuleError        error
	DescribeRulesError     error
	DescribeListenersError error
}

func NewDummyActions() []types.Action {
	return []types.Action{
		{
			ForwardConfig: &types.ForwardActionConfig{
				TargetGroups: []types.TargetGroupTuple{
					{
						TargetGroupArn: aws.String("oldTarget"),
						Weight:         aws.Int32(50),
					},
					{
						TargetGroupArn: aws.String("newTarget"),
						Weight:         aws.Int32(50),
					},
				},
			},
		},
	}
}

func (m elbv2Mock) ModifyListener(ctx context.Context, params *elbv2.ModifyListenerInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyListenerOutput, error) {
	return &elbv2.ModifyListenerOutput{}, m.ModifyListenerError
}
func (m elbv2Mock) ModifyRule(ctx context.Context, params *elbv2.ModifyRuleInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyRuleOutput, error) {
	return &elbv2.ModifyRuleOutput{}, m.ModifyRuleError
}
func (m elbv2Mock) DescribeRules(ctx context.Context, params *elbv2.DescribeRulesInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeRulesOutput, error) {
	return &elbv2.DescribeRulesOutput{
		Rules: []types.Rule{
			{
				RuleArn: &params.RuleArns[0],
				Actions: NewDummyActions(),
			},
		},
	}, m.DescribeRulesError
}
func (m elbv2Mock) DescribeListeners(ctx context.Context, params *elbv2.DescribeListenersInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeListenersOutput, error) {
	return &elbv2.DescribeListenersOutput{
		Listeners: []types.Listener{
			{
				ListenerArn:    &params.ListenerArns[0],
				DefaultActions: NewDummyActions(),
			},
		},
	}, m.DescribeListenersError
}
func (m elbv2Mock) DescribeTargetGroups(ctx context.Context, params *elbv2.DescribeTargetGroupsInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeTargetGroupsOutput, error) {
	return &elbv2.DescribeTargetGroupsOutput{}, m.DescribeListenersError
}

func TestExecSwitch(t *testing.T) {
	cases := []struct {
		Name         string
		IsForce      bool
		IsError      bool
		Targets      map[TargetType]Targets
		TargetWeight Weight
	}{
		{
			Name:    "success",
			IsError: false,
			IsForce: false,
			Targets: map[TargetType]Targets{
				"targets_mock": targetsMock([]targetMock{
					{
						Name: "target_mock",
						CurrentWeight: Weight{
							Old: 50,
							New: 50,
						},
					},
				}),
			},
			TargetWeight: Weight{
				Old: 0,
				New: 100,
			},
		},
		{
			Name:    "skip",
			IsError: false,
			IsForce: false,
			Targets: map[TargetType]Targets{
				"targets_mock": targetsMock([]targetMock{
					{
						Name: "target_mock",
						CurrentWeight: Weight{
							Old: 0,
							New: 100,
						},
					},
				}),
			},
			TargetWeight: Weight{
				Old: 100,
				New: 0,
			},
		},
	}
	for _, c := range cases {
		targets := c.Targets
		err := execSwitch(targets, c.TargetWeight, c.IsForce, Config{
			client: Client{
				elbv2: elbv2Mock{},
			},
			io: IOStreams{
				in:  bytes.NewBufferString(""),
				out: bytes.NewBufferString(""),
				err: bytes.NewBufferString(""),
			},
		})

		if err != nil {
			t.Errorf("expected no error, but throw %v", err)
		}
	}
}
