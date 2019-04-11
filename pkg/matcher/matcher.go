package matcher

import (
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"strings"
)

const (
	// Text matcher type
	Text = "text"
	// Contains matcher type
	Contains = "contains"
	// Equal matcher type
	Equal = "equal"
	// Not contains matcher type
	NotContains = "notcontains"
)

// NewMatcher creates a new matcher by type
func NewMatcher(matcher string) Matcher {
	switch matcher {
	case Text:
		return TextMatcher{}
	case Contains:
		return ContainsMatcher{}
	case Equal:
		return EqualMatcher{}
	case NotContains:
		return NotContainsMatcher{}
	default:
		panic(fmt.Sprintf("Validator '%s' does not exist!", matcher))
	}
}

// Matcher interface which is implemented by all matchers
type Matcher interface {
	Match(actual interface{}, expected interface{}) MatcherResult
}

// MatcherResult will be returned by the executed matcher
type MatcherResult struct {
	Success bool
	Diff    string
}

// TextMatcher matches if two texts are equals
type TextMatcher struct {
}

// Match matches both texts
func (m TextMatcher) Match(got interface{}, expected interface{}) MatcherResult {
	result := true
	if got != expected {
		result = false
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(got.(string)),
		B:        difflib.SplitLines(expected.(string)),
		FromFile: "Got",
		ToFile:   "Expected",
		Context:  3,
	}
	diffText, _ := difflib.GetUnifiedDiffString(diff)

	return MatcherResult{
		Diff:    diffText,
		Success: result,
	}
}

// ContainsMatcher tests if the expected value is in the got variable
type ContainsMatcher struct {
}

// Match matches on the given values
func (m ContainsMatcher) Match(got interface{}, expected interface{}) MatcherResult {
	result := true
	if !strings.Contains(got.(string), expected.(string)) {
		result = false
	}

	diff := `
Expected

%s

to contain

%s
`
	diff = fmt.Sprintf(diff, got, expected)

	return MatcherResult{
		Success: result,
		Diff:    diff,
	}
}

//EqualMatcher matches if given values are equal
type EqualMatcher struct {
}

//Match matches the values if they are equal
func (m EqualMatcher) Match(got interface{}, expected interface{}) MatcherResult {
	if got == expected {
		return MatcherResult{
			Success: true,
		}
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(fmt.Sprintf("%v", got)),
		B:        difflib.SplitLines(fmt.Sprintf("%v", expected)),
		FromFile: "Got",
		ToFile:   "Expected",
		Context:  3,
	}
	diffText, _ := difflib.GetUnifiedDiffString(diff)

	return MatcherResult{
		Success: false,
		Diff:    diffText,
	}
}

type NotContainsMatcher struct {
}

func (m NotContainsMatcher) Match(got interface{}, expected interface{}) MatcherResult {
	result := true
	if strings.Contains(got.(string), expected.(string)) {
		result = false
	}

	diff := `
Expected

%s

to not contain

%s
`
	diff = fmt.Sprintf(diff, got, expected)

	return MatcherResult{
		Success: result,
		Diff:    diff,
	}
}
