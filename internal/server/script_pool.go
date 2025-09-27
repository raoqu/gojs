package server

import (
	"fmt"

	"github.com/raoqu/gojs/internal/api"
	"github.com/raoqu/gojs/internal/util"
	"github.com/raoqu/gojs/internal/util/mysql"
)

// Global script pool for caching compiled scripts
var scriptPool *api.ScriptPool

// initScriptPool initializes the 	global script pool if not already initialized
func InitScriptPool(poolName string) *api.ScriptPool {
	var store api.ScriptStore = api.NewScriptLocalStore()
	if util.RedisConfig != nil {
		store = api.NewScriptRedisStore(poolName, util.RedisConfig)
	}
	scriptPool = api.NewScriptPool(poolName, store)

	// Inject console functions
	scriptPool.Inject("console.log", api.Console_log)
	scriptPool.Inject("console.error", api.Console_error)

	// Inject Redis functions
	scriptPool.Inject("redis.set", api.Redis_set)
	scriptPool.Inject("redis.get", api.Redis_get)
	scriptPool.Inject("redis.keys", api.Redis_keys)
	scriptPool.Inject("redis.hgetall", api.Redis_hgetall)

	// Inject Redis set operations
	scriptPool.Inject("redis.sadd", api.Redis_sadd)
	scriptPool.Inject("redis.srem", api.Redis_srem)
	scriptPool.Inject("redis.scard", api.Redis_scard)
	scriptPool.Inject("redis.smembers", api.Redis_smembers)

	// Inject MySQL functions
	if mysql.MYSQL_CLIENT != nil {
		scriptPool.Inject("mysql.query", api.MySQL_query)
		scriptPool.Inject("mysql.exec", api.MySQL_exec)
		scriptPool.Inject("mysql.queryRow", api.MySQL_queryRow)
		scriptPool.Inject("mysql.transaction", api.MySQL_transaction)
	}

	// Inject Net functions
	scriptPool.Inject("net.fetch", api.Net_fetch)

	// Inject Sys functions
	scriptPool.Inject("sys.command", api.Sys_command)

	return scriptPool
}

var EnableScript = true

// executeJavaScript runs a JavaScript code using goja with ScriptPool for caching
func executeJavaScript(name string, ctx map[string]interface{}) (interface{}, error) {
	if !EnableScript {
		return "", fmt.Errorf("script disabled")
	}

	code, err := scriptPool.Cache.GetScript(name)
	if err != nil {
		return "", fmt.Errorf("failed to get script: %s %v", name, err)
	}

	// Set the script in the pool (compiles and caches it)
	err = scriptPool.SetScript(name, code)
	if err != nil {
		return "", fmt.Errorf("failed to compile script: %v", err)
	}

	// Run the script from the pool
	result, err := scriptPool.RunScript(name, ctx)
	if err != nil {
		return "", fmt.Errorf("failed to run script: %v", err)
	}

	return result.Value, nil
}
