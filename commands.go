package main

import (
	"sync"

	"github.com/Luisgustavom1/build-redis-from-scratch/resp"
)

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: resp.String, Str: "PONG"}
	}

	return resp.Value{Typ: resp.String, Str: args[0].Bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: resp.Error, Str: "Expected 2 arguments"}
	}

	k := args[0].Bulk
	v := args[1].Bulk

	SETsMu.Lock()
	SETs[k] = v
	SETsMu.Unlock()

	return resp.Value{Typ: resp.String, Str: "OK"}
}
