package tentez

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	rgt "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
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

func (t targetsMock) fetchData(cfg Config) (TargetsData, error) {
	return []Target{}, nil
}

type elbv2MockModifyListenerResult struct {
	Value elbv2.ModifyListenerOutput
	Error error
}

type elbv2MockModifyRuleResult struct {
	Value elbv2.ModifyRuleOutput
	Error error
}

type elbv2MockDescribeRulesResult struct {
	Value elbv2.DescribeRulesOutput
	Error error
}

type elbv2MockDescribeListenersResult struct {
	Value elbv2.DescribeListenersOutput
	Error error
}

type elbv2MockDescribeTargetGroupsResult struct {
	Value elbv2.DescribeTargetGroupsOutput
	Error error
}
type elbv2Mock struct {
	ModifyListenerResult       elbv2MockModifyListenerResult
	ModifyRuleResult           elbv2MockModifyRuleResult
	DescribeRulesResult        elbv2MockDescribeRulesResult
	DescribeListenersResult    elbv2MockDescribeListenersResult
	DescribeTargetGroupsResult elbv2MockDescribeTargetGroupsResult
}

type rgtMockGetResourcesResult struct {
	Value rgt.GetResourcesOutput
	Error error
}

type rgtMock struct {
	GetResourcesResult rgtMockGetResourcesResult
}

type clockMock struct{}

func (c clockMock) Sleep(time.Duration) {}

func NewDummyActions() []elbv2Types.Action {
	return []elbv2Types.Action{
		{
			Type: elbv2Types.ActionTypeEnumForward,
			ForwardConfig: &elbv2Types.ForwardActionConfig{
				TargetGroups: []elbv2Types.TargetGroupTuple{
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
	return &m.ModifyListenerResult.Value, m.ModifyListenerResult.Error
}
func (m elbv2Mock) ModifyRule(ctx context.Context, params *elbv2.ModifyRuleInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyRuleOutput, error) {
	return &m.ModifyRuleResult.Value, m.ModifyRuleResult.Error
}
func (m elbv2Mock) DescribeRules(ctx context.Context, params *elbv2.DescribeRulesInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeRulesOutput, error) {
	return &m.DescribeRulesResult.Value, m.DescribeRulesResult.Error
}
func (m elbv2Mock) DescribeListeners(ctx context.Context, params *elbv2.DescribeListenersInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeListenersOutput, error) {
	return &m.DescribeListenersResult.Value, m.DescribeListenersResult.Error
}
func (m elbv2Mock) DescribeTargetGroups(ctx context.Context, params *elbv2.DescribeTargetGroupsInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeTargetGroupsOutput, error) {
	return &m.DescribeTargetGroupsResult.Value, m.DescribeTargetGroupsResult.Error
}

func (m rgtMock) GetResources(ctx context.Context, params *rgt.GetResourcesInput, optFns ...func(*rgt.Options)) (*rgt.GetResourcesOutput, error) {
	return &m.GetResourcesResult.Value, m.GetResourcesResult.Error
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
			clock: clockMock{},
		})

		if err != nil {
			t.Errorf("expected no error, but throw %v", err)
		}
	}
}
