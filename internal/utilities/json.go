package utilities

import (
	"encoding/json"
	"fmt"

	"github.com/degreane/octopus/internal/utilities/debug"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

// DecodeJSON attempts to decode a JSON string into a map.
// Returns the decoded map on success, nil on failure.
//
// Parameters:
//   - jsonString: string - The JSON string to decode
//
// Returns:
//   - map[string]interface{} - Decoded JSON object or nil if decoding fails
//
// Example:
//
//	result := DecodeJSON(`{"name": "John", "age": 30}`)
func DecodeJSON(jsonString string) interface{} {
	// Try to decode as array first
	var arrayResult []map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &arrayResult)
	if err == nil {
		// fmt.Printf("Decoded as array: %v\n", arrayResult)
		return arrayResult
	} else {
		// debug.Debug(debug.Error, fmt.Sprintf("We have level1 Error %+v\n%s", err, jsonString))
		// If not array, try as single object
		debug.Debug(debug.Info, "We have level1 Error, trying as object")
		var objectResult map[string]interface{}
		err = json.Unmarshal([]byte(jsonString), &objectResult)
		if err == nil {
			// fmt.Printf("Decoded as object: %v\n", objectResult)
			return objectResult
		} else {
			debug.Debug(debug.Error, fmt.Sprintf("We have level2 Error %+v\n%s", err, jsonString))
		}
	}

	return nil
}
func handleMapValue(L *lua.LState, v map[string]interface{}) *lua.LTable {
	table := L.NewTable()
	for k, val := range v {
		switch innerVal := val.(type) {
		case string:
			table.RawSetString(k, lua.LString(innerVal))
		case float64:
			table.RawSetString(k, lua.LNumber(innerVal))
		case bool:
			table.RawSetString(k, lua.LBool(innerVal))
		case map[string]interface{}:
			// Recursively handle nested maps
			nestedTable := handleMapValue(L, innerVal)
			table.RawSetString(k, nestedTable)
		case []interface{}:
			// Handle slices (arrays)
			nestedTable := handleSliceValue(L, innerVal)
			table.RawSetString(k, nestedTable)
		case nil:
			// Handle nil values
			table.RawSetString(k, lua.LNil)
		default:
			// Handle unsupported types (optional: log or ignore)
			table.RawSetString(k, lua.LNil)
		}
	}
	return table
}

// Helper function to handle slices (arrays)
func handleSliceValue(L *lua.LState, slice []interface{}) *lua.LTable {
	table := L.NewTable()
	for i, val := range slice {
		switch innerVal := val.(type) {
		case string:
			table.RawSetInt(i+1, lua.LString(innerVal)) // Lua arrays are 1-indexed
		case float64:
			table.RawSetInt(i+1, lua.LNumber(innerVal))
		case bool:
			table.RawSetInt(i+1, lua.LBool(innerVal))
		case map[string]interface{}:
			// Recursively handle nested maps
			nestedTable := handleMapValue(L, innerVal)
			table.RawSetInt(i+1, nestedTable)
		case []interface{}:
			// Recursively handle nested slices
			nestedTable := handleSliceValue(L, innerVal)
			table.RawSetInt(i+1, nestedTable)
		case nil:
			// Handle nil values
			table.RawSetInt(i+1, lua.LNil)
		default:
			// Handle unsupported types (optional: log or ignore)
			table.RawSetInt(i+1, lua.LNil)
		}
	}
	return table
}

// GetDecodeJSON returns a Lua function that decodes JSON strings.
// The function is exposed to Lua scripts and provides JSON decoding functionality.
//
// Parameters passed from Lua:
//   - jsonString: string - The JSON string to decode
//
// Returns to Lua:
//   - table - Decoded JSON as Lua table on success
//   - nil - If decoding fails
//
// Usage in Lua:
//
//	local jsonTable = decodeJSON('{"name": "John", "age": 30}')
func GetDecodeJSON(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		jsonString := L.ToString(1)
		decoded := DecodeJSON(jsonString)

		if decoded == nil {
			L.Push(lua.LNil)
			return 1
		}

		// Handle array type
		if arr, ok := decoded.([]map[string]interface{}); ok {
			table := L.NewTable()
			for i, item := range arr {
				itemTable := handleMapValue(L, item)
				table.RawSetInt(i+1, itemTable)
			}
			L.Push(table)
			return 1
		}

		// Handle object type
		if obj, ok := decoded.(map[string]interface{}); ok {
			table := handleMapValue(L, obj)
			debug.Debug(debug.Info, "Mapping Object")
			// log.Printf("%+v", table)
			L.Push(table)
			return 1
		}

		L.Push(lua.LNil)
		return 1
	}
}

// func GetEncodeJSON(c *fiber.Ctx) lua.LGFunction {
// 	return func(L *lua.LState) int {
// 		// Get the table from first argument
// 		tbl := L.CheckTable(1)

// 		// Convert Lua table to Go map
// 		data := make(map[string]interface{})
// 		tbl.ForEach(func(k, v lua.LValue) {
// 			switch v.Type() {
// 			case lua.LTString:
// 				data[k.String()] = v.String()
// 			case lua.LTNumber:
// 				data[k.String()] = float64(v.(lua.LNumber))
// 			case lua.LTBool:
// 				data[k.String()] = bool(v.(lua.LBool))
// 			case lua.LTTable:
// 				// Handle nested tables recursively
// 				nestedTbl := make(map[string]interface{})
// 				v.(*lua.LTable).ForEach(func(nk, nv lua.LValue) {
// 					nestedTbl[nk.String()] = nv
// 				})
// 				data[k.String()] = nestedTbl
// 			}
// 		})

// 		// Convert to JSON
// 		jsonBytes, err := json.Marshal(data)
// 		if err != nil {
// 			L.Push(lua.LNil)
// 			L.Push(lua.LString(err.Error()))
// 			return 2
// 		}

// 		L.Push(lua.LString(string(jsonBytes)))
// 		return 1
// 	}
// }

func GetEncodeJSON(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		// Get the table from first argument
		tbl := L.CheckTable(1)

		// Convert Lua table to Go value (handling arrays/objects recursively)
		data := convertLuaTable(L, tbl)

		// Convert to JSON
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LString(string(jsonBytes)))
		return 1
	}
}

// Recursive function to convert Lua values to Go values
func convertLuaValue(L *lua.LState, lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case lua.LString:
		return v.String()
	case lua.LNumber:
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LTable:
		return convertLuaTable(L, v)
	default:
		return v.String()
	}
}

// Converts Lua table to appropriate Go type (slice or map)
func convertLuaTable(L *lua.LState, tbl *lua.LTable) interface{} {
	// Check if table is array-like (sequential integer keys starting at 1)
	maxIndex := 0
	isArray := true
	count := 0

	tbl.ForEach(func(k, v lua.LValue) {
		if k.Type() == lua.LTNumber {
			num := float64(k.(lua.LNumber))
			if num == float64(int(num)) {
				index := int(num)
				if index == maxIndex+1 {
					maxIndex = index
					count++
				} else {
					isArray = false
				}
			} else {
				isArray = false
			}
		} else {
			isArray = false
		}
	})

	// Handle array-like tables
	if isArray && count > 0 {
		array := make([]interface{}, maxIndex)
		for i := 1; i <= maxIndex; i++ {
			val := tbl.RawGetInt(i)
			array[i-1] = convertLuaValue(L, val)
		}
		return array
	}

	// Handle object-like tables
	obj := make(map[string]interface{})
	tbl.ForEach(func(k, v lua.LValue) {
		key := k.String()
		obj[key] = convertLuaValue(L, v)
	})
	return obj
}
