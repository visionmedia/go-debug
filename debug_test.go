package debug

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func assertContains(t *testing.T, str, substr string) {
	if !strings.Contains(str, substr) {
		t.Fatalf("expected %q to contain %q", str, substr)
	}
}

func assertNotContains(t *testing.T, str, substr string) {
	if strings.Contains(str, substr) {
		t.Fatalf("expected %q to not contain %q", str, substr)
	}
}

func TestDefault(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	debug := Debug("foo")
	debug.Debug("something")
	debug.Debug("here")
	debug.Debug("whoop")

	if buf.Len() != 0 {
		t.Fatalf("buffer should be empty")
	}
}

func TestDefaultLazy(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	debug := Debug("foo")
	debug.Debug(func() string { return "something" })
	debug.Debug(func() string { return "here" })
	debug.Debug(func() string { return "whoop" })

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
	debug.Debug("something")
	debug.Debug("here")
	debug.Debug("whoop")
	debug.Debug(func() string { return "lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := string(buf.Bytes())
	assertContains(t, str, "something")
	assertContains(t, str, "here")
	assertContains(t, str, "whoop")
	assertContains(t, str, "lazy")
}

func TestMultipleOneEnabled(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo")

	foo := Debug("foo")
	foo.Debug("foo")
	foo.Debug(func() string { return "foo lazy" })

	bar := Debug("bar")
	bar.Debug("bar")
	bar.Debug(func() string { return "bar lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := string(buf.Bytes())
	assertContains(t, str, "foo")
	assertContains(t, str, "foo lazy")
	assertNotContains(t, str, "bar")
	assertNotContains(t, str, "bar lazy")
}

func TestMultipleEnabled(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo,bar")

	foo := Debug("foo")
	foo.Debug("foo")
	foo.Debug(func() string { return "foo lazy" })

	bar := Debug("bar")
	bar.Debug("bar")
	bar.Debug(func() string { return "bar lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := string(buf.Bytes())
	assertContains(t, str, "foo")
	assertContains(t, str, "foo lazy")
	assertContains(t, str, "bar")
	assertContains(t, str, "bar lazy")
}

func TestEnableDisable(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo,bar")
	Disable()

	foo := Debug("foo")
	foo.Debug("foo")
	foo.Debug(func() string { return "foo" })

	bar := Debug("bar")
	bar.Debug("bar")
	bar.Debug(func() string { return "bar" })

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
		debug.Debug("sending mail")
		debug.Debug("send email to %s", "tobi@segment.io")
		debug.Debug("send email to %s", "loki@segment.io")
		debug.Debug("send email to %s", "jane@segment.io")
		time.Sleep(500 * time.Millisecond)
	}
}

func BenchmarkDisabled(b *testing.B) {
	debug := Debug("something")
	for i := 0; i < b.N; i++ {
		debug.Debug("stuff")
	}
}

func BenchmarkDisabledLazy(b *testing.B) {
	debug := Debug("something")
	for i := 0; i < b.N; i++ {
		debug.Debug(func() string { return "lazy" })
	}
}

func BenchmarkNonMatch(b *testing.B) {
	debug := Debug("something")
	Enable("nonmatch")
	for i := 0; i < b.N; i++ {
		debug.Debug("stuff")
	}
}

func BenchmarkLargeNonMatch(b *testing.B) {
	debug := Debug("large:not:lazy")

	abs, _ := filepath.Abs("./crashes.json")
	file := GetFileBytes(abs)

	Enable("nonmatch")
	for i := 0; i < b.N; i++ {
		debug.Debug(string(file))
	}
}

func BenchmarkLargeLazyNonMatch(b *testing.B) {
	debug := Debug("large:lazy")

	abs, _ := filepath.Abs("./crashes.json")
	file := GetFileBytes(abs)

	Enable("nonmatch")
	for i := 0; i < b.N; i++ {
		debug.Debug(func() string {
			return string(file)
		})
	}
}

func BenchmarkLargeLazyMatch(b *testing.B) {
	debug := Debug("large:lazy")

	abs, _ := filepath.Abs("./crashes.json")
	file := GetFileBytes(abs)

	Enable("large:lazy")
	for i := 0; i < b.N; i++ {
		debug.Debug(func() string {
			return string(file)
		})
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
