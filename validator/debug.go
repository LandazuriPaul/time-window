//go:build debug

package validator

import (
	"os"
)

// init allows us to pass multiline arguments when debugging the CLI
func init() {
	os.Args = append(os.Args,
		"-commit_message",
		"my nice commit message",
		"-force_valid_regexp",
		"force\\-time\\-window",
		"-timestamp",
		"1729537200", // Monday, October 21, 2024 at 08:00:00 PM UTC
		"-allowed",
		`
- name: Office hours
  cronExpression: "0 9 * * 1-5"
  duration: 8h
  timezone: Europe/London`,
		"-blocked",
		`
- name: Christmas holidays
  cronExpression: "0 0 24 12 *"
  duration: 96h
  timezone: Europe/London`,
	)
}
