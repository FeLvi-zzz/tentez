package tentez

import "testing"

type tentezMock struct {
	Steps      []Step
	Targets    map[string]Targets
	planCount  int
	applyCount int
	getCount   int
	helpCount  int
}

func (t *tentezMock) plan() error {
	t.planCount++
	return nil
}

func (t *tentezMock) apply() error {
	t.applyCount++
	return nil
}

func (t *tentezMock) get() error {
	t.getCount++
	return nil
}

func (t *tentezMock) help() {
	t.helpCount++
}

func TestExec(t *testing.T) {
	type expected struct {
		isError    bool
		planCount  int
		applyCount int
		getCount   int
		helpCount  int
	}
	cases := []struct {
		cmd      string
		expected expected
	}{
		{
			cmd: "plan",
			expected: expected{
				planCount: 1,
			},
		},
		{
			cmd: "apply",
			expected: expected{
				planCount:  1,
				applyCount: 1,
			},
		},
		{
			cmd: "get",
			expected: expected{
				getCount: 1,
			},
		},
		{
			cmd: "help",
			expected: expected{
				helpCount: 1,
			},
		},
		{
			cmd: "",
			expected: expected{
				helpCount: 1,
			},
		},
		{
			cmd: "hoge",
			expected: expected{
				isError:   true,
				helpCount: 1,
			},
		},
	}

	for _, c := range cases {
		te := &tentezMock{
			Steps:   []Step{},
			Targets: map[string]Targets{},
		}
		err := Exec(te, c.cmd)
		if !((err != nil) == c.expected.isError &&
			te.applyCount == c.expected.applyCount &&
			te.planCount == c.expected.planCount &&
			te.getCount == c.expected.getCount &&
			te.helpCount == c.expected.helpCount) {
			t.Errorf("%s: expect %+v, but got %+v", c.cmd, c.expected, struct {
				isError    bool
				planCount  int
				applyCount int
				getCount   int
				helpCount  int
			}{
				isError:    err != nil,
				planCount:  te.planCount,
				applyCount: te.applyCount,
				getCount:   te.getCount,
				helpCount:  te.helpCount,
			})
		}
	}
}
