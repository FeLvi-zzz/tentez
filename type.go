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

type YamlStruct struct {
	Steps            []Step           `yaml:"steps"`
	AwsListeners     AwsListeners     `yaml:"aws_listeners"`
	AwsListenerRules AwsListenerRules `yaml:"aws_listener_rules"`
}

type AwsTargetGroupTuple struct {
	TargetGroupArn string `yaml:"arn"`
	Weight         int32  `yaml:"weight"`
}
