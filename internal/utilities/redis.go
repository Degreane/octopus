package utilities

import (
	"context"
	"fmt"
	"time"

	"github.com/degreane/octopus/internal/database"
	"github.com/degreane/octopus/internal/utilities/debug"
	lua "github.com/yuin/gopher-lua"
)

// GetRedisValueLua retrieves a value from Redis by key
// Usage in Lua: local value = eocto.getRedis("mykey")
func GetRedisValueLua(L *lua.LState) int {
	key := L.CheckString(1)

	// Get Redis client from the database package
	redisClient := database.GetRedisClient()
	if redisClient == nil {
		debug.Debug(debug.Error, fmt.Sprintf("Redis client not initialized"))
		L.Push(lua.LNil)
		return 1
	}

	ctx := context.Background()
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		// Key doesn't exist or other error
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(val))
	return 1
}

// SetRedisValueLua sets a value in Redis with an optional TTL
// Usage in Lua: eocto.setRedis("mykey", "myvalue", 3600) -- with TTL in seconds
// Usage in Lua: eocto.setRedis("mykey", "myvalue") -- without TTL (persistent)
func SetRedisValueLua(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckString(2)

	// Optional TTL parameter (in seconds)
	var ttl time.Duration
	if L.GetTop() >= 3 {
		ttlSeconds := L.CheckNumber(3)
		ttl = time.Duration(ttlSeconds) * time.Second
	}

	// Get Redis client from the database package
	redisClient := database.GetRedisClient()
	if redisClient == nil {
		debug.Debug(debug.Error, fmt.Sprintf("Redis client not initialized"))
		L.Push(lua.LBool(false))
		return 1
	}

	ctx := context.Background()
	var err error

	if ttl > 0 {
		// Set with expiration
		err = redisClient.Set(ctx, key, value, ttl).Err()
	} else {
		// Set without expiration
		err = redisClient.Set(ctx, key, value, 0).Err()
	}

	if err != nil {
		debug.Debug(debug.Error, fmt.Sprintf("Error setting Redis key %s: %v", key, err))
		L.Push(lua.LBool(false))
		return 1
	}

	L.Push(lua.LBool(true))
	return 1
}

// DeleteRedisKeyLua deletes a key from Redis
// Usage in Lua: local success = eocto.delRedis("mykey")
func DeleteRedisKeyLua(L *lua.LState) int {
	key := L.CheckString(1)

	// Get Redis client from the database package
	redisClient := database.GetRedisClient()
	if redisClient == nil {
		debug.Debug(debug.Error, fmt.Sprintf("Redis client not initialized"))
		L.Push(lua.LBool(false))
		return 1
	}

	ctx := context.Background()
	_, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		debug.Debug(debug.Error, fmt.Sprintf("Error deleting Redis key %s: %v", key, err))
		L.Push(lua.LBool(false))
		return 1
	}

	L.Push(lua.LBool(true))
	return 1
}
