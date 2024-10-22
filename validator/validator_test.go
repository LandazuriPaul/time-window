package validator

import (
	"testing"
	"time"
)

func TestValidator(t *testing.T) {
	tests := []struct {
		name                  string
		allowedInput          string
		blockedInput          string
		forceValidRegexpInput string
		commitMessageInput    string
		timestampInput        int64
		expected              Result
	}{
		{
			"simple valid case",
			`- name: Office hours
  cronExpression: "0 9 * * 1-5"
  duration: 8h`,
			"",
			`force\-time\-window`,
			`normal commit message`,
			1729525158, // 2024-10-21T15:39:18.000Z,
			Result{
				"",
				true,
				"timestamp in allowed window 'Office hours'",
				1729525158, // 2024-10-21T15:39:18.000Z,
			},
		},
		{
			"simple invalid case",
			`- name: Office hours
  cronExpression: "0 9 * * 1-5"
  duration: 8h`,
			"",
			`force\-time\-window`,
			`normal commit message`,
			1729537200, // 2024-10-21T19:00:00.000Z
			Result{
				"",
				false,
				"the timestamp isn't in any allowed window",
				1729537200, // 2024-10-21T19:00:00.000Z
			},
		},
		{
			"force valid case",
			`- name: Office hours
  cronExpression: "0 9 * * 1-5"
  duration: 8h`,
			"",
			`force\-time\-window`,
			"normal commit message force-time-window",
			1729537200, // 2024-10-21T19:00:00.000Z
			Result{
				"",
				true,
				`forced validated with regexp force\-time\-window and commit message normal commit message force-time-window`,
				1729537200, // 2024-10-21T19:00:00.000Z
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v, err := NewValidator(&test.allowedInput, &test.blockedInput, &test.forceValidRegexpInput, &test.commitMessageInput)
			if err != nil {
				t.Errorf("unexpected error getting new validator for %s: %s", test.name, err)
			}
			r := v.Validate(time.Unix(test.timestampInput, 0))
			if r != test.expected {
				t.Errorf("expected result %s, got %s", test.expected, r)
			}
		})
	}
}
