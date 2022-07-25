package tentez

import "fmt"

type Weight struct {
	Old int32
	New int32
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

type Step struct {
	Type         string `yaml:"type"`
	Weight       Weight `yaml:"weight"`
	SleepSeconds int    `yaml:"sleepSeconds"`
}

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
