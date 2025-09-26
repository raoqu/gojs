package api

import (
	"context"

	"github.com/dop251/goja"
	"github.com/raoqu/gojs/internal/util"
)

func toGojaResult(rt *goja.Runtime, val interface{}, err error) (goja.Value, error) {
	if err != nil {
		return goja.Undefined(), err
	}
	return rt.ToValue(val), nil
}

// Redis_set (group, key, value)
func Redis_set(rt *goja.Runtime, call goja.FunctionCall) (goja.Value, error) {
	var group string
	var key string
	var value string
	// var outputBuffer strings.Builder
	if len(call.Arguments) < 2 {
		return goja.Undefined(), nil
	} else if len(call.Arguments) == 2 {
		key = call.Arguments[0].String()
		value = call.Arguments[1].String()
		util.RedisData.Set(context.Background(), key, value, 0)
	} else {
		group = call.Arguments[0].String()
		key = call.Arguments[1].String()
		value = call.Arguments[2].String()
		util.RedisData.HSet(context.Background(), group, key, value)
	}

	return goja.Undefined(), nil
}

// Redis_get (group, key) -> string
func Redis_get(rt *goja.Runtime, call goja.FunctionCall) (goja.Value, error) {
	var group string
	var key string
	if len(call.Arguments) < 2 {
		return goja.Undefined(), nil
	} else if len(call.Arguments) == 2 {
		group = call.Arguments[0].String()
		key = call.Arguments[1].String()
		val, err := util.RedisData.HGet(context.Background(), group, key).Result()
		return toGojaResult(rt, val, err)
	} else {
		key = call.Arguments[0].String()
		val, err := util.RedisData.Get(context.Background(), key).Result()
		return toGojaResult(rt, val, err)
	}
}

func Redis_keys(rt *goja.Runtime, call goja.FunctionCall) (goja.Value, error) {
	if len(call.Arguments) < 1 {
		return goja.Undefined(), nil
	}
	group := call.Arguments[0].Export().(string)
	val, err := util.RedisData.HKeys(context.Background(), group).Result()
	return toGojaResult(rt, val, err)
}

func Redis_hgetall(rt *goja.Runtime, call goja.FunctionCall) (goja.Value, error) {
	if len(call.Arguments) < 1 {
		return goja.Undefined(), nil
	}
	group := call.Arguments[0].Export().(string)
	val, err := util.RedisData.HGetAll(context.Background(), group).Result()
	return toGojaResult(rt, val, err)
}
