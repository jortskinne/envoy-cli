package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TypeRule struct {
	Key     string
	Type    string // "int", "bool", "url", "email", "regex"
	Pattern string // used when Type == "regex"
}

type TypeCheckOptions struct {
	Rules []TypeRule
}

type TypeCheckIssue struct {
	Key      string
	Value    string
	Expected string
	Message  string
}

func DefaultTypeCheckOptions() TypeCheckOptions {
	return TypeCheckOptions{}
}

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
var urlRe = regexp.MustCompile(`^https?://[^\s]+$`)

func TypeCheck(entries []Entry, opts TypeCheckOptions) []TypeCheckIssue {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	var issues []TypeCheckIssue
	for _, rule := range opts.Rules {
		val, ok := lookup[rule.Key]
		if !ok {
			continue
		}
		switch strings.ToLower(rule.Type) {
		case "int":
			if _, err := strconv.Atoi(val); err != nil {
				issues = append(issues, TypeCheckIssue{Key: rule.Key, Value: val, Expected: "int", Message: "value is not a valid integer"})
			}
		case "bool":
			lower := strings.ToLower(val)
			if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
				issues = append(issues, TypeCheckIssue{Key: rule.Key, Value: val, Expected: "bool", Message: "value is not a valid boolean"})
			}
		case "url":
			if !urlRe.MatchString(val) {
				issues = append(issues, TypeCheckIssue{Key: rule.Key, Value: val, Expected: "url", Message: "value is not a valid URL"})
			}
		case "email":
			if !emailRe.MatchString(val) {
				issues = append(issues, TypeCheckIssue{Key: rule.Key, Value: val, Expected: "email", Message: "value is not a valid email"})
			}
		case "regex":
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				issues = append(issues, TypeCheckIssue{Key: rule.Key, Value: val, Expected: "regex", Message: fmt.Sprintf("invalid pattern: %s", err)})
				continue
			}
			if !re.MatchString(val) {
				issues = append(issues, TypeCheckIssue{Key: rule.Key, Value: val, Expected: rule.Pattern, Message: "value does not match required pattern"})
			}
		}
	}
	return issues
}
