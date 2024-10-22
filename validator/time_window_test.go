package validator

import (
	"testing"
	"time"
)

func TestTimestampInTimeWindows(t *testing.T) {
	tests := []struct {
		name      string
		config    string
		timestamp int64
		expected  bool
	}{
		{
			"Inside office hours",
			`- name: Office hours
  cronExpression: "0 9 * * 1-5"
  duration: 8h`,
			1729537200, // 2024-10-21T19:00:00.000Z
			false,
		},
		{
			"After office hours",
			`- name: Office hours
  cronExpression: "0 9 * * 1-5"
  duration: 8h`,
			1729525158, // 2024-10-21T15:39:18.000Z
			true,
		},
		{
			"Before office hours",
			`- name: Office hours
  cronExpression: "0 9 * * 1-5"
  duration: 8h`,
			1729669158, // 2024-10-23T08:39:18.000Z
			false,
		},
		{
			"Weekend",
			`- name: Office hours
  cronExpression: "0 9 * * 1-5"
  duration: 8h`,
			1730027418, // 2024-10-27T11:10:18.000Z
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tws, err := NewTimeWindows(test.config)
			if err != nil {
				t.Errorf("unexpected error when creating time windows: %s", err)
			}
			if len(tws) != 1 {
				t.Errorf("expected to have exactly 1 time window, got %d", len(tws))
			}
			isTimeIn := tws[0].isTimeIn(time.Unix(test.timestamp, 0))
			if test.expected != isTimeIn {
				t.Errorf("unexpected result for test %s: expected %t but got %t", test.name, test.expected, isTimeIn)
			}
		})
	}
}
