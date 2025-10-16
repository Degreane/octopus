package utilities

import (
	"encoding/json"
	"fmt"

	"github.com/degreane/octopus/internal/utilities/debug"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

// GetPostBody returns a Lua function that extracts POST body data from the Fiber context
// and converts it into a Lua table. It handles multiple body formats:
//
// 1. JSON body - Unmarshals JSON data into the Lua table
// 2. Multipart form data - Extracts form values into the Lua table
// 3. Regular POST form data - Processes standard POST parameters
//
// Parameters passed to the Lua function:
//   - None required
//
// Returns to Lua:
//   - table: A Lua table containing all body data with string keys and values
//
// Usage in Lua:
//
//	local bodyData = eocto.getPostBody()
// func GetPostBody(c *fiber.Ctx) lua.LGFunction {
// return func(L *lua.LState) int {
// 	// Create Lua table for body data
// 	bodyTable := L.NewTable()

// 	// Handle JSON body
// 	var jsonData map[string]interface{}
// 	// debug.Debug(debug.Important, string(c.Body()))
// 	err := json.Unmarshal(c.Body(), &jsonData)
// 	if err == nil {
// 		debug.Debug(debug.Info, "Json Body")
// 		for k, v := range jsonData {
// 			bodyTable.RawSetString(k, lua.LString(fmt.Sprintf("%v", v)))
// 		}
// 	}

// 	// Handle form data
// 	form, err := c.MultipartForm()
// 	if err == nil {
// 		debug.Debug(debug.Info, "Form Body")
// 		for k, v := range form.Value {
// 			if len(v) > 0 {
// 				bodyTable.RawSetString(k, lua.LString(v[0]))
// 			}
// 		}
// 	}

// 	// Handle regular POST form
// 	c.Request().PostArgs().VisitAll(func(key, value []byte) {
// 		bodyTable.RawSetString(string(key), lua.LString(string(value)))
// 	})

//		L.Push(bodyTable)
//		return 1
//	}
//
// }
func GetPostBody(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		bodyTable := L.NewTable()
		rawBody := c.Body()

		// First try to parse as JSON (either object or array)
		var jsonValue interface{}
		if err := json.Unmarshal(rawBody, &jsonValue); err == nil {
			debug.Debug(debug.Info, "JSON Body")
			luaValue := convertToLua(L, jsonValue)
			if tbl, ok := luaValue.(*lua.LTable); ok {
				// If it's a table, use it directly
				L.Push(tbl)
			} else {
				// For simple values, store in a table
				bodyTable.RawSetString("_value", luaValue)
				L.Push(bodyTable)
			}
			return 1
		}

		// Handle form data if JSON parsing failed
		form, err := c.MultipartForm()
		if err == nil {
			debug.Debug(debug.Info, "Form Body")
			for k, v := range form.Value {
				if len(v) > 0 {
					bodyTable.RawSetString(k, lua.LString(v[0]))
				}
			}
		}

		// Handle regular POST form
		c.Request().PostArgs().VisitAll(func(key, value []byte) {
			bodyTable.RawSetString(string(key), lua.LString(string(value)))
		})

		L.Push(bodyTable)
		return 1
	}
}

// Convert Go interface to Lua value (recursive)
func convertToLua(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case map[string]interface{}:
		tbl := L.NewTable()
		for key, val := range v {
			tbl.RawSetString(key, convertToLua(L, val))
		}
		return tbl
	case []interface{}:
		tbl := L.NewTable()
		for i, item := range v {
			tbl.RawSetInt(i+1, convertToLua(L, item)) // Lua arrays are 1-indexed
		}
		return tbl
	case string:
		return lua.LString(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case nil:
		return lua.LNil
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}
