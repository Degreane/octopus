package utilities

import (
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

// GetRenderJson returns a Lua function that converts a Lua table to JSON and sends it via Fiber
func GetRenderJson(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		// Get the Lua table from the first argument
		luaTable := L.CheckTable(1)

		// Optional: Get status code from second argument (default to 200)
		statusCode := 200
		if L.GetTop() >= 2 {
			statusCode = L.CheckInt(2)
		}

		// Convert Lua table to Go map
		jsonData := luaTableToMap(luaTable)

		// Send JSON response via Fiber
		err := c.Status(statusCode).JSON(jsonData)
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// Mark that response was rendered from Lua to prevent default template rendering
		c.Locals("rendered_from_lua", true)

		L.Push(lua.LBool(true))
		return 1
	}
}

// Helper function to convert Lua table to Go map recursively
func luaTableToMap(table *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})

	table.ForEach(func(key, value lua.LValue) {
		keyStr := key.String()

		switch v := value.(type) {
		case *lua.LTable:
			// Recursively handle nested tables
			result[keyStr] = luaTableToMap(v)
		case lua.LString:
			result[keyStr] = string(v)
		case lua.LNumber:
			result[keyStr] = float64(v)
		case lua.LBool:
			result[keyStr] = bool(v)
		case *lua.LNilType:
			result[keyStr] = nil
		default:
			result[keyStr] = value.String()
		}
	})

	return result
}
