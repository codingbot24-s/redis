package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET" : Set,
	"GET" : Get,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

var SETs = map[string]string{}
var SETsm = sync.RWMutex{}

func Set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "invalid number of arguments"}
	}
	key := args[0].bulk
	val := args[1].bulk

	SETsm.Lock()
	SETs[key] = val
	SETsm.Unlock()

	return Value{typ: "string", str: "OK"}
}

func Get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "invalid number of arguments"}
	}

	key := args[0].bulk
	SETsm.Lock()
	val, ok := SETs[key]
	SETsm.Unlock()
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}
