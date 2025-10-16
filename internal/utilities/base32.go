package utilities

import (
	"encoding/base32"

	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

func GetEncodeBase32(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		input := L.ToString(1)
		encoded := base32.StdEncoding.EncodeToString([]byte(input))
		L.Push(lua.LString(encoded))
		return 1
	}
}

func GetDecodeBase32(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		input := L.ToString(1)
		decoded, err := base32.StdEncoding.DecodeString(input)
		if err != nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(lua.LString(string(decoded)))
		return 1
	}
}
