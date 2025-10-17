// Package utilities provides helper functions and utilities for the application.
//
// This file implements Fiber Locals integration with Lua scripts,
// allowing seamless access and modification of request-scoped local variables.
package utilities

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/degreane/octopus/internal/middleware"
	"github.com/degreane/octopus/internal/utilities/debug"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

func convertLuaTableToGo(t *lua.LTable) interface{} {
	goMap := make(map[string]interface{})
	goSlice := make([]interface{}, 0)

	isArray := true
	maxIndex := 0

	// Check if the table is an array (sequential integer keys starting from 1)
	t.ForEach(func(k, v lua.LValue) {
		if isArray {
			if num, ok := k.(lua.LNumber); ok {
				index := int(num)
				if index == maxIndex+1 {
					maxIndex = index
				} else {
					isArray = false
				}
			} else {
				isArray = false
			}
		}
	})

	// Convert to a Go slice if it's an array
	if isArray {
		t.ForEach(func(k, v lua.LValue) {
			goSlice = append(goSlice, convertLuaValueToGo(v))
		})
		return goSlice
	}

	// Otherwise, convert to a Go map
	t.ForEach(func(k, v lua.LValue) {
		goMap[k.String()] = convertLuaValueToGo(v)
	})
	return goMap
}

func convertLuaValueToGo(v lua.LValue) interface{} {
	switch v.Type() {
	case lua.LTString:
		return v.String()
	case lua.LTNumber:
		return float64(v.(lua.LNumber))
	case lua.LTBool:
		return bool(v.(lua.LBool))
	case lua.LTTable:
		// log.Printf("Lua Table ~~~~~~~~~~~~~~~~~~~~~~~~~~~\n%+v\n", v)
		return convertLuaTableToGo(v.(*lua.LTable))
	case lua.LTNil:
		return nil
	default:
		// log.Printf("Lua NIL ~~~~~~~~~~~~~~~~~~~~~~~~~~~\n%+v\n", v)
		return nil // Unsupported type
	}
}

// GetLocal returns a Lua function that retrieves values from Fiber Locals.
// The function is exposed to Lua scripts and provides access to request-scoped variables.
//
// Parameters passed from Lua:
//   - key: string - The locals key to retrieve
//
// Returns to Lua:
//   - string value if key exists
//   - nil if key doesn't exist
//
// Usage in Lua:
//
//	local value = getLocal("user_data")
func GetLocal(c *fiber.Ctx, basePath string) lua.LGFunction {
	return func(L *lua.LState) int {
		key := L.ToString(1)
		sess, err := middleware.Store.Get(c)
		if err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error getting session: %v", err))
			L.Push(lua.LNil)
			return 1
		}

		// Retrieve the JSON string stored under basePath + "_" + key
		raw := sess.Get(fmt.Sprintf("%s_%s", basePath, key))
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

// convertGoToLua converts generic Go values (from JSON) into lua.LValue
func convertGoToLua(L *lua.LState, v interface{}) lua.LValue {
	switch val := v.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(val)
	case float64:
		return lua.LNumber(val)
	case float32:
		return lua.LNumber(val)
	case int:
		return lua.LNumber(val)
	case int64:
		return lua.LNumber(val)
	case int32:
		return lua.LNumber(val)
	case json.Number:
		if f, err := val.Float64(); err == nil {
			return lua.LNumber(f)
		}
		return lua.LString(val.String())
	case string:
		return lua.LString(val)
	case []interface{}:
		t := L.NewTable()
		for i, item := range val {
			// Lua arrays are 1-based
			t.RawSetInt(i+1, convertGoToLua(L, item))
		}
		return t
	case map[string]interface{}:
		t := L.NewTable()
		for k, item := range val {
			t.RawSetString(k, convertGoToLua(L, item))
		}
		return t
	default:
		// Fallback to string representation
		return lua.LString(fmt.Sprintf("%v", val))
	}
}

func GetWsLocal(c *socketio.Websocket) lua.LGFunction {
	return func(L *lua.LState) int {

		key := L.ToString(1)
		val := c.Locals(key)
		if val == nil {
			log.Printf("Key %s not found in locals", key)
			L.Push(lua.LNil)
			return 1
		}
		L.Push(lua.LString(fmt.Sprintf("%v", val)))
		return 1
	}
}

