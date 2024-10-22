package main

import (
	"flag"
	"fmt"
	"github.com/landazuripaul/time-window-validator/validator"
	"time"
)

const defaultForceValidRegexp = `force\-time\-window`

func main() {
	// parse flagged inputs
	allowedInput := flag.String("allowed", "", "a YAML string representing the allowed windows")
	blockedInput := flag.String("blocked", "", "a YAML string representing the blocked windows")
	forceValidRegexpInput := flag.String("force_valid_regexp", defaultForceValidRegexp, "a regular expression to check the commit message to bypass time window checks")
	commitMessageInput := flag.String("commit_message", "", "the commit message")
	timestampInput := flag.Int64("timestamp", 0, "Unix timestamp (in seconds) used to check the time windows")
	flag.Parse()

	v, err := validator.NewValidator(allowedInput, blockedInput, forceValidRegexpInput, commitMessageInput)
	if err != nil {
		validator.Result{
			Error:   fmt.Sprintf("parsing the flag inputs: %s", err),
			Message: "invalid flag inputs",
		}.WriteToGithubOutput()
		return
	}

	// timestamp
	// FIXME: Timezone!
	timestamp := time.Now()
	if timestampInput != nil && *timestampInput != 0 {
		timestamp = time.Unix(*timestampInput, 0)
	}

	r := v.Validate(timestamp)
	r.WriteToGithubOutput()
}
