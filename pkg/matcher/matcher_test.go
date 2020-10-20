package matcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewMatcher_Text(t *testing.T) {
	got := NewMatcher(Text)
	assert.IsType(t, TextMatcher{}, got)
}

func Test_NewMatcher_Contains(t *testing.T) {
	got := NewMatcher(Contains)
	assert.IsType(t, ContainsMatcher{}, got)
}

func Test_NewMatcher_Equal(t *testing.T) {
	got := NewMatcher(Equal)
	assert.IsType(t, EqualMatcher{}, got)
}

func TestTextMatcher_Validate(t *testing.T) {
	m := TextMatcher{}
	got := m.Match("test", "test")
	assert.True(t, got.Success)
}

func TestNewMatcher_Fail(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			assert.Contains(t, r, "Validator 'no' does not exist!")
		}
		assert.NotNil(t, r)
	}()

	_ = NewMatcher("no")
}

func TestTextMatcher_ValidateFails(t *testing.T) {
	m := TextMatcher{}
	got := m.Match("test", "unequal")
	assert.False(t, got.Success)
	assert.Contains(t, got.Diff, "+unequal")
	assert.Contains(t, got.Diff, "-test")
}

func TestEqualMatcher_Validate(t *testing.T) {
	m := EqualMatcher{}
	got := m.Match(2, 2)
	assert.True(t, got.Success)
}

func TestEqualMatcher_ValidateFails(t *testing.T) {
	m := EqualMatcher{}
	got := m.Match(2, 3)
	assert.False(t, got.Success)
	assert.Contains(t, got.Diff, "+3")
	assert.Contains(t, got.Diff, "-2")
}

func TestContainsMatcher_Match(t *testing.T) {
	m := ContainsMatcher{}
	got := m.Match("hello world", "hello")
	assert.True(t, got.Success)
}

func TestContainsMatcher_MatchFail(t *testing.T) {
	m := ContainsMatcher{}
	got := m.Match("hello world", "nope")
	assert.False(t, got.Success)
}

func TestNewMatcher_NotContains(t *testing.T) {
	m := NewMatcher(NotContains)
	assert.IsType(t, NotContainsMatcher{}, m)
}

func TestNotContainsMatcher(t *testing.T) {
	m := NotContainsMatcher{}
	r := m.Match("lore ipsum donor", "hello")
	assert.True(t, r.Success)
	assert.Equal(t, "\nExpected\n\nlore ipsum donor\n\nto not contain\n\nhello\n", r.Diff)
}

func TestNotContainsMatcher_Fails(t *testing.T) {
	m := NotContainsMatcher{}
	r := m.Match("lore ipsum donor", "donor")
	assert.False(t, r.Success)

	diffText := "\nExpected\n\nlore ipsum donor\n\nto not contain\n\ndonor\n"
	assert.Equal(t, diffText, r.Diff)
}

func TestXMLMatcher_Match(t *testing.T) {
	m := XMLMatcher{}
	r := m.Match("<book>test</book>", map[string]string{"/book": "test"})

	assert.True(t, r.Success)
	assert.Equal(t, "", r.Diff)
}

func TestXMLMatcher_DoesNotMatch(t *testing.T) {
	m := XMLMatcher{}
	r := m.Match("<book>another</book>", map[string]string{"/book": "test"})

	diff := `Expected xml path "/book" with result

test

to be equal to

another`
	assert.False(t, r.Success)
	assert.Equal(t, diff, r.Diff)
}

func TestXMLMatcher_DoesNotFindPath(t *testing.T) {
	m := XMLMatcher{}
	r := m.Match("<comic>test</comic>", map[string]string{"/book": "test"})

	assert.False(t, r.Success)
	assert.Equal(t, `Query "/book" did not match a path`, r.Diff)
}

func TestXMLMatcher_MatchArray(t *testing.T) {
	m := NewMatcher(XML)
	html := "<books><book>test1</book><book>test2</book></books>"

	r := m.Match(html, map[string]string{"/books": "test1test2"})

	assert.True(t, r.Success)
	assert.Equal(t, "", r.Diff)
}

func TestXMLMatcher_WithXPATHError(t *testing.T) {
	m := XMLMatcher{}
	r := m.Match("<book>test</book>", map[string]string{"<!/book": "test"})

	assert.False(t, r.Success)
	assert.Equal(t, `Error occured: expression must evaluate to a node-set`, r.Diff)
}

func TestJSONMatcher_Match(t *testing.T) {
	m := JSONMatcher{}
	r := m.Match(`{"book": "test"}`, map[string]string{"book": "test"})

	assert.True(t, r.Success)
	assert.Equal(t, "", r.Diff)
}

func TestJSONMatcher_DoesNotFindPath(t *testing.T) {
	m := JSONMatcher{}
	r := m.Match(`{"comic": "test"}`, map[string]string{"book": "test"})

	assert.False(t, r.Success)
	assert.Equal(t, `Query "book" did not match a path`, r.Diff)
}

func TestJSONMatcher_MatchArray(t *testing.T) {
	m := JSONMatcher{}
	json := `{"books": [ "test1", "test2" ]}`

	r := m.Match(json, map[string]string{"books": "[test1 test2]"})

	assert.True(t, r.Success)
	assert.Equal(t, "", r.Diff)
}

func TestJSONMatcher_DoesNotMatch(t *testing.T) {
	m := NewMatcher(JSON)
	r := m.Match(`{"book": "another"}`, map[string]string{"book": "test"})

	diff := `Expected json path "book" with result

test

to be equal to

another`

	assert.False(t, r.Success)
	assert.Equal(t, diff, r.Diff)
}

func TestFileMatcher_Validate(t *testing.T) {
	m := FileMatcher{}
	got := m.Match("zoom", "zoom.txt")
	assert.True(t, got.Success)
	assert.Equal(t, "", got.Diff)
}

func TestFileMatcher_ValidateFails(t *testing.T) {
	m := FileMatcher{}
	got := m.Match("zoom", "zoologist.txt")
	assert.False(t, got.Success)
	println(got.Diff)
	assert.Contains(t, got.Diff, "+zoologist")
	assert.Contains(t, got.Diff, "-zoom")
}
