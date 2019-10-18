package matcher

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/tidwall/gjson"
	"log"
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
	JSON        = "json"
	XML         = "xml"
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
	case JSON:
		return JSONMatcher{}
	case XML:
		return XMLMatcher{}
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

type JSONMatcher struct {
}

func (m JSONMatcher) Match(got interface{}, expected interface{}) MatcherResult {
	result := MatcherResult{
		Success: true,
		Diff:    "",
	}

	json := expected.(map[string]string)
	for q, e := range json {
		r := gjson.Get(got.(string), q)
		if r.Value() == nil {
			result.Success = false
			result.Diff = fmt.Sprintf(`Query "%s" did not match a path`, q)
			break
		}

		if fmt.Sprintf("%v", r.Value()) != e {
			result.Success = false
			result.Diff = fmt.Sprintf(`Expected json path "%s" with result

%s

to be equal to

%s`, q, e, r.Value())
			break
		}
	}

	return result
}

type XMLMatcher struct {
}

func (m XMLMatcher) Match(got interface{}, expected interface{}) MatcherResult {
	result := MatcherResult{
		Success: true,
		Diff:    "",
	}

	xml := expected.(map[string]string)
	for q, e := range xml {
		doc, err := xmlquery.Parse(strings.NewReader(got.(string)))
		if err != nil {
			log.Fatal(err)
		}

		node, err := xmlquery.Query(doc, q)
		if err != nil {
		    return MatcherResult{
		        Success: false,
		        Diff: fmt.Sprintf("Error occured: %s", err),
            }
		}

		if node == nil {
            return MatcherResult{
                Success: false,
                Diff:    fmt.Sprintf(`Query "%s" did not match a path`, q),
            }
        }

        if node.InnerText() != e {
            result.Success = false
            result.Diff = fmt.Sprintf(`Expected xml path "%s" with result

%s

to be equal to

%s`, q, e, node.InnerText())
        }
    }

	return result
}
