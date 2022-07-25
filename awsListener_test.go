package tentez

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestAwsListener_execSwitch(t *testing.T) {
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
				ModifyListenerError: fmt.Errorf("error"),
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
		},
	}

	for _, c := range cases {
		err := c.awsListener.execSwitch(c.weight, c.isForce, Config{
			client: Client{
				elbv2: c.elbv2Mock,
			},
			io: IOStreams{
				in:  bytes.NewBufferString(""),
				out: bytes.NewBufferString(""),
				err: bytes.NewBufferString(""),
			},
		})

		if c.isError != (err != nil) {
			t.Errorf("%s: expect isError == %t, but got %v", c.awsListener.Name, c.isError, err)
		}
	}
}

func TestAwsListeners_fetchData(t *testing.T) {
	cases := []struct {
		isError      bool
		expect       interface{}
		awsListeners AwsListeners
		elbv2Mock    elbv2Mock
	}{
		{
			isError: false,
			expect: struct {
				AwsListeners []AwsListenerData `yaml:"aws_listeners"`
			}{
				AwsListeners: []AwsListenerData{
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
								TargetGroupArn: "NewTarget",
								Weight:         50,
								Type:           "unknown",
							},
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
				DescribeListenersError: fmt.Errorf("error"),
			},
		},
	}

	for _, c := range cases {
		got, gotErr := c.awsListeners.fetchData(Config{
			client: Client{
				elbv2: c.elbv2Mock,
			},
			io: IOStreams{
				in:  bytes.NewBufferString(""),
				out: bytes.NewBufferString(""),
				err: bytes.NewBufferString(""),
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
