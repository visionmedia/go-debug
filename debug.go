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

	goCache "github.com/patrickmn/go-cache"
)

var (
	writer    io.Writer = os.Stderr
	reg       *regexp.Regexp
	m         sync.Mutex
	enabled   = false
	cache     *goCache.Cache
	hasColors           = true
	hasTime             = true
	formatter Formatter = &TextFormatter{}
)

// Debugger function.
type DebugFunction func(...interface{})

type Fields map[string]interface{}

type Debugger struct {
	name       string
	prevGlobal time.Time
	prev       time.Time
	fields     Fields
	WithFields func(map[string]interface{}) Debugger
	color      string

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
	cacheMinStr := os.Getenv("DEBUG_CACHE_MINUTES")
	colorOffStr := os.Getenv("DEBUG_COLOR_OFF")
	timeOffStr := os.Getenv("DEBUG_TIME_OFF")

	if "" != env {
		Enable(env)
	}

	SetHasColors(colorOffStr == "")
	SetHasTime(timeOffStr == "")

	err := SetCache(cacheMinStr)

	if err != nil {
		panic(err)
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

/*
	Enable the given debug `pattern`. Patterns take a glob-like form,
	for example if you wanted to enable everything, just use "*", or
	if you had a library named mongodb you could use "mongodb:connection",
	or "mongodb:*". Multiple matches can be made with a comma, for
	example "mongo*,redis*".

	This function is thread-safe.
*/
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

/*
	Initialize the namespace cache which defaults to 60 min lifespan
*/
func SetCache(cacheMinStr string) error {
	var err error
	cachedMin := 60

	m.Lock()
	defer m.Unlock()
	if cacheMinStr != "" {
		cachedMin, err = strconv.Atoi(cacheMinStr)
		if err != nil {
			return err
		}
	}

	cache = goCache.New(time.Duration(cachedMin)*time.Minute, 10*time.Minute)
	return nil
}

func setBoolWithLock(work func(bool)) func(bool) {
	return func(isOn bool) {
		m.Lock()
		defer m.Unlock()
		work(isOn)
	}
}

var SetHasColors = setBoolWithLock(func(isOn bool) {
	hasColors = isOn
})

var SetHasTime = setBoolWithLock(func(isOn bool) {
	hasTime = isOn
})

// Debug creates a debug function for `name` which you call
// with printf-style arguments in your application or library.
func Debug(name string) Debugger {
	entry, cached := cache.Get(name)

	if cached {
		dbg, _ := entry.(Debugger)
		return dbg
	}

	dbg := Debugger{}

	dbg.name = name
	dbg.prevGlobal = time.Now()
	dbg.color = colors[rand.Intn(len(colors))]
	dbg.prev = time.Now()

	dbg.Spawn = func(ns string) Debugger {
		return Debug(name + ":" + ns)
	}

	dbg.Log = func(args ...interface{}) {
		var strOrFunc interface{}
		var msg string
		var isString bool

		if !enabled {
			return
		}

		if !reg.MatchString(name) {
			return
		}

		if len(args) >= 1 {
			strOrFunc = args[0]
			args = args[1:]

			msg, isString = strOrFunc.(string)

			if !isString {
				lazy, isFunc := strOrFunc.(func() string)
				if !isFunc {
					// coerce to string
					msg = fmt.Sprint(strOrFunc)
				} else {
					msg = lazy()
				}
			}
		}

		dbg.WithFields = func(fields map[string]interface{}) Debugger {
			return dbg.Spawn("junk")
		}

		preppedMsg := formatter.Format(&dbg, msg)

		fmt.Fprintf(writer, preppedMsg, args...)
		dbg.prevGlobal = time.Now()
		dbg.prev = time.Now()
	}

	cache.Set(name, dbg, goCache.DefaultExpiration)

	return dbg
}

func getColorStr(color string, isOn bool) string {
	if !isOn {
		return ""
	}
	return "\033[" + color + "m"
}

func getColorOff(isOn bool) string {
	if !isOn {
		return ""
	}
	return "\033[0m"
}

func getTime(prevGlobal time.Time, prev time.Time, color string, isOn bool) string {
	if !isOn {
		return ""
	}
	d := deltas(prevGlobal, prev, color)

	return d
}

// Return formatting for deltas.
func deltas(prevGlobal, prev time.Time, color string) string {
	now := time.Now()
	global := now.Sub(prevGlobal).Nanoseconds()
	delta := now.Sub(prev).Nanoseconds()
	ts := now.UTC().Format("15:04:05.000")
	deltas := fmt.Sprintf("%s %-6s "+getColorStr(color, hasColors)+"%-6s", ts, humanizeNano(global), humanizeNano(delta))
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
