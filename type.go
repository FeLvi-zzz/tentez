package tentez

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type TargetType string

const (
	TargetTypeAwsListener     TargetType = "aws_listeners"
	TargetTypeAwsListenerRule TargetType = "aws_listener_rules"
)

const (
	maxDescribeTargetGroupsItems = 20
	maxDescribeRulesItems        = 20
	maxDescribeListenersItems    = 20
)

type Weight struct {
	Old int32
	New int32
}

func (w Weight) CalcOldRatio() float64 {
	return float64(w.Old) / float64(w.New+w.Old)
}

func (w Weight) CalcNewRatio() float64 {
	return float64(w.New) / float64(w.New+w.Old)
}

type Switch struct {
	Old string
	New string
}

func (s Switch) getType(t string) string {
	switch t {
	case s.Old:
		return "old"
	case s.New:
		return "new"
	default:
		return "unknown"
	}
}

func makeNewActions(actions []elbv2Types.Action, targetSwitch Switch, targetWeight Weight) []elbv2Types.Action {
	res := make([]elbv2Types.Action, len(actions))

	for i, action := range actions {
		if action.Type != elbv2Types.ActionTypeEnumForward {
			if action.Type == elbv2Types.ActionTypeEnumAuthenticateOidc {
				action.AuthenticateOidcConfig.UseExistingClientSecret = aws.Bool(true)
			}
			res[i] = action
			continue
		}

		switch {
		case targetWeight.Old == 0:
			res[i] = elbv2Types.Action{
				Type:           elbv2Types.ActionTypeEnumForward,
				TargetGroupArn: aws.String(targetSwitch.New),
				Order:          action.Order,
			}
		case targetWeight.New == 0:
			res[i] = elbv2Types.Action{
				Type:           elbv2Types.ActionTypeEnumForward,
				TargetGroupArn: aws.String(targetSwitch.Old),
				Order:          action.Order,
			}
		default:
			res[i] = elbv2Types.Action{
				Type: elbv2Types.ActionTypeEnumForward,
				ForwardConfig: &elbv2Types.ForwardActionConfig{
					TargetGroups: []elbv2Types.TargetGroupTuple{
						{
							TargetGroupArn: aws.String(targetSwitch.Old),
							Weight:         aws.Int32(targetWeight.Old),
						},
						{
							TargetGroupArn: aws.String(targetSwitch.New),
							Weight:         aws.Int32(targetWeight.New),
						},
					},
				},
				Order: action.Order,
			}
		}
	}

	return res
}

type Step struct {
	Type         StepType `yaml:"type"`
	Weight       Weight   `yaml:"weight,omitempty"`
	SleepSeconds int      `yaml:"sleepSeconds,omitempty"`
}

type StepType string

const (
	StepTypePause  StepType = "pause"
	StepTypeSleep  StepType = "sleep"
	StepTypeSwitch StepType = "switch"
)

type YamlStruct struct {
	Steps            []Step           `yaml:"steps"`
	AwsListeners     AwsListeners     `yaml:"aws_listeners"`
	AwsListenerRules AwsListenerRules `yaml:"aws_listener_rules"`
}

type AwsTargetGroupTuple struct {
	TargetGroupArn string `yaml:"arn"`
	Weight         int32  `yaml:"weight"`
	Type           string `yaml:"type"`
}

type SkipSwitchError struct {
	Message string
}

func (s SkipSwitchError) Error() string {
	return fmt.Sprintf("skip switching: %s", s.Message)
}

type FailedFetchTargetGroupsError struct {
	tgs []string
}

func NewFailedFetchTargetGroupsError(tgs []string) error {
	if len(tgs) == 0 {
		return nil
	}
	return &FailedFetchTargetGroupsError{tgs: tgs}
}

func (f *FailedFetchTargetGroupsError) Error() string {
	return fmt.Sprintf("TargetGroupsNotFound: %v", f.tgs)
}

func (f *FailedFetchTargetGroupsError) Is(target error) bool {
	_, ok := target.(*FailedFetchTargetGroupsError)
	return f != nil && ok
}
