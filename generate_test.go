package tentez

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestChangeValueAwsListener_GetSwitchTarget(t *testing.T) {
	cs := []struct {
		c        ChangeValueAwsListener
		expected string
		isError  bool
	}{
		{
			c: ChangeValueAwsListener{
				Arn: "listener_arn",
				DefaultActions: []ChangeValueElbV2TargetAction{
					{
						TargetGroupArn: "tg_arn",
					},
				},
			},
			expected: "tg_arn",
		},
		{
			c: ChangeValueAwsListener{
				Arn: "listener_arn",
				DefaultActions: []ChangeValueElbV2TargetAction{
					{
						Forward: []ChangeValueElbV2TargetActionForward{
							{
								TargetGroups: []ChangeValueElbV2TargetActionForwardTargetGroups{
									{
										Arn:    "tg_new_arn",
										Weight: 1,
									},
									{
										Arn:    "tg_old_arn",
										Weight: 0,
									},
								},
							},
						},
					},
				},
			},
			expected: "tg_new_arn",
		},
		{
			c: ChangeValueAwsListener{
				Arn:            "listener_arn",
				DefaultActions: []ChangeValueElbV2TargetAction{},
			},
			isError: true,
		},
		{
			c: ChangeValueAwsListener{
				Arn: "listener_arn",
				DefaultActions: []ChangeValueElbV2TargetAction{
					{},
				},
			},
			isError: true,
		},
	}
	for _, c := range cs {
		got, err := c.c.GetSwitchTarget()
		if c.isError != (err != nil) {
			t.Errorf("expected isError is %t, got %v", c.isError, err)
		}
		if got != c.expected {
			t.Errorf("expected %s, got %s", c.expected, got)
		}
	}
}

func TestChangeValueAwsListener_GetTarget(t *testing.T) {
	cs := []struct {
		c        ChangeValueAwsListener
		expected string
	}{
		{
			c: ChangeValueAwsListener{
				Arn: "listener_arn",
			},
			expected: "listener_arn",
		},
	}

	for _, c := range cs {
		st, _ := c.c.GetTarget()
		if st != c.expected {
			t.Errorf("expected %s, got %s", c.expected, st)
		}
	}
}

func TestChangeValueAwsListenerRule_GetSwitchTarget(t *testing.T) {
	cs := []struct {
		c        ChangeValueAwsListenerRule
		expected string
		isError  bool
	}{
		{
			c: ChangeValueAwsListenerRule{
				Arn: "rule_arn",
				Actions: []ChangeValueElbV2TargetAction{
					{
						TargetGroupArn: "tg_arn",
					},
				},
			},
			expected: "tg_arn",
		},
		{
			c: ChangeValueAwsListenerRule{
				Arn: "rule_arn",
				Actions: []ChangeValueElbV2TargetAction{
					{
						Forward: []ChangeValueElbV2TargetActionForward{
							{
								TargetGroups: []ChangeValueElbV2TargetActionForwardTargetGroups{
									{
										Arn:    "tg_new_arn",
										Weight: 1,
									},
									{
										Arn:    "tg_old_arn",
										Weight: 0,
									},
								},
							},
						},
					},
				},
			},
			expected: "tg_new_arn",
		},
		{
			c: ChangeValueAwsListenerRule{
				Arn:     "rule_arn",
				Actions: []ChangeValueElbV2TargetAction{},
			},
			isError: true,
		},
		{
			c: ChangeValueAwsListenerRule{
				Arn: "listener_arn",
				Actions: []ChangeValueElbV2TargetAction{
					{},
				},
			},
			isError: true,
		},
	}
	for _, c := range cs {
		got, err := c.c.GetSwitchTarget()
		if c.isError != (err != nil) {
			t.Errorf("expected isError is %t, got %v", c.isError, err)
		}
		if got != c.expected {
			t.Errorf("expected %s, got %s", c.expected, got)
		}
	}
}

