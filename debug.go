package debug

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	writer  io.Writer = os.Stderr
	reg     *regexp.Regexp
	m       sync.Mutex
	enabled = false
)

// Debugger function.
type DebugFunction func(...interface{})

type Debugger struct {
	Log   DebugFunction
	Spawn func(ns string) Debugger
}

// Terminal colors used at random.
var colors []string = []string{
	"31",
	"32",
	"33",
	"34",
	"35",
	"36",
}

// Initialize with DEBUG environment variable.
func init() {
	env := os.Getenv("DEBUG")

	if "" != env {
		Enable(env)
	}
}

// SetWriter replaces the default of os.Stderr with `w`.
func SetWriter(w io.Writer) {
	m.Lock()
	defer m.Unlock()
	writer = w
}

// Disable all pattern matching. This function is thread-safe.
func Disable() {
	m.Lock()
	defer m.Unlock()
	enabled = false
}

// Enable the given debug `pattern`. Patterns take a glob-like form,
// for example if you wanted to enable everything, just use "*", or
// if you had a library named mongodb you could use "mongodb:connection",
// or "mongodb:*". Multiple matches can be made with a comma, for
// example "mongo*,redis*".
//
// This function is thread-safe.
func Enable(pattern string) {
	m.Lock()
	defer m.Unlock()
	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.Replace(pattern, "\\*", ".*?", -1)
	pattern = strings.Replace(pattern, ",", "|", -1)
	pattern = "^(" + pattern + ")$"
	reg = regexp.MustCompile(pattern)
	enabled = true
}

// Debug creates a debug function for `name` which you call
// with printf-style arguments in your application or library.
func Debug(name string) Debugger {
	prevGlobal := time.Now()
	color := colors[rand.Intn(len(colors))]
	prev := time.Now()

	dbg := Debugger{}

	dbg.Spawn = func(ns string) Debugger {
		return Debug(name + ":" + ns)
	}

	dbg.Log = func(args ...interface{}) {
		var strOrFunc interface{}
		var format string
		var isString bool

		if !enabled {
			return
		}

		if !reg.MatchString(name) {
			return
		}

		if len(args) >= 1 {
			strOrFunc = args[0]

			format, isString = strOrFunc.(string)

			if !isString {
				lazy, isFunc := strOrFunc.(func() string)
				if !isFunc {
					panic("invalid first argument type for Debug, must either be a string or lazy function")
				}
				format = lazy()
			}
		}

		d := deltas(prevGlobal, prev, color)
		fmt.Fprintf(writer, d+" \033["+color+"m"+name+"\033[0m - "+format+"\n", args...)
		prevGlobal = time.Now()
		prev = time.Now()
	}

	return dbg
}

// Return formatting for deltas.
func deltas(prevGlobal, prev time.Time, color string) string {
	now := time.Now()
	global := now.Sub(prevGlobal).Nanoseconds()
	delta := now.Sub(prev).Nanoseconds()
	ts := now.UTC().Format("15:04:05.000")
	deltas := fmt.Sprintf("%s %-6s \033["+color+"m%-6s", ts, humanizeNano(global), humanizeNano(delta))
	return deltas
}

// Humanize nanoseconds to a string.
func humanizeNano(n int64) string {
	var suffix string

	switch {
	case n > 1e9:
		n /= 1e9
		suffix = "s"
	case n > 1e6:
		n /= 1e6
		suffix = "ms"
	case n > 1e3:
		n /= 1e3
		suffix = "us"
	default:
		suffix = "ns"
	}

	return strconv.Itoa(int(n)) + suffix
}
