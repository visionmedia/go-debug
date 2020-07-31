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
	FieldKeyError     = "debug_error"
)

// highly inspired by logrus
type Formatter interface {
	Format(*Debugger, string) string
	GetHasFieldsOnly() bool
}

type TextFormatter struct {
	HasColor         bool
	ForceQuote       bool
	QuoteEmptyFields bool
	DisableQuote     bool
	HasFieldsOnly    bool
	SortingFunc      func(keys []string)
}

func (t *TextFormatter) Format(dbg *Debugger, msg string) string {
	mainMsg := ""
	fields := ""

	var keys []string

	finalized := finalizeFields(dbg, msg, HAS_COLORS && t.HasColor, func(k string, v interface{}) interface{} {
		keys = append(keys, k)
		return nil
	})

	if t.SortingFunc == nil {
		sort.Strings(keys)
	} else {
		t.SortingFunc(keys)
	}

	// build fields string in specified order
	for _, k := range keys {
		v := finalized.Fields[k]
		fields = t.appendKeyValue(fields, k, v)
	}

	if !t.GetHasFieldsOnly() {
		mainMsg = BasicFormat(finalized.TimeString, finalized.Delta, finalized.Namespace, msg)
		if fields != "" {
			fields = "    " + fields
		}
	}

	if fields != "" {
		fields += "\n"
	}

	return mainMsg + fields
}
func (t *TextFormatter) GetHasFieldsOnly() bool {
	return t.HasFieldsOnly
}

func BasicFormat(ts string, delta string, ns string, msg string) string {
	time := getTime(ts, delta, HAS_TIME)
	head := time + ns

	if head != "" {
		head += " - "
	}

	return head + msg + "\n"
}

func (f *TextFormatter) appendKeyValue(s string, key string, value interface{}) string {
	if len(s) > 0 {
		s += " "
	}
	return s + key + "=" + f.appendValue(key, value)
}

func (f *TextFormatter) appendValue(key string, value interface{}) string {
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

type Finalized struct {
	Fields     Fields
	Namespace  string
	TimeString string
	Delta      string
}

func finalizeFields(dbg *Debugger, msg string, hasColor bool, cb func(string, interface{}) interface{}) *Finalized {
	ts, delta := deltas(dbg.prev)
	ns := getColorStr(dbg.color, hasColor) + dbg.name + getColorOff(hasColor)

	fields := Fields{}

	for k, v := range dbg.fields {
		switch {
		case k == FieldKeyNamespace:
			fields[k] = ns
		case k == FieldKeyMsg:
			fields[k] = msg
		case k == FieldKeyTime:
			fields[k] = ts
		case k == FieldKeyDelta:
			fields[k] = delta
		default:
			fields[k] = v
		}
		if cb != nil {
			transformed := cb(k, fields[k])
			if transformed != nil {
				fields[k] = transformed
			}
		}
	}
	return &Finalized{Fields: fields, Namespace: ns, TimeString: ts, Delta: delta}
}
