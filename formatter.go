package debug

import (
	"fmt"
	"sort"
)

const (
	FieldKeyMsg       = "msg"
	FieldKeyNamespace = "namespace"
	FieldKeyTime      = "time"
	FieldKeyDelta     = "delta"
)

// highly inspired by logrus
type Formatter interface {
	Format(*Debugger, string) string
	GetHasFieldsOnly() bool
}

type TextFormatter struct {
	HasColor         bool
	HasTime          bool
	ForceQuote       bool
	QuoteEmptyFields bool
	DisableQuote     bool
	HasFieldsOnly    bool
	SortingFunc      func(keys []string)
}

func (t *TextFormatter) Format(dbg *Debugger, msg string) string {
	mainMsg := ""
	color := dbg.color

	timestring, delta := deltas(dbg.prev)
	ns := getColorStr(color, hasColors) + dbg.name + getColorOff(hasColors)

	if !t.HasFieldsOnly {
		mainMsg = basicFormat(timestring, delta, ns, msg)
	}

	fields := ""

	var keys []string

	for k := range dbg.fields {
		keys = append(keys, k)
	}

	if t.SortingFunc == nil {
		sort.Strings(keys)
	} else {
		t.SortingFunc(keys)
	}

	for _, k := range keys {
		v := dbg.fields[k]
		switch {
		case k == FieldKeyNamespace:
			v = ns
		case k == FieldKeyMsg:
			v = msg
		case k == FieldKeyTime:
			v = timestring
		case k == FieldKeyDelta:
			v = delta
		default:
		}
		fields = t.appendKeyValue(fields, k, v)
	}

	if fields != "" {
		fields += "\n"
	}

	return fmt.Sprintf("%s%s", mainMsg, fields)
}
func (t *TextFormatter) GetHasFieldsOnly() bool {
	return t.HasFieldsOnly
}

func basicFormat(ts string, delta string, ns string, msg string) string {
	time := getTime(ts, delta, hasTime)
	head := fmt.Sprintf("%s%s", time, ns)

	if head != "" {
		head += " - "
	}

	return fmt.Sprintf("%s%s\n", head, msg)
}

func (f *TextFormatter) appendKeyValue(s string, key string, value interface{}) string {
	if len(s) > 0 {
		s += " "
	}
	return fmt.Sprintf("%s%s=%s", s, key, f.appendValue(s, key, value))
}

func (f *TextFormatter) appendValue(s string, key string, value interface{}) string {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if key == FieldKeyDelta || key == FieldKeyTime || key == FieldKeyNamespace || !f.needsQuoting(stringVal) {
		return stringVal
	}

	return fmt.Sprintf("%q", stringVal)
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.ForceQuote {
		return true
	}
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if f.DisableQuote {
		return false
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}
