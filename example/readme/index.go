package main

import (
	"flag"
	"time"

	. "github.com/nmccready/go-debug"
	. "github.com/nmccready/go-debug/example/readme/pkg"
)

type CliOpts struct {
	json       bool
	pretty     bool
	fieldsOnly bool
}

func main() {
	opts := CliOpts{}

	flag.BoolVar(&opts.json, "json", false, "set true for JSONFormatter")
	flag.BoolVar(&opts.pretty, "pretty", false, "set true to make json pretty")
	flag.BoolVar(&opts.fieldsOnly, "fieldsOnly", false, "set true to make text formatter fields only")

	flag.Parse()

	if opts.json {
		SetFormatter(&JSONFormatter{PrettyPrint: opts.pretty})
	}

	if opts.fieldsOnly {
		SetFormatter(&TextFormatter{HasFieldsOnly: true})
	}

	var debug = Spawn("main")
	var sibling = Spawn("sibling")

	for {
		// app-name:main ....
		debug.WithField("key", "value").Log("sending mail")
		debug.Log(func() string { return "sending mail" })
		debug.Log("send email to %s", "tobi@segment.io")
		debug.Log("send email to %s", "loki@segment.io")
		debug.Log("send email to %s", "jane@segment.io")
		debug.Error("oh noes")
		debug.Spawn("child").Log("hi from child")
		sibling.Spawn("a").Log("hi")
		sibling.Spawn("b").Log("hi")
		sibling.Spawn("b").Error("sad")
		time.Sleep(500 * time.Millisecond)
	}
}
