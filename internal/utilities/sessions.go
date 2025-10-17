// Package utilities provides helper functions and utilities for the application.
//
// This package specifically handles session management integration with Lua scripts,
// allowing seamless session access and modification from Lua code.
package utilities

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/degreane/octopus/internal/middleware"
	"github.com/degreane/octopus/internal/utilities/debug"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	gob.Register(map[string]interface{}{})
}

// GetSession returns a Lua function that retrieves session values.
// The function is exposed to Lua scripts and provides access to the HTTP session.
//
// Parameters passed from Lua:
//   - key: string - The session key to retrieve
//
// Returns to Lua:
//   - string value if key exists
//   - nil if key doesn't exist or on error
//
// Usage in Lua:
//
//	local value = getSession("user_id")
func GetSession(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		sess, err := middleware.Store.Get(c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			L.Push(lua.LNil)
			return 1
		}
		// keys := sess.Keys()
		// log.Printf("Session keys: Get %v", keys)
		// for _, key := range keys {
		// 	log.Printf("Key: %s, Value: %v", key, sess.Get(key))
		// }
		key := L.ToString(1)
		raw := sess.Get(fmt.Sprintf("%s", key))
		if raw == nil {
			debug.Debug(debug.Info, fmt.Sprintf("Key \"%s\" not found in session locals", key))
			L.Push(lua.LNil)
			return 1
		}

		var jsonStr string
		switch v := raw.(type) {
		case string:
			jsonStr = v
		case []byte:
			jsonStr = string(v)
		default:
			// Fallback: if value isn't a JSON string, treat it as string
			jsonStr = fmt.Sprintf("%v", v)
		}

		var decoded interface{}
		if err := json.Unmarshal([]byte(jsonStr), &decoded); err != nil {
			// If it's not valid JSON, return as plain string for robustness
			debug.Debug(debug.Info, fmt.Sprintf("Session value for key '%s' is not JSON or unmarshal failed: %v", key, err))
			L.Push(lua.LString(jsonStr))
			return 1
		}

		// Convert decoded Go value into a corresponding lua.LValue
		L.Push(convertGoToLua(L, decoded))
		return 1
	}
}

// SetSession returns a Lua function that stores values in the session.
// The function is exposed to Lua scripts and allows modification of HTTP session data.
//
// Parameters passed from Lua:
//   - key: string - The session key to set
//   - value: string - The value to store
//
// Returns to Lua:
//   - nothing (void function)
//
// Usage in Lua:
//
//	setSession("user_id", "123")
//
// The function handles:
//   - Session retrieval and validation
//   - Value storage
//   - Session persistence
//   - Error logging
func SetSession(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		sess, err := middleware.Store.Get(c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			return 0
		}

		key := L.ToString(1)
		var value interface{}

		switch L.Get(2).Type() {
		case lua.LTString:
			value = L.ToString(2)
		case lua.LTNumber:
			value = float64(L.ToNumber(2))
		case lua.LTNil:
			value = nil
		case lua.LTBool:
			value = L.ToBool(2)
		case lua.LTTable:
			value = convertLuaTableToGo(L.ToTable(2))
		}
		// log.Printf("Setting session key %s to value %+v", key, value)
		jsonBytes, mErr := json.Marshal(value)
		if mErr != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error marshaling value to JSON: %v", mErr))
		} else {
			sess.Set(fmt.Sprintf("%s", key), string(jsonBytes))
			sess.Fresh()
			err = sess.Save()
			if err != nil {
				debug.Debug(debug.Error, fmt.Sprintf("Error saving session: %v", err))
			}
		}
		return 0
	}
}

// SetSessionExpiry exposes session expiration functionality to Lua scripts.
// It allows setting session values with custom expiration times directly from Lua code.
//
// Expected Lua parameters:
//  1. key (string): The session key to store
//  2. value (any): The value to store (supports string, number, nil, boolean, table)
//  3. expiry (number): Expiration time in seconds
//
// Usage in Lua:
//
//	eocto.setSessionExpiry("user_token", "abc123", 3600) -- expires in 1 hour
//
// Returns nothing to Lua state but logs any errors encountered
func SetSessionExpiry(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		sess, err := middleware.Store.Get(c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			return 0
		}

		key := L.ToString(1)
		var value interface{}
		expirySeconds := L.ToNumber(3)

		switch L.Get(2).Type() {
		case lua.LTString:
			value = L.ToString(2)
		case lua.LTNumber:
			value = float64(L.ToNumber(2))
		case lua.LTNil:
			value = nil
		case lua.LTBool:
			value = L.ToBool(2)
		case lua.LTTable:
			mapValue := make(map[string]interface{})
			L.ToTable(2).ForEach(func(k, v lua.LValue) {
				mapValue[k.String()] = v.String()
			})
			value = mapValue
		}

		// log.Printf("Setting session key %s to value %+v with expiry %v seconds", key, value, expirySeconds)
		sess.Set(key, value)
		sess.SetExpiry(time.Duration(expirySeconds) * time.Second)
		sess.Fresh()

		err = sess.Save()
		if err != nil {
			log.Printf("Error saving session: %v", err)
			return 0
		}

		return 0
	}
}

// DeleteSession returns a Lua function that deletes a value from the session.
// The function is exposed to Lua scripts and allows deletion of HTTP session data.
//
// Parameters passed from Lua:
//   - key: string - The session key to delete
//
// Returns to Lua:
//   - success: bool - Whether the delete operation was successful
//   - error: string - The error message if the delete operation failed
//
// Usage in Lua:
//
//	success, err = deleteSession("user_id")
//
// The function handles:
//   - Session retrieval and validation
//   - Value deletion
//   - Session persistence
//   - Error logging
func DeleteSession(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		sess, err := middleware.Store.Get(c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			L.Push(lua.LBool(false))
			L.Push(lua.LString(err.Error()))
			return 2
		}

		key := L.ToString(1)
		if key == "" {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("key cannot be empty"))
			return 2
		}

		sess.Delete(key)
		err = sess.Save()
		if err != nil {
			log.Printf("Error saving session after delete: %v", err)
			L.Push(lua.LBool(false))
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LBool(true))
		return 1
	}
}

// func GetAllSessions(L *lua.LState) int {
// 	store := middleware.Store
// 	tbl := L.NewTable()

// 	// Use the Storage interface methods
// 	keys, err := store.Storage.Keys()
// 	if err != nil {
// 		L.Push(lua.LNil)
// 		return 1
// 	}

// 	for _, key := range keys {
// 		value, err := store.Storage.Get(key)
// 		if err != nil {
// 			continue
// 		}
// 		tbl.RawSetString(key, lua.LString(string(value)))
// 	}

// 	L.Push(tbl)
// 	return 1
// }
