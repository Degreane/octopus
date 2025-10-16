package utilities

import (
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

// GetResponse returns a Lua function that stores HTTP response data in Fiber Locals.
//
// Parameters passed from Lua:
//   - status: number - HTTP status code
//   - body: table - Response body data
//
// Usage in Lua:
//
//	response(500, {error = "Internal Server Error"})
//	response(200, {data = "Success"})
func GetResponse(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		status := L.ToInt(1)
		body := L.ToTable(2)

		responseMap := fiber.Map{}
		body.ForEach(func(k, v lua.LValue) {
			responseMap[k.String()] = v.String()
		})

		c.Locals("lua_response", map[int]fiber.Map{
			status: responseMap,
		})
		return 0
	}
}
