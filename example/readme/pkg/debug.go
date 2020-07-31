package readme

import dbg "github.com/nmccready/go-debug"

var rootDebug = dbg.Debug("app-name")
var Spawn = rootDebug.Spawn
var Log = rootDebug.Log
