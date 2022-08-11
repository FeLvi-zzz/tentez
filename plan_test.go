package tentez

import (
	"reflect"
	"sort"
	"testing"
)

func TestGetTargetNames(t *testing.T) {
	cases := []struct {
		targets  map[TargetType]Targets
		expected []string
	}{
		{
			map[TargetType]Targets{},
			[]string{},
		},
		{
			map[TargetType]Targets{
				TargetTypeAwsListenerRule: AwsListenerRules{
					AwsListenerRule{
						Name:   "hoge",
						Target: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:listener/app/my-lb/0123456789abcdef/0123456789abcdef/0123456789abcdef",
						Switch: Switch{
							Old: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef",
							New: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210",
						},
					},
					AwsListenerRule{
						Name:   "foo",
						Target: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:listener/app/my-lb/0123456789abcdef/0123456789abcdef/0123456789abcdef",
						Switch: Switch{
							Old: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef",
							New: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210",
						},
					},
				},
			},
			[]string{"hoge", "foo"},
		},
		{
			map[TargetType]Targets{
				TargetTypeAwsListenerRule: AwsListenerRules{
					AwsListenerRule{
						Name:   "hoge",
						Target: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:listener/app/my-lb/0123456789abcdef/0123456789abcdef/0123456789abcdef",
						Switch: Switch{
							Old: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef",
							New: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210",
						},
					},
					AwsListenerRule{
						Name:   "foo",
						Target: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:listener/app/my-lb/0123456789abcdef/0123456789abcdef/0123456789abcdef",
						Switch: Switch{
							Old: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef",
							New: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210",
						},
					},
				},
				TargetTypeAwsListener: AwsListeners{
					AwsListener{
						Name:   "hoge",
						Target: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:listener/app/my-lb/0123456789abcdef/0123456789abcdef/0123456789abcdef",
						Switch: Switch{
							Old: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef",
							New: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210",
						},
					},
					AwsListener{
						Name:   "",
						Target: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:listener/app/my-lb/0123456789abcdef/0123456789abcdef/0123456789abcdef",
						Switch: Switch{
							Old: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets1/0123456789abcdef",
							New: "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/my-targets2/fedcba9876543210",
						},
					},
				},
			},
			[]string{"hoge", "foo", "hoge", ""},
		},
	}

	for _, c := range cases {
		got := getTargetNames(c.targets)
		sortSlice(c.expected)
		sortSlice(got)
		if !reflect.DeepEqual(c.expected, got) {
			t.Errorf("expected %v, but got %v", c.expected, got)
		}
	}
}

func sortSlice(s []string) {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}
