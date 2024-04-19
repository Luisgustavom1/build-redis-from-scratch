package main

import (
	"sync"

	"github.com/Luisgustavom1/build-redis-from-scratch/resp"
)

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
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

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: resp.Error, Str: "Expected 1 arguments"}
	}

	k := args[0].Bulk

	SETsMu.RLock()
	v, ok := SETs[k]
	SETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: resp.Null}
	}

	return resp.Value{Typ: resp.Bulk, Bulk: v}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: resp.Error, Str: "Expected 3 arguments"}
	}

	h := args[0].Bulk
	if _, ok := HSETs[h]; !ok {
		HSETs[h] = map[string]string{}
	}
	k := args[1].Bulk
	v := args[2].Bulk

	HSETsMu.Lock()
	HSETs[h][k] = v
	HSETsMu.Unlock()

	return resp.Value{Typ: resp.String, Str: "OK"}
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: resp.Error, Str: "Expected 2 arguments"}
	}

	h := args[0].Bulk
	k := args[1].Bulk

	HSETsMu.RLock()
	v, ok := HSETs[h][k]
	HSETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: resp.Null}
	}

	return resp.Value{Typ: resp.Bulk, Bulk: v}
}

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: resp.Error, Str: "Expected 1 arguments"}
	}

	h := args[0].Bulk

	HSETsMu.RLock()
	maps, ok := HSETs[h]
	HSETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: resp.Null}
	}

	values := make([]resp.Value, len(maps)*2)
	i := 0
	for k, v := range maps {
		values[i] = resp.Value{Typ: resp.Bulk, Bulk: k}
		i++
		values[i] = resp.Value{Typ: resp.Bulk, Bulk: v}
		i++
	}

	return resp.Value{Typ: resp.Array, Array: values}
}
