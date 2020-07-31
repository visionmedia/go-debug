# go-debug [![build status][travis-image]][travis-url]

Conditional debug logging for Go libraries.

View the [docs](http://godoc.org/github.com/tj/go-debug).

## Installation

```
$ go get github.com/nmccready/go-debug
```

## [Example](./example/readme/index.go)

### Show me everthing

```bash
$ DEBUG=* go run ./example/readme/index.go
18:18:25.215 4us   app-name:main - sending mail
    key=value
18:18:25.215 986ns app-name:main - sending mail
    key=value
18:18:25.215 516ns app-name:main - send email to tobi@segment.io
    key=value
18:18:25.215 550ns app-name:main - send email to loki@segment.io
    key=value
18:18:25.215 486ns app-name:main - send email to jane@segment.io
    key=value
18:18:25.215 1us   error:app-name:main - oh noes
18:18:25.215 852ns app-name:main:child - hi from child
18:18:25.215 804ns app-name:sibling:a - hi
18:18:25.215 804ns app-name:sibling:b - hi
18:18:25.215 850ns error:app-name:sibling:b - sad
```

### App Domain only

```bash
$ DEBUG=app-name* go run ./example/readme/index.go
18:19:00.567 5us   app-name:main - sending mail
    key=value
18:19:00.567 916ns app-name:main - sending mail
    key=value
18:19:00.567 493ns app-name:main - send email to tobi@segment.io
    key=value
18:19:00.567 458ns app-name:main - send email to loki@segment.io
    key=value
18:19:00.567 423ns app-name:main - send email to jane@segment.io
    key=value
18:19:00.567 736ns app-name:main:child - hi from child
18:19:00.567 793ns app-name:sibling:a - hi
18:19:00.567 713ns app-name:sibling:b - hi
```

### Filter out errors

```bash
$ DEBUG=*,-error* go run ./example/readme/index.go
18:19:20.595 5us   app-name:main - sending mail
    key=value
18:19:20.595 1us   app-name:main - sending mail
    key=value
18:19:20.595 611ns app-name:main - send email to tobi@segment.io
    key=value
18:19:20.595 584ns app-name:main - send email to loki@segment.io
    key=value
18:19:20.595 505ns app-name:main - send email to jane@segment.io
    key=value
18:19:20.595 801ns app-name:main:child - hi from child
18:19:20.595 779ns app-name:sibling:a - hi
18:19:20.595 794ns app-name:sibling:b - hi
```

### Errors Only

```bash
$ DEBUG=error* go run ./example/readme/index.go
17:35:13.401 1us   error:app-name:main - oh noes
17:35:13.401 1us   error:app-name:sibling:b - sad
```

### Everything but omit some

Omitting sibling but report their errors

```bash
$ DEBUG=*,-app-name:sibling* go run ./example/readme/index.go
18:19:48.034 5us   app-name:main - sending mail
    key=value
18:19:48.034 1us   app-name:main - sending mail
    key=value
18:19:48.034 730ns app-name:main - send email to tobi@segment.io
    key=value
18:19:48.034 648ns app-name:main - send email to loki@segment.io
    key=value
18:19:48.034 533ns app-name:main - send email to jane@segment.io
    key=value
18:19:48.034 1us   error:app-name:main - oh noes
18:19:48.034 1us   app-name:main:child - hi from child
18:19:48.034 1us   error:app-name:sibling:b - sad
```

## Formatters

By Default the TextFormatter is w/o Fields as seen above. Fields can be added via `WithField` or `WithFields`. If `TextFormatter{HasFieldsOnly:true}` then all
fields namespace, time, msg, and delta will be inlined as fields.

Otherwise the message is printed as above but with fields below it.

### TextFormatter HasFieldsOnly

```bash
$ DEBUG=* go run ./example/readme/index.go -fieldsOnly
delta=6us key=value msg="sending mail" namespace=app-name:main time=18:14:49.462
delta=925ns key=value msg="sending mail" namespace=app-name:main time=18:14:49.462
delta=532ns key=value msg="send email to tobi@segment.io" namespace=app-name:main time=18:14:49.462
delta=528ns key=value msg="send email to loki@segment.io" namespace=app-name:main time=18:14:49.462
delta=466ns key=value msg="send email to jane@segment.io" namespace=app-name:main time=18:14:49.462
delta=1us msg="oh noes" namespace=error:app-name:main time=18:14:49.462
delta=1us msg="hi from child" namespace=app-name:main:child time=18:14:49.462
delta=3us msg=hi namespace=app-name:sibling:a time=18:14:49.462
delta=1us msg=hi namespace=app-name:sibling:b time=18:14:49.462
delta=7us msg=sad namespace=error:app-name:sibling:b time=18:14:49.462
```

### JSONFormatter

```json
$ DEBUG=* go run ./example/readme/index.go -json
{"delta":"6us","key":"value","msg":"sending mail","namespace":"app-name:main","time":"18:15:17.113"}
{"delta":"892ns","key":"value","msg":"sending mail","namespace":"app-name:main","time":"18:15:17.113"}
{"delta":"640ns","key":"value","msg":"send email to tobi@segment.io","namespace":"app-name:main","time":"18:15:17.113"}
{"delta":"604ns","key":"value","msg":"send email to loki@segment.io","namespace":"app-name:main","time":"18:15:17.113"}
{"delta":"505ns","key":"value","msg":"send email to jane@segment.io","namespace":"app-name:main","time":"18:15:17.113"}
{"delta":"1us","msg":"oh noes","namespace":"error:app-name:main","time":"18:15:17.113"}
{"delta":"1us","msg":"hi from child","namespace":"app-name:main:child","time":"18:15:17.113"}
{"delta":"1us","msg":"hi","namespace":"app-name:sibling:a","time":"18:15:17.113"}
{"delta":"1us","msg":"hi","namespace":"app-name:sibling:b","time":"18:15:17.113"}
{"delta":"4us","msg":"sad","namespace":"error:app-name:sibling:b","time":"18:15:17.113"}
```

#### W/ PrettyPrint

```json
$ DEBUG=* go run ./example/readme/index.go -json -pretty
{
  "delta": "6us",
  "key": "value",
  "msg": "sending mail",
  "namespace": "app-name:main",
  "time": "18:16:13.351"
}
{
  "delta": "912ns",
  "key": "value",
  "msg": "sending mail",
  "namespace": "app-name:main",
  "time": "18:16:13.351"
}
{
  "delta": "561ns",
  "key": "value",
  "msg": "send email to tobi@segment.io",
  "namespace": "app-name:main",
  "time": "18:16:13.351"
}
{
  "delta": "580ns",
  "key": "value",
  "msg": "send email to loki@segment.io",
  "namespace": "app-name:main",
  "time": "18:16:13.351"
}
{
  "delta": "518ns",
  "key": "value",
  "msg": "send email to jane@segment.io",
  "namespace": "app-name:main",
  "time": "18:16:13.351"
}
{
  "delta": "1us",
  "msg": "oh noes",
  "namespace": "error:app-name:main",
  "time": "18:16:13.351"
}
{
  "delta": "1us",
  "msg": "hi from child",
  "namespace": "app-name:main:child",
  "time": "18:16:13.351"
}
{
  "delta": "1us",
  "msg": "hi",
  "namespace": "app-name:sibling:a",
  "time": "18:16:13.351"
}
{
  "delta": "1us",
  "msg": "hi",
  "namespace": "app-name:sibling:b",
  "time": "18:16:13.351"
}
{
  "delta": "4us",
  "msg": "sad",
  "namespace": "error:app-name:sibling:b",
  "time": "18:16:13.351"
}
```

## Timestamp & Delta

A timestamp and one delta . The timestamp consists of hour, minute, second and microseconds. The delta is relative to the previous debug call of any name.

## The environment variables

### DEBUG

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

### DEBUG_CACHE_MINUTES

Integer in minutes to set the debug caching. By default each debug namespace is cached 60 minutes via Spawn.

### DEBUG_COLOR_OFF

Default is false and color is on. Truthy ie anything not an empty string is true. Thus, turning color off.

### DEBUG_TIME_OFF

Default is false and timestamp & delta are on. Truthy ie anything not an empty string is true. Thus, turning times off.

## License

MIT

## Development

Install:

- yq - `brew install yq`
- golines - `go get -u github.com/segmentio/golines`
- golangci-lint - `brew install golangci/tap/golangci-lint`

[travis-image]: https://img.shields.io/travis/nmccready/go-debug.svg
[travis-url]: https://travis-ci.org/nmccready/go-debug
