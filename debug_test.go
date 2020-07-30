package debug

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	colors = []string{"31"}
	os.Exit(m.Run())
}

func TestDefault(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	debug := Debug("foo")
	debug.Log("something")
	debug.Log("here")
	debug.Log("whoop")
	debug.Log(os.Args[:1]) // can log non strings

	if buf.Len() != 0 {
		t.Fatalf("buffer should be empty")
	}
}

func TestDefaultLazy(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	debug := Debug("foo")
	debug.Log(func() string { return "something" })
	debug.Log(func() string { return "here" })
	debug.Log(func() string { return "whoop" })

	if buf.Len() != 0 {
		t.Fatalf("buffer should be empty")
	}
}

func TestEnable(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo")

	debug := Debug("foo")
	// assert.Equal(t, "foo", debug.name)
	debug.Log("something")
	debug.Log("here")
	debug.Log("whoop")
	debug.Log(func() string { return "lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
	assert.Contains(t, str, "something")
	assert.Contains(t, str, "here")
	assert.Contains(t, str, "whoop")
	assert.Contains(t, str, "lazy")
}

func TestEnableNegative(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("*,-foo")

	debug := Debug("foo")
	// assert.Equal(t, "foo", debug.name)
	debug.Log("something")
	debug.Log("here")
	debug.Log("whoop")
	debug.Log(func() string { return "lazy" })

	assert.Equal(t, buf.Len(), 0)
}

func TestColorsEnable(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo")

	debug := Debug("foo")
	debug.Log("something")

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
	assert.Contains(t, str, "something")
	assert.Contains(t, str, getColorStr(colors[0], true))
	assert.Contains(t, str, "\033")
}

func TestColorsDisable(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo")
	SetHasColors(false)

	debug := Debug("foo")
	debug.Log("something")

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
	assert.Contains(t, str, "something")
	assert.NotContains(t, str, getColorStr(colors[0], true))

	SetHasColors(true)
}

func TestMultipleOneEnabled(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo")

	foo := Debug("foo")
	foo.Log("foo")
	foo.Log(func() string { return "foo lazy" })

	bar := Debug("bar")
	bar.Log("bar")
	bar.Log(func() string { return "bar lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
	assert.Contains(t, str, "foo")
	assert.Contains(t, str, "foo lazy")
	assert.NotContains(t, str, "bar")
	assert.NotContains(t, str, "bar lazy")
}

func TestMultipleEnabled(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo,bar")

	foo := Debug("foo")
	foo.Log("foo")
	foo.Log(func() string { return "foo lazy" })

	bar := Debug("bar")
	bar.Log("bar")
	bar.Log(func() string { return "bar lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
	assert.Contains(t, str, "foo")
	assert.Contains(t, str, "foo lazy")
	assert.Contains(t, str, "bar")
	assert.Contains(t, str, "bar lazy")
}

func TestSpawnMultipleEnabled(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo*,bar*")

	//nolint
	var foo IDebugger

	foo = Debug("foo").Spawn("child").Spawn("grandChild")
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
	assert.Contains(t, str, "foo")
	assert.Contains(t, str, "foo lazy")
	assert.Contains(t, str, "bar:child:grandChild")
	assert.Contains(t, str, "bar")
	assert.Contains(t, str, "bar lazy")
}

func TestEnableDisable(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo,bar")
	Disable()

	foo := Debug("foo")
	foo.Log("foo")
	foo.Log(func() string { return "foo" })

	bar := Debug("bar")
	bar.Log("bar")
	bar.Log(func() string { return "bar" })

	if buf.Len() != 0 {
		t.Fatalf("buffer should not have output")
	}
}

func TestSpawnEnableDisable(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	cache.Flush()

	Enable("foo*,bar*")
	Disable()

	foo := Debug("foo").Spawn("child").Spawn("grandChild")
	foo.Log("foo")
	foo.Log(func() string { return "foo" })

	bar := Debug("bar").Spawn("child").Spawn("grandChild")
	bar.Log("bar")
	bar.Log(func() string { return "bar" })

	// run again to test cache to make sure it does not overflow
	Debug("bar").Spawn("child").Spawn("grandChild")

	// fmt.Println("items", cache.Items())
	assert.Equal(t, len(cache.Items()), 6)

	if buf.Len() != 0 {
		t.Fatalf("buffer should not have output")
	}
}

func ExampleEnable() {
	Enable("mongo:connection")
	Enable("mongo:*")
	Enable("foo,bar,baz")
	Enable("*")
}

func ExampleDebug() {
	var debug = Debug("single")

	for {
		debug.Log("sending mail")
		debug.Log("send email to %s", "tobi@segment.io")
		debug.Log("send email to %s", "loki@segment.io")
		debug.Log("send email to %s", "jane@segment.io")
		time.Sleep(500 * time.Millisecond)
	}
}

func GetFileBytes(filename string) []byte {
	file, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)

	return bytes
}

func TestBuildPattern(t *testing.T) {
	tests := []struct {
		in        string
		pattern   string
		negatives []string
	}{
		{
			in:      "a,b,c",
			pattern: RegExWrap("a|b|c"),
		},
		{
			in:        "a,-b,-c",
			pattern:   RegExWrap("a"),
			negatives: []string{"b", "c"},
		},
		{
			in:        "*,-b,-c",
			pattern:   RegExWrap(".*?"),
			negatives: []string{"b", "c"},
		},
		{
			in:        "parent:child*,-parent:child:a,parent:child:b,-parent:child:c",
			pattern:   RegExWrap(`parent:child.*?|parent:child:b`),
			negatives: []string{"parent:child:a", "parent:child:c"},
		},
	}

	for _, tt := range tests {
		tt := tt
		actual, actualNegs := BuildPattern(tt.in)

		assert.Equal(t, tt.pattern, actual, "pattern matches")
		assert.Equal(t, tt.negatives, actualNegs, "pattern matches")
	}
}
