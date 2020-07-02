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

	"github.com/stretchr/testify/assert"
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
	debug.Log("something")
	debug.Log("here")
	debug.Log("whoop")
	debug.Log(func() string { return "lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
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
	foo.Log("foo")
	foo.Log(func() string { return "foo lazy" })

	bar := Debug("bar")
	bar.Log("bar")
	bar.Log(func() string { return "bar lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
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
	foo.Log("foo")
	foo.Log(func() string { return "foo lazy" })

	bar := Debug("bar")
	bar.Log("bar")
	bar.Log(func() string { return "bar lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
	assertContains(t, str, "foo")
	assertContains(t, str, "foo lazy")
	assertContains(t, str, "bar")
	assertContains(t, str, "bar lazy")
}

func TestSpawnMultipleEnabled(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	SetWriter(buf)

	Enable("foo*,bar*")

	foo := Debug("foo").Spawn("child").Spawn("grandChild")
	foo.Log("foo")
	foo.Log(func() string { return "foo lazy" })

	bar := Debug("bar").Spawn("child").Spawn("grandChild")
	bar.Log("bar")
	bar.Log(func() string { return "bar lazy" })

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := buf.String()
	assertContains(t, str, "foo:child:grandChild")
	assertContains(t, str, "foo")
	assertContains(t, str, "foo lazy")
	assertContains(t, str, "bar:child:grandChild")
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

func BenchmarkDisabled(b *testing.B) {
	debug := Debug("something")
	for i := 0; i < b.N; i++ {
		debug.Log("stuff")
	}
}

func BenchmarkMatch(b *testing.B) {
	debug := Debug("something")
	Enable("something")
	for i := 0; i < b.N; i++ {
		debug.Log("stuff")
	}
}

func BenchmarkMatchNonStringNonFunc(b *testing.B) {
	debug := Debug("something")
	Enable("something")
	for i := 0; i < b.N; i++ {
		debug.Log(os.Args[:1])
	}
}

func BenchmarkDisabledLazy(b *testing.B) {
	debug := Debug("something")
	for i := 0; i < b.N; i++ {
		debug.Log(func() string { return "lazy" })
	}
}

func BenchmarkNonMatch(b *testing.B) {
	debug := Debug("something")
	Enable("nonmatch")
	for i := 0; i < b.N; i++ {
		debug.Log("stuff")
	}
}

func BenchmarkLargeNonMatch(b *testing.B) {
	debug := Debug("large:not:lazy")

	abs, _ := filepath.Abs("./crashes.json")
	file := GetFileBytes(abs)

	Enable("nonmatch")
	for i := 0; i < b.N; i++ {
		debug.Log(string(file))
	}
}

func BenchmarkLargeLazyNonMatch(b *testing.B) {
	debug := Debug("large:lazy")

	abs, _ := filepath.Abs("./crashes.json")
	file := GetFileBytes(abs)

	Enable("nonmatch")
	for i := 0; i < b.N; i++ {
		debug.Log(func() string {
			return string(file)
		})
	}
}

func BenchmarkLargeMatch(b *testing.B) {
	debug := Debug("large:lazy")

	abs, _ := filepath.Abs("./crashes.json")
	file := GetFileBytes(abs)

	Enable("large:lazy")
	for i := 0; i < b.N; i++ {
		debug.Log(string(file))
	}
}

func BenchmarkLargeLazyMatch(b *testing.B) {
	debug := Debug("large:lazy")

	abs, _ := filepath.Abs("./crashes.json")
	file := GetFileBytes(abs)

	Enable("large:lazy")
	for i := 0; i < b.N; i++ {
		debug.Log(func() string {
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
