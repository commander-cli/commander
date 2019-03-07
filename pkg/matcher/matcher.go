package matcher

import (
    "fmt"
    "github.com/pmezard/go-difflib/difflib"
    "strings"
)

const (
    Text     = "text"
    Contains = "contains"
    Equal    = "equal"
)

func NewMatcher(matcher string) Matcher {
    switch matcher {
    case Text:
        return TextMatcher{}
    case Contains:
        return ContainsMatcher{}
    case Equal:
        return EqualMatcher{}
    default:
        panic(fmt.Sprintf("Validator '%s' does not exist!", matcher))
    }
}

type Matcher interface {
    Match(actual interface{}, expected interface{}) MatcherResult
}

type MatcherResult struct {
    Success bool
    Diff    string
}

type TextMatcher struct {
}

func (m TextMatcher) Match(got interface{}, expected interface{}) MatcherResult {
    result := true
    if got != expected {
        result = false
    }

    diff := difflib.UnifiedDiff{
        A: difflib.SplitLines(got.(string)),
        B: difflib.SplitLines(expected.(string)),
        FromFile: "Got",
        ToFile: "Expected",
        Context: 3,
    }
    diffText, _ := difflib.GetUnifiedDiffString(diff)

    return MatcherResult{
        Diff:       diffText,
        Success:    result,
    }
}

type ContainsMatcher struct {
}

func (m ContainsMatcher) Match(got interface{}, expected interface{}) MatcherResult {
    result := true
    if !strings.Contains(got.(string), expected.(string)) {
        result = false
    }

    diff := difflib.UnifiedDiff{
        A: difflib.SplitLines(fmt.Sprintf("%s", got.(string))),
        B: difflib.SplitLines(fmt.Sprintf("%s", expected.(string))),
        FromFile: "Got",
        ToFile: "Expected",
        Context: 3,
    }
    diffText, _ := difflib.GetUnifiedDiffString(diff)

    return MatcherResult{
        Success: result,
        Diff: diffText,
    }
}

type EqualMatcher struct {
}

func (m EqualMatcher) Match(got interface{}, expected interface{}) MatcherResult {
    if got == expected {
        return MatcherResult{
            Success: true,
        }
    }

    diff := difflib.UnifiedDiff{
        A: difflib.SplitLines(fmt.Sprintf("%d", got.(int))),
        B: difflib.SplitLines(fmt.Sprintf("%d", expected.(int))),
        FromFile: "Got",
        ToFile: "Expected",
        Context: 3,
    }
    diffText, _ := difflib.GetUnifiedDiffString(diff)

    return MatcherResult{
        Success: false,
        Diff: diffText,
    }
}
