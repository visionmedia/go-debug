# go-debug [![build status][travis-image]][travis-url]

Conditional debug logging for Go libraries.

View the [docs](http://godoc.org/github.com/tj/go-debug).

## Installation

```
$ go get github.com/tj/go-debug
```

## [Example](./example/readme/index.go)

[`DEBUG=* go run ./example/readme/index.go`](./example/readme/index.go)

```
21:16:16.870 4us   app-name:main - sending mail
21:16:16.870 765ns app-name:main - sending mail
21:16:16.870 542ns app-name:main - send email to tobi@segment.io
21:16:16.870 455ns app-name:main - send email to loki@segment.io
21:16:16.870 485ns app-name:main - send email to jane@segment.io
21:16:16.870 50us  app-name:sibling - hi
21:16:16.870 846ns app-name:main:helper - hi
21:16:17.371 501ms app-name:main - sending mail
21:16:17.371 1us   app-name:main - sending mail
21:16:17.371 694ns app-name:main - send email to tobi@segment.io
21:16:17.371 667ns app-name:main - send email to loki@segment.io
21:16:17.371 573ns app-name:main - send email to jane@segment.io
21:16:17.371 501ms app-name:sibling - hi
21:16:17.371 501ms app-name:main:helper - hi
21:16:17.871 500ms app-name:main - sending mail
21:16:17.871 1us   app-name:main - sending mail
21:16:17.871 761ns app-name:main - send email to tobi@segment.io
21:16:17.871 583ns app-name:main - send email to loki@segment.io
21:16:17.871 592ns app-name:main - send email to jane@segment.io
21:16:17.871 500ms app-name:sibling - hi
21:16:17.871 1s    app-name:main:helper - hi
21:16:18.371 500ms app-name:main - sending mail
21:16:18.371 865ns app-name:main - sending mail
21:16:18.371 601ns app-name:main - send email to tobi@segment.io
21:16:18.371 558ns app-name:main - send email to loki@segment.io
21:16:18.371 547ns app-name:main - send email to jane@segment.io
21:16:18.371 500ms app-name:sibling - hi
21:16:18.371 1s    app-name:main:helper - hi
21:16:18.874 502ms app-name:main - sending mail
21:16:18.874 26us  app-name:main - sending mail
21:16:18.874 896ns app-name:main - send email to tobi@segment.io
21:16:18.874 716ns app-name:main - send email to loki@segment.io
21:16:18.874 652ns app-name:main - send email to jane@segment.io
21:16:18.874 503ms app-name:sibling - hi
21:16:18.874 2s    app-name:main:helper - hi
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
If your swamped you can also start omitting namespaces ie: `DEBUG=*,-mongo,-redis`.

The name given _should_ be the package name, however you can use whatever you like.

## License

MIT

## Development

Install:

- yq - `brew install yq`
- golines - `go get -u github.com/segmentio/golines`
- golangci-lint - `brew install golangci/tap/golangci-lint`

[travis-image]: https://img.shields.io/travis/nmccready/go-debug.svg
[travis-url]: https://travis-ci.org/nmccready/go-debug
