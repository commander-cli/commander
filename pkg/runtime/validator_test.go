package runtime

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func Test_NewValidator_Text(t *testing.T) {
    got := NewValidator(Text)
    assert.IsType(t, TextValidator{}, got)
}

func Test_NewValidator_Contains(t *testing.T) {
    got := NewValidator(Contains)
    assert.IsType(t, ContainsValidator{}, got)
}

func Test_NewValidator_Equal(t *testing.T) {
    got := NewValidator(Equal)
    assert.IsType(t, EqualValidator{}, got)
}

func TestTextValidator_Validate(t *testing.T) {
    v := TextValidator{}
    got := v.Validate("test", "test")
    assert.True(t, got.Success)
}

func TestTextValidator_ValidateFails(t *testing.T) {
    v := TextValidator{}
    got := v.Validate("test", "unequal")
    assert.False(t, got.Success)
    assert.Contains(t, got.Diff, "+unequal")
    assert.Contains(t, got.Diff, "-test")
}

func TestEqualValidator_Validate(t *testing.T) {
    v := EqualValidator{}
    got := v.Validate(1, 1)
    assert.True(t, got.Success)
}

func TestEqualValidator_ValidateFails(t *testing.T) {
    v := EqualValidator{}
    got := v.Validate(1, 0)
    assert.False(t, got.Success)
    assert.Contains(t, got.Diff, "+0")
    assert.Contains(t, got.Diff, "-1")
}