// SetLocal returns a Lua function that stores values in Fiber Locals.
// The function is exposed to Lua scripts and allows modification of request-scoped variables.
//
// Parameters passed from Lua:
//   - key: string - The locals key to set
//   - value: string - The value to store
//
// Returns to Lua:
//   - true on successful operation
//
// Usage in Lua:
//
//	setLocal("user_data", "123")
func SetLocal(c *fiber.Ctx, basePath string) lua.LGFunction {
	return func(L *lua.LState) int {
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
			// mapValue := make(map[string]interface{})
			// L.ToTable(2).ForEach(func(k, v lua.LValue) {
			// 	mapValue[k.String()] = v.String()
			// })
			value = convertLuaTableToGo(L.ToTable(2))
		}
		// value := L.ToString(2)
		// debug.Debug(debug.Important, fmt.Sprintf("Setting local key %s to value %+v", key, value))
		c.Locals(key, value)
		sess, err := middleware.Store.Get(c)
		if err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error getting session: %v", err))
			//log.Printf("Error getting session: %v", err)
		} else {
			debug.Debug(debug.Info, fmt.Sprintf("setting  session: %s", fmt.Sprintf("%s_%s", basePath, key)))
			jsonBytes, mErr := json.Marshal(value)
			if mErr != nil {
				debug.Debug(debug.Error, fmt.Sprintf("Error marshaling value to JSON: %v", mErr))
			} else {
				sess.Set(fmt.Sprintf("%s_%s", basePath, key), string(jsonBytes))
				sess.Fresh()
				err = sess.Save()
				if err != nil {
					debug.Debug(debug.Error, fmt.Sprintf("Error saving session: %v", err))
				}
			}
		}

		L.Push(lua.LBool(true))
		return 1
	}
}

func SetWsLocal(c *socketio.Websocket) lua.LGFunction {
	return func(L *lua.LState) int {
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
			// mapValue := make(map[string]interface{})
			// L.ToTable(2).ForEach(func(k, v lua.LValue) {
			// 	mapValue[k.String()] = v.String()
			// })
			value = convertLuaTableToGo(L.ToTable(2))
		}
		// value := L.ToString(2)
		// debug.Debug(debug.Important, fmt.Sprintf("Setting local key %s to value %+v", key, value))
		c.Conn.Locals(key, value)
		L.Push(lua.LBool(true))
		return 1
	}
}

// DeleteLocal returns a Lua function that removes values from Fiber Locals.
// The function is exposed to Lua scripts and allows deletion of request-scoped variables.
//
// Parameters passed from Lua:
//   - key: string - The locals key to delete
//
// Returns to Lua:
//   - true if operation was successful
//   - false and error message if operation failed
//
// Usage in Lua:
//
//	local success, err = deleteLocal("user_data")
func DeleteLocal(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		key := L.ToString(1)
		if key == "" {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("key cannot be empty"))
			return 2
		}
		c.Locals(key, nil)
		L.Push(lua.LBool(true))
		return 1
	}
}
func DeleteWsLocal(c *socketio.Websocket) lua.LGFunction {
	return func(L *lua.LState) int {
		key := L.ToString(1)
		if key == "" {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("key cannot be empty"))
			return 2
		}
		c.Conn.Locals(key, nil)
		L.Push(lua.LBool(true))
		return 1
	}
}

// GetLocals returns a Lua function that retrieves all Fiber Locals as a table.
// The function is exposed to Lua scripts and provides access to all request-scoped variables.
//
// Parameters passed from Lua:
//   - none
//
// Returns to Lua:
//   - table containing all locals as key/value pairs
//
// Usage in Lua:
//
//	local allLocals = getLocals()
//	for k, v in pairs(allLocals) do
//	    print(k, v)
//	end
func GetLocals(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		localsTable := L.NewTable()

		c.Context().VisitUserValuesAll(func(a1, a2 any) {
			localsTable.RawSetString(fmt.Sprintf("%v", a1), lua.LString(fmt.Sprintf("%v", a2)))
		})
		L.Push(localsTable)
		return 1
	}
}

func GetWsLocals(c *socketio.Websocket) lua.LGFunction {
	return func(L *lua.LState) int {
		localsTable := L.NewTable()

		//c.Context().VisitUserValuesAll(func(a1, a2 any) {
		//	localsTable.RawSetString(fmt.Sprintf("%v", a1), lua.LString(fmt.Sprintf("%v", a2)))
		//})
		localsTable.RawSetString(fmt.Sprintf("%v", "ws"), lua.LString(fmt.Sprintf("%v", c.UUID)))
		L.Push(localsTable)
		return 1
	}
}
