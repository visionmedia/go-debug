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
	writer     io.Writer = os.Stderr
	reg        *regexp.Regexp
	neg        []*regexp.Regexp
	m          sync.Mutex
	enabled    = false
	cache      *goCache.Cache
	HAS_COLORS           = true
	HAS_TIME             = true
	formatter  Formatter = &TextFormatter{HasColor: true}
	negRegEx             = regexp.MustCompile(`^-.*?`)
)

type Fields map[string]interface{}

type Debugger struct {
	name   string
	prev   time.Time
	fields Fields
	color  string
}

type IDebugger interface {
	Log(...interface{})
	Error(...interface{})
	Spawn(ns string) *Debugger
	WithFields(fields map[string]interface{}) *Debugger
	WithField(key string, value interface{}) *Debugger
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

func SetFormatter(f Formatter) {
	m.Lock()
	defer m.Unlock()
	formatter = f
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
	pattern, neg = BuildPattern(pattern)
	reg = regexp.MustCompile(pattern)
	enabled = true
}

func BuildPattern(pattern string) (string, []*regexp.Regexp) {
	pattern = regexp.QuoteMeta(pattern)
	pattern = RegExWildCard(pattern)
	pattern = strings.Replace(pattern, ",", "|", -1)

	pattern, negatives := BuildNegativeMatches(pattern, `|`)
	return RegExWrap(pattern), negatives
}

func RegExWildCard(pattern string) string {
	return strings.Replace(pattern, "\\*", ".*?", -1)
}

func RegExWrap(pattern string) string {
	return "^(" + pattern + ")$"
}

func RegExWrapCompile(pattern string) *regexp.Regexp {
	return regexp.MustCompile(RegExWrap(pattern))
}

/*
	Find all negative namespaces and pull them out of the pattern.

	But build a slice of negatives

	example: pattern="*,-somenamespace"

	Returns: "*", ["somenamespace"]
*/
func BuildNegativeMatches(pattern string, orString string) (string, []*regexp.Regexp) {
	if orString == "" {
		orString = ","
	}
	var negatives []*regexp.Regexp
	maybeNegs := strings.Split(pattern, orString)

	// fmt.Printf("maybeNegs: %+v\n", maybeNegs)

	for _, s := range maybeNegs {
		if negRegEx.MatchString(s) {
			negatives = append(negatives, RegExWrapCompile(strings.Replace(s, "-", "", 1)))
			pattern = strings.Replace(pattern, orString+s, "", -1)
		}
	}
	return pattern, negatives
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
	HAS_COLORS = isOn
})

var SetHasTime = setBoolWithLock(func(isOn bool) {
	HAS_TIME = isOn
})

// Debug creates a debug function for `name` which you call
// with printf-style arguments in your application or library.
func Debug(name string) *Debugger {
	entry, cached := cache.Get(name)

	if cached {
		dbg, _ := entry.(Debugger)
		return &dbg
	}

	dbg := Debugger{name: name, prev: time.Now(), color: colors[rand.Intn(len(colors))]}

	if formatter.GetHasFieldsOnly() {
		dbg.WithFields(map[string]interface{}{"namespace": name, "msg": nil})

		if HAS_TIME {
			dbg.WithFields(map[string]interface{}{"time": nil, "delta": nil})
		}
	}

	cache.Set(name, dbg, goCache.DefaultExpiration)

	return &dbg
}

func (dbg *Debugger) Spawn(ns string) *Debugger {
	return Debug(dbg.name + ":" + ns)
}

func (dbg *Debugger) Log(args ...interface{}) {
	if !enabled {
		return
	}

	if !reg.MatchString(dbg.name) {
		return
	}

	for _, n := range neg {
		if n.MatchString(dbg.name) {
			return
		}
	}

	var strOrFunc interface{}
	var msg string
	var isString bool

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

	preppedMsg := formatter.Format(dbg, msg)

	fmt.Fprintf(writer, preppedMsg, args...)
	dbg.prev = time.Now()
}

func (dbg *Debugger) Error(args ...interface{}) {
	// prepend error name as it is easier to filter!
	Debug("error:" + dbg.name).Log(args...)
}

func (dbg *Debugger) WithFields(fields map[string]interface{}) *Debugger {
	if len(dbg.fields) == 0 {
		dbg.fields = fields
		return dbg
	}

	for k, v := range fields {
		dbg.fields[k] = v
	}
	return dbg
}

func (dbg *Debugger) WithField(key string, value interface{}) *Debugger {
	if len(dbg.fields) == 0 {
		dbg.fields = map[string]interface{}{}
	}

	dbg.fields[key] = value
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

func getTime(timestring string, delta string, isOn bool) string {
	if !isOn {
		return ""
	}

	return fmt.Sprintf("%s %-6s", timestring, delta)
}

// Return formatting for deltas.
func deltas(prev time.Time) (string, string) {
	now := time.Now()
	delta := now.Sub(prev).Nanoseconds()
	ts := now.UTC().Format("15:04:05.000")
	return ts, humanizeNano(delta)
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
