package tentez

import (
	"encoding/json"
	"fmt"
)

type TerraformPlanJson struct {
	FormatVersion   string           `json:"format_version"`
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Change  Change `json:"change"`
}

type Change struct {
	Actions ChangeActions   `json:"actions"`
	Before  json.RawMessage `json:"before"`
	After   json.RawMessage `json:"after"`
}

type ChangeActions []ChangeAction
type ChangeAction string

type ChangeValueInterface interface {
	GetTarget() (string, error)
	GetSwitchTarget() (string, error)
}

type ChangeValueAwsListener struct {
	Arn            string                         `json:"arn"`
	DefaultActions []ChangeValueElbV2TargetAction `json:"default_action"`
}

type ChangeValueAwsListenerRule struct {
	Arn     string                         `json:"arn"`
	Actions []ChangeValueElbV2TargetAction `json:"action"`
}

type ChangeValueElbV2TargetAction struct {
	TargetGroupArn string                                `json:"target_group_arn"`
	Forward        []ChangeValueElbV2TargetActionForward `json:"forward"`
}

type ChangeValueElbV2TargetActionForward struct {
	TargetGroups []ChangeValueElbV2TargetActionForwardTargetGroups `json:"target_group"`
}

type ChangeValueElbV2TargetActionForwardTargetGroups struct {
	Arn    string `json:"arn"`
	Weight int    `json:"weight"`
}

func (c ChangeValueAwsListener) GetSwitchTarget() (string, error) {
	if len(c.DefaultActions) != 1 {
		return "", fmt.Errorf("cannot get default_action; %+v", c)
	}

	if len(c.DefaultActions) == 1 && c.DefaultActions[0].TargetGroupArn != "" {
		return c.DefaultActions[0].TargetGroupArn, nil
	}

	if len(c.DefaultActions[0].Forward) == 1 {
		for _, tg := range c.DefaultActions[0].Forward[0].TargetGroups {
			if tg.Weight > 0 {
				return tg.Arn, nil
			}
		}
	}

	return "", fmt.Errorf("cannot find target group")
}

func (c ChangeValueAwsListener) GetTarget() (string, error) {
	return c.Arn, nil
}

func (c ChangeValueAwsListenerRule) GetSwitchTarget() (string, error) {
	if len(c.Actions) != 1 {
		return "", fmt.Errorf("cannot get action; %+v", c)
	}

	if len(c.Actions) == 1 && c.Actions[0].TargetGroupArn != "" {
		return c.Actions[0].TargetGroupArn, nil
	}

	if len(c.Actions[0].Forward) == 1 {
		for _, tg := range c.Actions[0].Forward[0].TargetGroups {
			if tg.Weight > 0 {
				return tg.Arn, nil
			}
		}
	}

	return "", fmt.Errorf("cannot find target group")
}

func (c ChangeValueAwsListenerRule) GetTarget() (string, error) {
	return c.Arn, nil
}

const (
	ChangeActionNoop   ChangeAction = "no-op"
	ChangeActionCreate ChangeAction = "create"
	ChangeActionRead   ChangeAction = "read"
	ChangeActionUpdate ChangeAction = "update"
	ChangeActionDelete ChangeAction = "delete"
)

func (c ChangeActions) IsUpdate() bool {
	return len(c) == 1 && c[0] == ChangeActionUpdate
}

func GenerateConfigFromTerraformPlanJsons(jsons []TerraformPlanJson) (YamlStruct, error) {
	allAwsListeners := []AwsListener{}
	allAwsListenerRules := []AwsListenerRule{}

	for _, json := range jsons {
		awsListeners, err := getAwsListenersFromTerraformJson(json)
		if err != nil {
			return YamlStruct{}, err
		}
		allAwsListeners = append(allAwsListeners, awsListeners...)

		awsListenerRules, err := getAwsListenerRulesFromTerraformJson(json)
		if err != nil {
			return YamlStruct{}, err
		}
		allAwsListenerRules = append(allAwsListenerRules, awsListenerRules...)
	}

	return YamlStruct{
		Steps:            defaultSteps,
		AwsListeners:     allAwsListeners,
		AwsListenerRules: allAwsListenerRules,
	}, nil
}

func getAwsListenersFromTerraformJson(tfplanjson TerraformPlanJson) ([]AwsListener, error) {
	awsListeners := []AwsListener{}
	for _, resourceChange := range tfplanjson.ResourceChanges {
		if !resourceChange.Change.Actions.IsUpdate() || resourceChange.Type != "aws_lb_listener" {
			continue
		}

		var before, after ChangeValueAwsListener
		if err := json.Unmarshal(resourceChange.Change.Before, &before); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(resourceChange.Change.After, &after); err != nil {
			return nil, err
		}

		target, err := before.GetTarget()
		if err != nil {
			return nil, err
		}

		switchOld, err := before.GetSwitchTarget()
		if err != nil {
			return nil, err
		}

		switchAfter, err := after.GetSwitchTarget()
		if err != nil {
			return nil, err
		}

		awsListeners = append(awsListeners, AwsListener{
			Name:   resourceChange.Address,
			Target: target,
			Switch: Switch{
				Old: switchOld,
				New: switchAfter,
			},
		})
	}

	return awsListeners, nil
}

func getAwsListenerRulesFromTerraformJson(tfplanjson TerraformPlanJson) ([]AwsListenerRule, error) {
	awsListenerRules := []AwsListenerRule{}
	for _, resourceChange := range tfplanjson.ResourceChanges {
		if !resourceChange.Change.Actions.IsUpdate() || resourceChange.Type != "aws_lb_listener_rule" {
			continue
		}

		var before, after ChangeValueAwsListenerRule
		if err := json.Unmarshal(resourceChange.Change.Before, &before); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(resourceChange.Change.After, &after); err != nil {
			return nil, err
		}

		target, err := before.GetTarget()
		if err != nil {
			return nil, err
		}

		switchOld, err := before.GetSwitchTarget()
		if err != nil {
			return nil, err
		}

		switchAfter, err := after.GetSwitchTarget()
		if err != nil {
			return nil, err
		}

		awsListenerRules = append(awsListenerRules, AwsListenerRule{
			Name:   resourceChange.Address,
			Target: target,
			Switch: Switch{
				Old: switchOld,
				New: switchAfter,
			},
		})
	}

	return awsListenerRules, nil
}
