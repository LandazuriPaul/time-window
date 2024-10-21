package validator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const githubOutputEnvVar = "GITHUB_OUTPUT"

type Result struct {
	Error     string `json:"error"`
	IsAllowed bool   `json:"isAllowed"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (r Result) GithubFormat() string {
	return fmt.Sprintf("error=%s\nis_allowed=%t\nmessage=%s\ntimestamp=%d\nresult=%s\n", r.Error, r.IsAllowed, r.Message, r.Timestamp, r)
}

func (r Result) Print() {
	if r.Error == "" {
		fmt.Printf("%s\n", r)
	} else {
		log.Fatalf("%s\n", r)
	}
}

func (r Result) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		log.Printf("marshalling the result to JSON: %s", err)
		return fmt.Sprintf("result={\"error\":\"%s\",\"isAllowed\":%t,\"message\":\"%s\",\"timestamp\":%d,\"marshallingError\":\"%s\"}", r.Error, r.IsAllowed, r.Message, r.Timestamp, err)
	}
	return string(b)
}

func (r Result) WriteToGithubOutput() {
	outputFilePath := os.Getenv(githubOutputEnvVar)
	if outputFilePath != "" {
		// If the file doesn't exist, create it, or append to the file
		f, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("appending to $%s at %s: %s", githubOutputEnvVar, outputFilePath, err)
		}
		if _, err := f.Write([]byte(r.GithubFormat())); err != nil {
			f.Close() // ignore error; Write error takes precedence
			log.Fatalf("writing to the $%s at %s: %s", githubOutputEnvVar, outputFilePath, err)
		}
		if err := f.Close(); err != nil {
			log.Fatalf("closing file $%s at %s: %s", githubOutputEnvVar, outputFilePath, err)
		}
	}

	r.Print()
}
