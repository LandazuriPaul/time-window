package validator

import (
	"flag"
	"fmt"
	"regexp"
	"time"
)

type Validator struct {
	Allowed          []*TimeWindow
	Blocked          []*TimeWindow
	CommitMessage    string
	ForceValidRegexp *regexp.Regexp
	Timestamp        time.Time
}

func NewValidator() (*Validator, error) {
	// parse flagged inputs
	allowedFlag := flag.String("allowed", "", "a YAML string representing the allowed windows")
	blockedFlag := flag.String("blocked", "", "a YAML string representing the blocked windows")
	forceValidRegexpFlag := flag.String("force_valid_regexp", "", "a regular expression to check the commit message to bypass time window checks")
	commitMessageFlag := flag.String("commit_message", "", "the commit message")
	timestampFlag := flag.Int64("timestamp", 0, "Unix timestamp (in seconds) used to check the time windows")
	flag.Parse()

	// init config fields
	commitMessage := ""
	timestamp := time.Now() // FIXME: Timezone!
	var (
		err              error
		allowed          []*TimeWindow
		blocked          []*TimeWindow
		forceValidRegexp *regexp.Regexp
	)

	// allowed
	if allowedFlag != nil && *allowedFlag != "" {
		allowed, err = NewTimeWindows(*allowedFlag)
		if err != nil {
			return nil, fmt.Errorf("parsing the allowed windows: %w", err)
		}
	}

	// blocked
	if blockedFlag != nil && *blockedFlag != "" {
		blocked, err = NewTimeWindows(*blockedFlag)
		if err != nil {
			return nil, fmt.Errorf("parsing the blocked windows: %w", err)
		}
	}

	// forceValidRegexp
	if forceValidRegexpFlag != nil && *forceValidRegexpFlag != "" {
		forceValidRegexp, err = regexp.Compile(*forceValidRegexpFlag)
		if err != nil {
			return nil, fmt.Errorf("compiling the force-valid-regexp: %w", err)
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
		ForceValidRegexp: forceValidRegexp,
		Timestamp:        timestamp,
	}, nil
}

func (v *Validator) IsForceValidated() bool {
	if v.ForceValidRegexp != nil && v.CommitMessage != "" {
		found := v.ForceValidRegexp.Find([]byte(v.CommitMessage))
		if found != nil {
			return true
		}
	}
	return false
}

func (v *Validator) ValidateTimestamp() Result {
	// init result with timestamp
	r := Result{
		IsValid:   false,
		Timestamp: v.Timestamp.Unix(),
	}

	// requires at least one of allowed or blocked
	if v.Allowed == nil && v.Blocked == nil {
		r.Error = "at least one of `allowed` or `blocked` is required"
		r.Message = "Missing `allowed` and `blocked`. This can be due to a misconfiguration"
		return r
	}

	// check if the timestamp is in any time window
	// we stop at the first matching window
	var isInAllowed, isInBlocked bool
	var allowedWindowName, blockedWindowName string
	for _, tw := range v.Allowed {
		if tw.isTimeIn(v.Timestamp) {
			isInAllowed = true
			allowedWindowName = tw.Name
			break
		}
	}
	for _, tw := range v.Blocked {
		if tw.isTimeIn(v.Timestamp) {
			isInBlocked = true
			blockedWindowName = tw.Name
			break
		}
	}

	// blocked window takes precedence
	if isInBlocked {
		r.Message = fmt.Sprintf("timestamp in blocked window %s", blockedWindowName)
	} else {
		if isInAllowed {
			// timestamp is only allowed if it is both:
			//- NOT IN a blocked window
			//- IN an allowed window
			r.IsValid = true
			r.Message = fmt.Sprintf("timestamp in allowed window %s", allowedWindowName)
		} else {
			r.Message = "the timestamp isn't in any allowed window"
		}
	}

	return r
}