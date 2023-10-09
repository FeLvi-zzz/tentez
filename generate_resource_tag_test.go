package tentez

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	rgt "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	rgtTypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
)

func TestGenerateConfigFromResourceTags(t *testing.T) {
	type args struct {
		filterTags      map[string]string
		matchingTagKeys []string
		switchKey       string
		oldValue        string
		newValue        string
		cfg             Config
	}
	tests := []struct {
		name    string
		args    args
		want    YamlStruct
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				filterTags: map[string]string{
					"filter_key": "filter_value",
				},
				matchingTagKeys: []string{
					"matching_tag_key",
				},
				switchKey: "switch_key",
				oldValue:  "old_value",
				newValue:  "new_value",
				cfg: Config{
					client: Client{
						elbv2: elbv2Mock{
							DescribeTargetGroupsResult: elbv2MockDescribeTargetGroupsResult{
								Value: elbv2.DescribeTargetGroupsOutput{
									TargetGroups: []elbv2Types.TargetGroup{
										{
											LoadBalancerArns: []string{"arn:aws:elasticloadbalancing:ap-northeast-1:0123456789012:loadbalancer/app/tentez-alb/0123456789012345"},
										},
									},
								},
							},
							DescribeListenersResult: elbv2MockDescribeListenersResult{
								Value: elbv2.DescribeListenersOutput{
									Listeners: []elbv2Types.Listener{
										{
											ListenerArn: aws.String("listener_arn"),
											Port:        aws.Int32(80),
											DefaultActions: []elbv2Types.Action{
												{
													Type:           elbv2Types.ActionTypeEnumForward,
													TargetGroupArn: aws.String("old_tg_arn"),
													ForwardConfig: &elbv2Types.ForwardActionConfig{
														TargetGroups: []elbv2Types.TargetGroupTuple{
															{
																TargetGroupArn: aws.String("old_tg_arn"),
															},
														},
													},
												},
											},
										},
									},
								},
							},
							DescribeRulesResult: elbv2MockDescribeRulesResult{
								Value: elbv2.DescribeRulesOutput{
									Rules: []elbv2Types.Rule{
										{
											IsDefault: true,
											RuleArn:   aws.String("default_rule_arn"),
											Priority:  aws.String("default"),
											Actions: []elbv2Types.Action{
												{
													Type:           elbv2Types.ActionTypeEnumForward,
													TargetGroupArn: aws.String("old_tg_arn"),
													ForwardConfig: &elbv2Types.ForwardActionConfig{
														TargetGroups: []elbv2Types.TargetGroupTuple{
															{
																TargetGroupArn: aws.String("old_tg_arn"),
															},
														},
													},
												},
											},
										},
										{
											RuleArn:  aws.String("rule_arn"),
											Priority: aws.String("1"),
											Actions: []elbv2Types.Action{
												{
													Type:           elbv2Types.ActionTypeEnumForward,
													TargetGroupArn: aws.String("old_tg_arn"),
													ForwardConfig: &elbv2Types.ForwardActionConfig{
														TargetGroups: []elbv2Types.TargetGroupTuple{
															{
																TargetGroupArn: aws.String("old_tg_arn"),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						rgt: rgtMock{
							GetResourcesResult: rgtMockGetResourcesResult{
								Value: rgt.GetResourcesOutput{
									ResourceTagMappingList: []rgtTypes.ResourceTagMapping{
										{
											ResourceARN: aws.String("old_tg_arn"),
											Tags: []rgtTypes.Tag{
												{
													Key:   aws.String("filter_key"),
													Value: aws.String("filter_value"),
												},
												{
													Key:   aws.String("matching_tag_key"),
													Value: aws.String("matching_tag_value"),
												},
												{
													Key:   aws.String("switch_key"),
													Value: aws.String("old_value"),
												},
											},
										},
										{
											ResourceARN: aws.String("new_tg_arn"),
											Tags: []rgtTypes.Tag{
												{
													Key:   aws.String("filter_key"),
													Value: aws.String("filter_value"),
												},
												{
													Key:   aws.String("matching_tag_key"),
													Value: aws.String("matching_tag_value"),
												},
												{
													Key:   aws.String("switch_key"),
													Value: aws.String("new_value"),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: YamlStruct{
				Steps: defaultSteps,
				AwsListeners: []AwsListener{
					{
						Name:   "tentez-alb:80, matching_tag_value",
						Target: "listener_arn",
						Switch: Switch{
							Old: "old_tg_arn",
							New: "new_tg_arn",
						},
					},
				},
				AwsListenerRules: []AwsListenerRule{
					{
						Name:   "tentez-alb:80-1, matching_tag_value",
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateConfigFromResourceTags(context.TODO(), tt.args.filterTags, tt.args.matchingTagKeys, tt.args.switchKey, tt.args.oldValue, tt.args.newValue, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateConfigFromResourceTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateConfigFromResourceTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
