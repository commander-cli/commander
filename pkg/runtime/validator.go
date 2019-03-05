package runtime

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

func NewValidator(validator string) Validator {
	switch validator {
	case Text:
		return TextValidator{}
	case Contains:
		return ContainsValidator{}
	case Equal:
		return EqualValidator{}
	default:
		panic(fmt.Sprintf("Validator '%s' does not exist!", validator))
	}
}

type Validator interface {
    Validate(got interface{}, expected interface{}) ValidationResult
}

type ValidationResult struct {
	Success    bool
	Diff       string
}

type TextValidator struct {
}

func (v TextValidator) Validate(got interface{}, expected interface{}) ValidationResult {
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

	return ValidationResult{
		Diff:       diffText,
		Success:    result,
	}
}

type ContainsValidator struct {
}

func (v ContainsValidator) Validate(got interface{}, expected interface{}) ValidationResult {
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

	return ValidationResult{
		Success: result,
		Diff: diffText,
	}
}

type EqualValidator struct {
}

func (v EqualValidator) Validate(got interface{}, expected interface{}) ValidationResult {
	if got == expected {
		return ValidationResult{
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

	return ValidationResult{
		Success: false,
		Diff: diffText,
	}
}
