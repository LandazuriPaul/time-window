package validator

import (
	"fmt"
	"regexp"
	"time"
)

type Validator struct {
	Allowed          []*TimeWindow
	Blocked          []*TimeWindow
	CommitMessage    string
	ForceValidRegexp *regexp.Regexp
}

func NewValidator(allowedInput, blockedInput, forceValidRegexpInput, commitMessageInput *string) (*Validator, error) {
	// init config fields
	commitMessage := ""
	var (
		err              error
		allowed          []*TimeWindow
		blocked          []*TimeWindow
		forceValidRegexp *regexp.Regexp
	)

	// allowed
	if allowedInput != nil && *allowedInput != "" {
		allowed, err = NewTimeWindows(*allowedInput)
		if err != nil {
			return nil, fmt.Errorf("parsing the allowed windows: %w", err)
		}
	}

	// blocked
	if blockedInput != nil && *blockedInput != "" {
		blocked, err = NewTimeWindows(*blockedInput)
		if err != nil {
			return nil, fmt.Errorf("parsing the blocked windows: %w", err)
		}
	}

	// forceValidRegexp
	if forceValidRegexpInput != nil && *forceValidRegexpInput != "" {
		forceValidRegexp, err = regexp.Compile(*forceValidRegexpInput)
		if err != nil {
			return nil, fmt.Errorf("compiling the force-valid-regexp: %w", err)
		}
	}

	// commit Message
	if commitMessageInput != nil && *commitMessageInput != "" {
		commitMessage = *commitMessageInput
	}

	return &Validator{
		Allowed:          allowed,
		Blocked:          blocked,
		CommitMessage:    commitMessage,
		ForceValidRegexp: forceValidRegexp,
	}, nil
}

func (v *Validator) Validate(timestamp time.Time) Result {
	if v.isForceValidated() {
		return Result{
			IsValid:   true,
			Message:   v.forceValidatedMessage(),
			Timestamp: timestamp.Unix(),
		}
	}
	return v.validateTimestamp(timestamp)
}

func (v *Validator) forceValidatedMessage() string {
	return fmt.Sprintf("forced validated with regexp %s and commit message %s", v.ForceValidRegexp, v.CommitMessage)
}

func (v *Validator) isForceValidated() bool {
	if v.ForceValidRegexp != nil && v.CommitMessage != "" {
		found := v.ForceValidRegexp.Find([]byte(v.CommitMessage))
		if found != nil {
			return true
		}
	}
	return false
}

func (v *Validator) validateTimestamp(timestamp time.Time) Result {
	// init result with timestamp
	r := Result{
		IsValid:   false,
		Timestamp: timestamp.Unix(),
	}

	// requires at least the allowed input
	if v.Allowed == nil {
		r.Error = "the required `allowed` input is missing"
		r.Message = "Missing `allowed` input. This can be due to a misconfiguration"
		return r
	}

	// check if the timestamp is in any time window
	// we stop at the first matching window
	var isInAllowed, isInBlocked bool
	var allowedWindowName, blockedWindowName string
	for _, tw := range v.Allowed {
		if tw.isTimeIn(timestamp) {
			isInAllowed = true
			allowedWindowName = tw.Name
			break
		}
	}
	for _, tw := range v.Blocked {
		if tw.isTimeIn(timestamp) {
			isInBlocked = true
			blockedWindowName = tw.Name
			break
		}
	}

	// blocked window takes precedence
	if isInBlocked {
		r.Message = fmt.Sprintf("timestamp in blocked window '%s'", blockedWindowName)
	} else {
		if isInAllowed {
			// timestamp is only allowed if it is both:
			//- NOT IN a blocked window
			//- IN an allowed window
			r.IsValid = true
			r.Message = fmt.Sprintf("timestamp in allowed window '%s'", allowedWindowName)
		} else {
			r.Message = "the timestamp isn't in any allowed window"
		}
	}

	return r
}
