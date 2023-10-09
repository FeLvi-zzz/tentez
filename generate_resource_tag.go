package tentez

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Type "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	rgt "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	rgtTypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
)

type GenerateConfigResourceTagVersion string

type GenerateConfigResourceTag struct {
	Version GenerateConfigResourceTagVersion `yaml:"version"`
}

const (
	GenerateConfigResourceTagVersionV1beta1 GenerateConfigResourceTagVersion = "tentez.ResourceTag.v1beta1"
)

type GenerateConfigResourceTagV1beta1 struct {
	Spec GenerateConfigResourceTagV1beta1Spec `yaml:"spec"`
}

type GenerateConfigResourceTagV1beta1Spec struct {
	FilterTags      map[string]string `yaml:"filterTags"`
	MatchingTagKeys []string          `yaml:"matchingTagKeys"`
	SwitchTag       struct {
		Key   string `yaml:"key"`
		Value struct {
			Old string `yaml:"old"`
			New string `yaml:"new"`
		} `yaml:"value"`
	} `yaml:"switchTag"`
}

func GenerateConfigFromResourceTags(
	ctx context.Context,
	filterTags map[string]string,
	matchingTagKeys []string,
	switchKey string,
	oldValue string,
	newValue string,
	cfg Config,
) (YamlStruct, error) {
	yaml := YamlStruct{
		Steps: defaultSteps,
	}

	// list targetGroups and filter them
	paginator := rgt.NewGetResourcesPaginator(cfg.client.rgt, &rgt.GetResourcesInput{
		ResourceTypeFilters: []string{"elasticloadbalancing:targetgroup"},
		TagFilters:          buildTagFilters(filterTags, matchingTagKeys, switchKey, oldValue, newValue),
	})
	resourceTagMappingList := []rgtTypes.ResourceTagMapping{}
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return YamlStruct{}, err
		}
		resourceTagMappingList = append(resourceTagMappingList, output.ResourceTagMappingList...)
	}

	tgArnKeyMap, tgKeyArnMap := buildTgMap(resourceTagMappingList, matchingTagKeys, switchKey, oldValue, newValue)

	// fetch listeners from filtered targetGroups
	for tgArn, switchMapKey := range tgArnKeyMap {
		lbArn, lbName, err := getLbFromTg(ctx, tgArn, cfg)
		if err != nil {
			return YamlStruct{}, err
		}
		if lbArn == "" {
			continue
		}

		listenerOutput, err := cfg.client.elbv2.DescribeListeners(ctx, &elbv2.DescribeListenersInput{
			LoadBalancerArn: aws.String(lbArn),
		})
		if err != nil {
			return YamlStruct{}, err
		}

		for _, listener := range listenerOutput.Listeners {
			ruleOutput, err := cfg.client.elbv2.DescribeRules(ctx, &elbv2.DescribeRulesInput{
				ListenerArn: aws.String(*listener.ListenerArn),
			})
			if err != nil {
				return YamlStruct{}, err
			}

			for _, rule := range ruleOutput.Rules {
				for _, action := range rule.Actions {
					if action.Type != elbv2Type.ActionTypeEnumForward {
						continue
					}
					for _, tg := range action.ForwardConfig.TargetGroups {
						if tgArn == *tg.TargetGroupArn && (tg.Weight == nil || *tg.Weight > 0) {
							if rule.IsDefault {
								yaml.AwsListeners = append(yaml.AwsListeners, AwsListener{
									Name:   fmt.Sprintf("%s:%d, %s", lbName, *listener.Port, switchMapKey),
									Target: aws.ToString(listener.ListenerArn),
									Switch: Switch{
										Old: tgKeyArnMap[switchMapKey].Old,
										New: tgKeyArnMap[switchMapKey].New,
									},
								})
							} else {
								yaml.AwsListenerRules = append(yaml.AwsListenerRules, AwsListenerRule{
									Name:   fmt.Sprintf("%s:%d-%s, %s", lbName, *listener.Port, *rule.Priority, switchMapKey),
									Target: aws.ToString(rule.RuleArn),
									Switch: Switch{
										Old: tgKeyArnMap[switchMapKey].Old,
										New: tgKeyArnMap[switchMapKey].New,
									},
								})
							}
						}
					}
				}
			}
		}
	}

	return yaml, nil
}

func buildTagFilters(filterTags map[string]string, matchingTagKeys []string, switchKey string, oldValue string, newValue string) []rgtTypes.TagFilter {
	tagFileters := []rgtTypes.TagFilter{}

	for k, v := range filterTags {
		tagFileters = append(tagFileters, rgtTypes.TagFilter{
			Key:    aws.String(k),
			Values: []string{v},
		})
	}

	for _, k := range matchingTagKeys {
		tagFileters = append(tagFileters, rgtTypes.TagFilter{
			Key: aws.String(k),
		})
	}

	tagFileters = append(tagFileters, []rgtTypes.TagFilter{
		{
			Key:    aws.String(switchKey),
			Values: []string{oldValue, newValue},
		},
	}...)

	return tagFileters
}

func getLbFromTg(ctx context.Context, tgArn string, cfg Config) (string, string, error) {
	tgOutput, err := cfg.client.elbv2.DescribeTargetGroups(ctx, &elbv2.DescribeTargetGroupsInput{
		TargetGroupArns: []string{tgArn},
	})
	if err != nil {
		return "", "", err
	}

	lbs := tgOutput.TargetGroups[0].LoadBalancerArns
	if len(lbs) == 0 {
		return "", "", nil
	}

	lbArn := lbs[0]

	lbARN, err := arn.Parse(lbArn)
	if err != nil {
		return "", "", err
	}
	lbName := strings.Split(lbARN.Resource, "/")[2]

	return lbArn, lbName, nil
}

func buildTgMap(resourceTagMappingList []rgtTypes.ResourceTagMapping, matchingTagKeys []string, switchKey string, oldValue string, newValue string) (
	map[string]string,
	map[string]struct {
		Old string
		New string
	},
) {
	tgArnKeyMap := map[string]string{} // key: tgArn, value: switchMapKey
	tgKeyArnMap := map[string]struct { // key: switchMapKey, value: tgArns
		Old string // tgArn
		New string // tgArn
	}{}
	for _, tg := range resourceTagMappingList {
		tgArn := *tg.ResourceARN

		tagMap := map[string]string{}
		for _, tag := range tg.Tags {
			tagMap[*tag.Key] = *tag.Value
		}

		keys := []string{}
		for _, keyTag := range matchingTagKeys {
			v, ok := tagMap[keyTag]
			if !ok {
				break
			}
			keys = append(keys, v)
		}
		if len(keys) != len(matchingTagKeys) {
			continue
		}

		switchMapKey := strings.Join(keys, ",")

		switch tagMap[switchKey] {
		case oldValue:
			t := tgKeyArnMap[switchMapKey]
			t.Old = tgArn
			tgKeyArnMap[switchMapKey] = t
		case newValue:
			t := tgKeyArnMap[switchMapKey]
			t.New = tgArn
			tgKeyArnMap[switchMapKey] = t
		}
		tgArnKeyMap[tgArn] = switchMapKey
	}

	return tgArnKeyMap, tgKeyArnMap
}