func TestChangeValueAwsListenerRule_GetTarget(t *testing.T) {
	cs := []struct {
		c        ChangeValueAwsListenerRule
		expected string
	}{
		{
			c: ChangeValueAwsListenerRule{
				Arn: "rule_arn",
			},
			expected: "rule_arn",
		},
	}

	for _, c := range cs {
		st, _ := c.c.GetTarget()
		if st != c.expected {
			t.Errorf("expected %s, got %s", c.expected, st)
		}
	}
}

func TestChangeActions_IsUpdate(t *testing.T) {
	cs := []struct {
		c        ChangeActions
		expected bool
	}{
		{
			c: []ChangeAction{
				ChangeActionUpdate,
			},
			expected: true,
		},
		{
			c: []ChangeAction{
				ChangeActionCreate,
			},
			expected: false,
		},
		{
			c:        []ChangeAction{},
			expected: false,
		},
	}

	for _, c := range cs {
		got := c.c.IsUpdate()
		if got != c.expected {
			t.Errorf("expected %t, got %t", c.expected, got)
		}
	}
}

func TestGenerateConfigFromTerraformPlanJson(t *testing.T) {
	cs := []struct {
		input    TerraformPlanJson
		expected YamlStruct
		isError  bool
	}{
		{
			input: TerraformPlanJson{
				ResourceChanges: []ResourceChange{
					{
						Address: "other_resource_address",
						Type:    "aws_hoge",
						Change: Change{
							Actions: []ChangeAction{
								ChangeActionUpdate,
							},
							After: json.RawMessage(
								([]byte)(`{}`),
							),
							Before: json.RawMessage(
								([]byte)(`{}`),
							),
						},
					},
					{
						Address: "other_change_action_address",
						Type:    "aws_lb_listener",
						Change: Change{
							Actions: []ChangeAction{
								ChangeActionCreate,
							},
							After: json.RawMessage(
								([]byte)(`{}`),
							),
							Before: json.RawMessage(
								([]byte)(`{}`),
							),
						},
					},
					{
						Address: "listener_address",
						Type:    "aws_lb_listener",
						Change: Change{
							Actions: []ChangeAction{
								ChangeActionUpdate,
							},
							After: json.RawMessage(
								([]byte)(`{"arn":"listener_arn","default_action":[{"target_group_arn":"new_tg_arn"}]}`),
							),
							Before: json.RawMessage(
								([]byte)(`{"arn":"listener_arn","default_action":[{"target_group_arn":"old_tg_arn"}]}`),
							),
						},
					},
					{
						Address: "rule_address",
						Type:    "aws_lb_listener_rule",
						Change: Change{
							Actions: []ChangeAction{
								ChangeActionUpdate,
							},
							After: json.RawMessage(
								([]byte)(`{"arn":"rule_arn","action":[{"target_group_arn":"new_tg_arn"}]}`),
							),
							Before: json.RawMessage(
								([]byte)(`{"arn":"rule_arn","action":[{"target_group_arn":"old_tg_arn"}]}`),
							),
						},
					},
				},
			},
			expected: YamlStruct{
				Steps: defaultSteps,
				AwsListeners: []AwsListener{
					{
						Name:   "listener_address",
						Target: "listener_arn",
						Switch: Switch{
							Old: "old_tg_arn",
							New: "new_tg_arn",
						},
					},
				},
				AwsListenerRules: []AwsListenerRule{
					{
						Name:   "rule_address",
						Target: "rule_arn",
						Switch: Switch{
							Old: "old_tg_arn",
							New: "new_tg_arn",
						},
					},
				},
			},
		},
	}

	for _, c := range cs {
		y, err := GenerateConfigFromTerraformPlanJson(c.input)
		if c.isError != (err != nil) {
			t.Errorf("expected isError is %t, got %v", c.isError, err)
		}
		if !reflect.DeepEqual(y, c.expected) {
			t.Errorf("expected %+v, got %+v", c.expected, y)
		}
	}
}
