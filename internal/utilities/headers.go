package utilities

import (
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

// GetHeaders returns a Lua function that retrieves all request headers.
// Returns to Lua: table containing all headers as key/value pairs
func GetHeaders(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		headersTable := L.NewTable()
		c.Request().Header.VisitAll(func(key, value []byte) {
			// fmt.Printf("Header: %s = %s\n", key, value)
			headersTable.RawSetString(string(key), lua.LString(string(value)))
		})
		L.Push(headersTable)
		return 1
	}
}

// GetHeader returns a Lua function that retrieves a specific header value.
// Parameters from Lua: key (string) - The header key to retrieve
// Returns to Lua: string value if exists, nil if not
func GetHeader(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		key := L.ToString(1)
		value := c.Get(key)
		if value == "" {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(lua.LString(value))
		return 1
	}
}

// SetHeader returns a Lua function that sets a header value.
// Parameters from Lua:
//   - key: string - The header key to set
//   - value: string - The value to set
func SetHeader(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		key := L.ToString(1)
		value := L.ToString(2)
		c.Set(key, value)
		return 0
	}
}

// DeleteHeader returns a Lua function that removes a header.
// Parameters from Lua: key (string) - The header key to delete
func DeleteHeader(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		key := L.ToString(1)
		c.Request().Header.Del(key)
		return 0
	}
}

// GetMethod returns a Lua function that retrieves the HTTP method of the current request.
// Returns to Lua: string containing the HTTP method (GET, POST, PUT, DELETE, etc.)
//
// Usage in Lua:
//
//	local method = getMethod()
//	print(method) -- "GET", "POST", etc.
func GetMethod(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		L.Push(lua.LString(c.Method()))
		return 1
	}
}

// GetPath returns a Lua function that retrieves the request path.
// Returns to Lua: string containing the request path
//
// Usage in Lua:
//
//	local path = getPath()
func GetPath(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		L.Push(lua.LString(c.Path()))
		return 1
	}
}

// GetHost returns a Lua function that retrieves the request host.
// Returns to Lua: string containing the request host
//
// Usage in Lua:
//
//	local host = getHost()
func GetHost(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		L.Push(lua.LString(c.Hostname()))
		return 1
	}
}

// GetSchema returns a Lua function that retrieves the request schema (http/https).
// Returns to Lua: string containing the request schema
//
// Usage in Lua:
//
//	local schema = getSchema()
func GetSchema(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		L.Push(lua.LString(c.Protocol()))
		return 1
	}
}
