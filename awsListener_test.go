package tentez

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2Types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"gopkg.in/yaml.v3"
)

func TestAwsListener_execSwitch(t *testing.T) {
	describeListenersResult := elbv2MockDescribeListenersResult{
		Value: elbv2.DescribeListenersOutput{
			Listeners: []elbv2Types.Listener{
				{
					ListenerArn: aws.String("validTarget"),
					DefaultActions: []elbv2Types.Action{
						NewDummyAuthOidcAction(),
						NewDummyForwardAction(),
					},
				},
			},
		},
	}
	noForwardDescribeListenersResult := elbv2MockDescribeListenersResult{
		Value: elbv2.DescribeListenersOutput{
			Listeners: []elbv2Types.Listener{
				{
					ListenerArn: aws.String("validTarget"),
					DefaultActions: []elbv2Types.Action{
						NewDummyAuthOidcAction(),
					},
				},
			},
		},
	}

	cases := []struct {
		isError     bool
		isForce     bool
		awsListener AwsListener
		weight      Weight
		elbv2Mock   elbv2Mock
	}{
		{
			isError: false,
			isForce: false,
			awsListener: AwsListener{
				Name:   "success",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			elbv2Mock: elbv2Mock{
				DescribeListenersResult: describeListenersResult,
			},
		},
		{
			isError: true,
			isForce: false,
			awsListener: AwsListener{
				Name:   "success",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			elbv2Mock: elbv2Mock{
				DescribeListenersResult: noForwardDescribeListenersResult,
			},
		},
		{
			isError: true,
			isForce: false,
			awsListener: AwsListener{
				Name:   "api_error",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			weight: Weight{
				Old: 30,
				New: 70,
			},
			elbv2Mock: elbv2Mock{
				DescribeListenersResult: describeListenersResult,
				ModifyListenerResult: elbv2MockModifyListenerResult{
					Error: fmt.Errorf("error"),
				},
			},
		},
		{
			isError: true,
			isForce: false,
			awsListener: AwsListener{
				Name:   "skip_switch_error",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			weight: Weight{
				Old: 100,
				New: 0,
			},
			elbv2Mock: elbv2Mock{
				DescribeListenersResult: describeListenersResult,
			},
		},
		{
			isError: false,
			isForce: true,
			awsListener: AwsListener{
				Name:   "success_force_switch",
				Target: "validTarget",
				Switch: Switch{
					Old: "oldTarget",
					New: "newTarget",
				},
			},
			weight: Weight{
				Old: 100,
				New: 0,
			},
			elbv2Mock: elbv2Mock{
				DescribeListenersResult: describeListenersResult,
			},
		},
	}

	for _, c := range cases {
		err := c.awsListener.execSwitch(context.TODO(), c.weight, c.isForce, Config{
			client: Client{
				elbv2: c.elbv2Mock,
			},
			clock: clockMock{},
		})

		if c.isError != (err != nil) {
			t.Errorf("%s: expect isError == %t, but got %v", c.awsListener.Name, c.isError, err)
		}
	}
}

func TestAwsListeners_fetchData(t *testing.T) {
	describeListenersResult := elbv2MockDescribeListenersResult{
		Value: elbv2.DescribeListenersOutput{
			Listeners: []elbv2Types.Listener{
				{
					ListenerArn: aws.String("validTarget"),
					DefaultActions: []elbv2Types.Action{
						NewDummyForwardAction(),
					},
				},
			},
		},
	}
	cases := []struct {
		isError      bool
		expect       interface{}
		awsListeners AwsListeners
		elbv2Mock    elbv2Mock
	}{
		{
			isError: false,
			expect: []AwsListenerData{
				{
					Name:       "success",
					ListnerArn: "validTarget",
					Weights: []AwsTargetGroupTuple{
						{
							TargetGroupArn: "oldTarget",
							Weight:         50,
							Type:           "old",
						},
						{
							TargetGroupArn: "newTarget",
							Weight:         50,
							Type:           "new",
						},
					},
				},
			},
			awsListeners: AwsListeners{
				AwsListener{
					Name:   "success",
					Target: "validTarget",
					Switch: Switch{
						Old: "oldTarget",
						New: "newTarget",
					},
				},
			},
			elbv2Mock: elbv2Mock{
				DescribeListenersResult: describeListenersResult,
			},
		},
		{
			isError: true,
			expect:  nil,
			awsListeners: AwsListeners{
				AwsListener{
					Name:   "success",
					Target: "validTarget",
					Switch: Switch{
						Old: "oldTarget",
						New: "newTarget",
					},
				},
			},
			elbv2Mock: elbv2Mock{
				DescribeListenersResult: elbv2MockDescribeListenersResult{
					Error: fmt.Errorf("error"),
				},
			},
		},
	}

	for _, c := range cases {
		got, gotErr := c.awsListeners.fetchData(context.TODO(), Config{
			client: Client{
				elbv2: c.elbv2Mock,
			},
		})

		expectedYaml, err := yaml.Marshal(c.expect)
		if err != nil {
			t.Errorf("cannot yaml.Marshal: %s", err)
		}
		gotYaml, err := yaml.Marshal(got)
		if err != nil {
			t.Errorf("cannot yaml.Marshal: %s", err)
		}

		if !reflect.DeepEqual(expectedYaml, gotYaml) {
			t.Errorf("expect %+v, but got %+v", c.expect, got)
		}
		if c.isError != (gotErr != nil) {
			t.Errorf("expect isError == %t, but got %v", c.isError, gotErr)
		}
	}
}
