package debug

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextFormatterFieldSpawnMultipleEnabled(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo*,bar*")

	foo := Debug("foo").Spawn("child").Spawn("grandChild").WithFields(
		Fields{"field": 1, "field2": "two", "field3": "multiple strings"})
	foo.Log("foo")
	foo.Log(func() string { return "foo lazy" })

	bar := Debug("bar").Spawn("child").Spawn("grandChild")
	bar.Log("bar")
	bar.Log(func() string { return "bar lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
	assert.Contains(t, str, "foo:child:grandChild")
	assert.Contains(t, str, "field=1")
	assert.Contains(t, str, "field2=two")
	assert.Contains(t, str, `field3="multiple strings"`)
	assert.Contains(t, str, "foo")
	assert.Contains(t, str, "foo lazy")
	assert.Contains(t, str, "bar:child:grandChild")
	assert.Contains(t, str, "bar")
	assert.Contains(t, str, "bar lazy")
}

func TestBasicTextFormatterFieldsStrict(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo*")

	//nolint
	reg, err := regexp.Compile(`\d\d:\d\d:\d\d\.\d\d\d\s\d{1,4}(s|ms|us|ns)\s{1,3}.*foo.*-.*hello.*\n.*field=1 field2=two field3="multiple strings"`)

	assert.Nil(t, err, "regex error")

	foo := Debug("foo").WithFields(
		Fields{"field": 1, "field2": "two", "field3": "multiple strings"})
	foo.Log("hello")

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()

	assert.Contains(t, str, "foo", "namespace")
	assert.Contains(t, str, "hello", "message")
	assert.Contains(t, str, "field=1")
	assert.Contains(t, str, "field2=two")
	assert.Contains(t, str, `field3="multiple strings"`)

	assert.Regexp(t, reg, str, "strict match")
}
