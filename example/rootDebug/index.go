package rootDebug

import debug "github.com/nmccready/go-debug"

var _rootDebug = debug.Debug("example")

var Spawn = _rootDebug.Spawn
var Log = _rootDebug.Log
