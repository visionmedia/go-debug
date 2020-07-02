# go-debug

Conditional debug logging for Go libraries.

View the [docs](http://godoc.org/github.com/tj/go-debug).

## Installation

```
$ go get github.com/tj/go-debug
```

## Example

### debug.go (your app)

```go
package debug

import debug "github.com/nmccready/go-debug"

const Debug = Debug("app-name")
```

### main.go

```go
package main

import dbg "github.com/YOUR-APP/WHEREVER_DEBUG_GO_IS"
import "time"

var debug = dbg.Spawn("main")

func main() {
  for {
    // app-name:main ....
    debug.Log("sending mail")
    debug.Log(func()string { return "sending mail" })
    debug.Log("send email to %s", "tobi@segment.io")
    debug.Log("send email to %s", "loki@segment.io")
    debug.Log("send email to %s", "jane@segment.io")
    helper()
    time.Sleep(500 * time.Millisecond)
  }
}

func helper() {
  debug.Spawn("helper").Log("hi") // -> app-name:main:helper hi
}
```

If you run the program with the `DEBUG=*` environment variable you will see:

```
15:58:15.115 34us   33us   app-name:main - sending mail
15:58:15.115 34us   33us   app-name:main - sending mail
15:58:15.116 3us    3us    app-name:main - send email to tobi@segment.io
15:58:15.116 1us    1us    app-name:main - send email to loki@segment.io
15:58:15.116 1us    1us    app-name:main - send email to jane@segment.io
15:58:15.117 1us    1us    app-name:main:helper - hi
15:58:15.620 504ms  504ms  app-name:main - sending mail
15:58:15.620 504ms  504ms  app-name:main - sending mail
15:58:15.620 6us    6us    app-name:main - send email to tobi@segment.io
15:58:15.620 4us    4us    app-name:main - send email to loki@segment.io
15:58:15.620 4us    4us    app-name:main - send email to jane@segment.io
15:58:15.621 5us    5us    app-name:main:helper - hi
15:58:16.123 503ms  503ms  app-name:main - sending mail
15:58:16.123 504ms  504ms  app-name:main - sending mail
15:58:16.123 7us    7us    app-name:main - send email to tobi@segment.io
15:58:16.123 4us    4us    app-name:main - send email to loki@segment.io
15:58:16.123 4us    4us    app-name:main - send email to jane@segment.io
15:58:16.124 5us    5us    app-name:main:helper - hi
15:58:16.625 501ms  501ms  app-name:main - sending mail
15:58:15.620 504ms  504ms  app-name:main - sending mail
15:58:16.625 4us    4us    app-name:main - send email to tobi@segment.io
15:58:16.625 4us    4us    app-name:main - send email to loki@segment.io
15:58:16.625 5us    5us    app-name:main - send email to jane@segment.io
15:58:16.626 5us    5us    app-name:main:helper - hi
```

A timestamp and two deltas are displayed. The timestamp consists of hour, minute, second and microseconds. The left-most delta is relative to the previous debug call of any name, followed by a delta specific to that debug function. These may be useful to identify timing issues and potential bottlenecks.

## The DEBUG environment variable

Executables often support `--verbose` flags for conditional logging, however
libraries typically either require altering your code to enable logging,
or simply omit logging all together. go-debug allows conditional logging
to be enabled via the **DEBUG** environment variable, where one or more
patterns may be specified.

For example suppose your application has several models and you want
to output logs for users only, you might use `DEBUG=models:user`. In contrast
if you wanted to see what all database activity was you might use `DEBUG=models:*`,
or if you're love being swamped with logs: `DEBUG=*`. You may also specify a list of names delimited by a comma, for example `DEBUG=mongo,redis:*`.

The name given _should_ be the package name, however you can use whatever you like.

## License

MIT

## Development

Install:

- yq - `brew install yq`
- golines - `go get -u github.com/segmentio/golines`
- golangci-lint - `brew install golangci/tap/golangci-lint`
