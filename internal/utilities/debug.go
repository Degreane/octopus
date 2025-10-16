package utilities

import (
	"fmt"
	"strings"

	"github.com/degreane/octopus/internal/utilities/debug"
	lua "github.com/yuin/gopher-lua"
)

// tableToString converts a Lua table to a string representation
func tableToString(L *lua.LState, table *lua.LTable) string {
	var parts []string
	table.ForEach(func(key, value lua.LValue) {
		keyStr := key.String()
		var valueStr string

		switch value.Type() {
		case lua.LTTable:
			valueStr = tableToString(L, value.(*lua.LTable))
		default:
			valueStr = value.String()
		}

		parts = append(parts, fmt.Sprintf("%s: %s", keyStr, valueStr))
	})

	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}
func Debug(L *lua.LState) int {
	level := L.CheckString(1)
	msgValue := L.Get(2) // Get the value without type checking

	var msg string

	// Handle different types of Lua values
	switch msgValue.Type() {
	case lua.LTString:
		msg = msgValue.String()
	case lua.LTNumber:
		msg = msgValue.String()
	case lua.LTBool:
		msg = msgValue.String()
	case lua.LTTable:
		// Convert table to JSON-like string representation
		msg = tableToString(L, msgValue.(*lua.LTable))
	case lua.LTNil:
		msg = "nil"
	default:
		msg = msgValue.String()
	}

	switch level {
	case "info":
		debug.Debug(debug.Info, msg)
	case "warning":
		debug.Debug(debug.Warning, msg)
	case "error":
		debug.Debug(debug.Error, msg)
	case "important":
		debug.Debug(debug.Important, msg)
	}
	return 0
}
