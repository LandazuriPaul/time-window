package main

import (
	"fmt"
	"github.com/landazuripaul/time-window-validator/validator"
)

func main() {
	v, err := validator.NewValidator()
	if err != nil {
		validator.Result{
			Error:   fmt.Sprintf("parsing the flag inputs: %s", err),
			Message: "invalid flag inputs",
		}.WriteToGithubOutput()
		return
	}

	if v.IsForceAllowed() {
		validator.Result{
			IsAllowed: true,
			Message:   fmt.Sprintf("forced allowed with regexp %s and commit message %s", v.ForceAllowRegexp, v.CommitMessage),
			Timestamp: v.Timestamp.Unix(),
		}.WriteToGithubOutput()
		return
	}

	r := v.ValidateTimestamp()
	r.WriteToGithubOutput()
}
