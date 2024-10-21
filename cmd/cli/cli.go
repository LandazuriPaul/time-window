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

	if v.IsForceValidated() {
		validator.Result{
			IsValid:   true,
			Message:   fmt.Sprintf("forced validated with regexp %s and commit message %s", v.ForceValidRegexp, v.CommitMessage),
			Timestamp: v.Timestamp.Unix(),
		}.WriteToGithubOutput()
		return
	}

	r := v.ValidateTimestamp()
	r.WriteToGithubOutput()
}
