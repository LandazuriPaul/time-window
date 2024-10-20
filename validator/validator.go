package validator

import (
	"flag"
	"fmt"
	"regexp"
	"time"
)

type Validator struct {
	Allowed          []TimeWindow
	Blocked          []TimeWindow
	CommitMessage    string
	ForceAllowRegexp *regexp.Regexp
	Timestamp        time.Time
}

func NewValidator() (*Validator, error) {
	// parse flagged inputs
	allowedFlag := flag.String("allowed", "", "a YAML string representing the allowed windows")
	blockedFlag := flag.String("blocked", "", "a YAML string representing the blocked windows")
	forceAllowRegexpFlag := flag.String("force_allow_regexp", "", "a regular expression to check the commit message to bypass time window checks")
	commitMessageFlag := flag.String("commit_message", "", "the commit message")
	timestampFlag := flag.Int64("timestamp", 0, "Unix timestamp used to check the time windows")
	flag.Parse()

	// init config fields
	commitMessage := ""
	timestamp := time.Now()
	var (
		err              error
		allowed          []TimeWindow
		blocked          []TimeWindow
		forceAllowRegexp *regexp.Regexp
	)

	// allowed
	if allowedFlag != nil && *allowedFlag != "" {
		allowed, err = parseTimeWindows(*allowedFlag)
		if err != nil {
			return nil, fmt.Errorf("parsing the allowed windows")
		}
	}

	// blocked
	if blockedFlag != nil && *blockedFlag != "" {
		blocked, err = parseTimeWindows(*blockedFlag)
		if err != nil {
			return nil, fmt.Errorf("parsing the blocked windows: %w", err)
		}
	}

	// forceAllowRegexp
	if forceAllowRegexpFlag != nil && *forceAllowRegexpFlag != "" {
		forceAllowRegexp, err = regexp.Compile(*forceAllowRegexpFlag)
		if err != nil {
			return nil, fmt.Errorf("compiling the force-allow-regexp: %w", err)
		}
	}

	// commit Message
	if commitMessageFlag != nil && *commitMessageFlag != "" {
		commitMessage = *commitMessageFlag
	}

	// timestamp
	if timestampFlag != nil && *timestampFlag != 0 {
		timestamp = time.Unix(*timestampFlag, 0)
	}

	return &Validator{
		Allowed:          allowed,
		Blocked:          blocked,
		CommitMessage:    commitMessage,
		ForceAllowRegexp: forceAllowRegexp,
		Timestamp:        timestamp,
	}, nil
}

func (v *Validator) IsForceAllowed() bool {
	if v.ForceAllowRegexp != nil && v.CommitMessage != "" {
		found := v.ForceAllowRegexp.Find([]byte(v.CommitMessage))
		if found != nil {
			return true
		}
	}
	return false
}

func (v *Validator) ValidateTimestamp() Result {
	// init result with timestamp
	r := Result{
		Timestamp: v.Timestamp.Unix(),
	}

	// requires at least one of allowed or blocked
	if v.Allowed == nil && v.Blocked == nil {
		r.Error = "at least one of `allowed` or `blocked` is required"
		r.Message = "Missing `allowed` and `blocked`. This can be due to a misconfiguration"
		return r
	}

	// TODO
	return r
}
