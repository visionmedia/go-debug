package main

import (
	"time"

	. "github.com/nmccready/go-debug/example/readme/pkg"
)

var debug = Spawn("main")
var sibling = Spawn("sibling")

func main() {
	for {
		// app-name:main ....
		debug.Log("sending mail")
		debug.Log(func() string { return "sending mail" })
		debug.Log("send email to %s", "tobi@segment.io")
		debug.Log("send email to %s", "loki@segment.io")
		debug.Log("send email to %s", "jane@segment.io")
		debug.Spawn("hi from child")
		sibling.Log("hi")
		helper()
		time.Sleep(500 * time.Millisecond)
	}
}

func helper() {
	debug.Spawn("helper").Log("hi") // -> app-name:main:helper hi
}
