package main

import (
	"time"

	"github.com/nmccready/go-debug/example/rootDebug"
)

func main() {
	var debug = rootDebug.Spawn("single")
	for {
		debug.Log("sending mail")
		debug.Log("send email to %s", "tobi@segment.io")
		debug.Log("send email to %s", "loki@segment.io")
		debug.Log("send email to %s", "jane@segment.io")
		time.Sleep(500 * time.Millisecond)
	}
}
