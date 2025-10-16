package utilities

import (
	"log"

	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

// GetQueryParams returns a Lua function that retrieves URL query parameters.
// The function is exposed to Lua scripts and provides access to request query parameters.
//
// Returns to Lua:
//   - table containing all query parameters as key/value pairs
//
// Usage in Lua:
//
//	local params = getQueryParams()
//	for k, v in pairs(params) do
//	    print(k, v)
//	end
func GetQueryParams(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		queryTable := L.NewTable()

		c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
			queryTable.RawSetString(string(key), lua.LString(string(value)))
		})

		L.Push(queryTable)
		return 1
	}
}

// GetPathParams returns a Lua function that extracts path parameters from the Fiber request
func GetPathParams(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		// Create a new Lua table to hold path parameters
		pathParamsTable := L.NewTable()

		// Get all path parameters from Fiber context
		allParams := c.AllParams()

		// Add each path parameter to the Lua table
		for key, value := range allParams {
			pathParamsTable.RawSetString(key, lua.LString(value))
		}

		// Push the table onto the Lua stack
		L.Push(pathParamsTable)
		return 1 // Return 1 value (the table)
	}
}

// GetPathParam returns a Lua function that gets a specific path parameter
func GetPathParam(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		// Get the parameter name from Lua
		paramName := L.CheckString(1)
		log.Printf("paramName: %s", paramName)
		// Get the parameter value from Fiber context
		paramValue := c.Params(paramName)
		log.Printf("paramValue: %s", paramValue)

		// Push the value onto the Lua stack
		L.Push(lua.LString(paramValue))
		return 1 // Return 1 value (the parameter value)
	}
}
