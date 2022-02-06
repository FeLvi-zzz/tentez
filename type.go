package tentez

type Weight struct {
	Old int32
	New int32
}

type Switch struct {
	Old string
	New string
}

type Step struct {
	Type         string `yaml:"type"`
	Weight       Weight `yaml:"weight"`
	SleepSeconds int    `yaml:"sleepSeconds"`
}

type AwsListenerRule struct {
	Name   string `yaml:"name"`
	Target string `yaml:"target"`
	Switch Switch `yaml:"switch"`
}
type AwsListenerRules []AwsListenerRule

type AwsListener struct {
	Name   string `yaml:"name"`
	Target string `yaml:"target"`
	Switch Switch `yaml:"switch"`
}
type AwsListeners []AwsListener

type YamlStruct struct {
	Steps            []Step           `yaml:"steps"`
	AwsListeners     AwsListeners     `yaml:"aws_listeners"`
	AwsListenerRules AwsListenerRules `yaml:"aws_listener_rules"`
}

type AwsListenerRuleData struct {
	ListnerRuleArn string                `yaml:"target"`
	Weights        []AwsTargetGroupTuple `yaml:"weights"`
}
type AwsListenerData struct {
	ListnerArn string                `yaml:"target"`
	Weights    []AwsTargetGroupTuple `yaml:"weights"`
}

type AwsTargetGroupTuple struct {
	TargetGroupArn string `yaml:"arn"`
	Weight         int32  `yaml:"weight"`
}
