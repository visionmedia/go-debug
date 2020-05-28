package main

import (
	"time"

	. "github.com/nmccready/go-debug"
)

var debug = Debug("single")

func main() {
	for {
		debug("sending mail")
		debug("send email to %s", "tobi@segment.io")
		debug("send email to %s", "loki@segment.io")
		debug("send email to %s", "jane@segment.io")
		time.Sleep(500 * time.Millisecond)
	}
}
