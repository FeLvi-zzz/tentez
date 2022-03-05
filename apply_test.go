package tentez

import (
	"bytes"
	"context"
	"testing"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type targetsMock struct{}

func (t targetsMock) targetsSlice() []Target {
	return []Target{}
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

func (m elbv2Mock) ModifyListener(ctx context.Context, params *elbv2.ModifyListenerInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyListenerOutput, error) {
	return &elbv2.ModifyListenerOutput{}, m.ModifyListenerError
}
func (m elbv2Mock) ModifyRule(ctx context.Context, params *elbv2.ModifyRuleInput, optFns ...func(*elbv2.Options)) (*elbv2.ModifyRuleOutput, error) {
	return &elbv2.ModifyRuleOutput{}, m.ModifyRuleError
}
func (m elbv2Mock) DescribeRules(ctx context.Context, params *elbv2.DescribeRulesInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeRulesOutput, error) {
	return &elbv2.DescribeRulesOutput{}, m.DescribeRulesError
}
func (m elbv2Mock) DescribeListeners(ctx context.Context, params *elbv2.DescribeListenersInput, optFns ...func(*elbv2.Options)) (*elbv2.DescribeListenersOutput, error) {
	return &elbv2.DescribeListenersOutput{}, m.DescribeListenersError
}

func TestExecSwitch(t *testing.T) {
	targets := map[string]Targets{
		"targets_mock": targetsMock{},
	}
	err := execSwitch(targets, Weight{}, Config{
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
