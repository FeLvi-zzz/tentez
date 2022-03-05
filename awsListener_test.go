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
		awsListener AwsListener
		weight      Weight
		elbv2Mock   elbv2Mock
	}{
		{
			isError: false,
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
			awsListener: AwsListener{
				Name:   "success",
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
	}

	for _, c := range cases {
		err := c.awsListener.execSwitch(c.weight, Config{
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
			t.Errorf("expect isError == %t, but got %v", c.isError, err)
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
			}{},
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
