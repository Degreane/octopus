package utilities

import (
	"fmt"

	"github.com/degreane/octopus/internal/utilities/debug"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

func GetRender(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		view := L.ToString(1)
		data := L.ToTable(2)

		// Convert Lua table to fiber.Map
		bindData := fiber.Map{}
		data.ForEach(func(k, v lua.LValue) {
			switch v.Type() {
			case lua.LTString:
				bindData[k.String()] = v.String()
			case lua.LTNumber:
				bindData[k.String()] = float64(v.(lua.LNumber))
			case lua.LTBool:
				bindData[k.String()] = bool(v.(lua.LBool))
			case lua.LTTable:
				nestedMap := luaTableToMap(v.(*lua.LTable))
				//	make(map[string]interface{})
				//
				//v.(*lua.LTable).ForEach(func(nk, nv lua.LValue) {
				//	nestedMap[nk.String()] = nv
				//})
				bindData[k.String()] = nestedMap
			}
		})
		// fmt.Printf("Bind data: %v\n", bindData)
		debug.Debug(debug.Important, fmt.Sprintf("View: %s\n", view))
		// fmt.Printf("View: %s\n", view)
		err := c.Render(view, bindData)
		if err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error rendering view: %v\n", err))
			L.Push(lua.LBool(false))
			L.Push(lua.LString(err.Error()))
			return 2
		}
		c.Locals("rendered_from_lua", true)
		L.Push(lua.LBool(true))
		return 1
	}
}
