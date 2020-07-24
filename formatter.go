package debug

import "fmt"

// highly inspired by logrus
type Formatter interface {
	Format(*Debugger, string) string
}

type TextFormatter struct {
	HasColor bool
	HasTime  bool
}

func (t *TextFormatter) Format(dbg *Debugger, msg string) string {
	color := dbg.color

	head := getTime(dbg.prevGlobal, dbg.prev, color, hasTime) +
		getColorStr(color, hasColors) + dbg.name + getColorOff(hasColors)

	if head != "" {
		head += " - "
	}

	return fmt.Sprintf("%s%s\n", head, msg)
}
