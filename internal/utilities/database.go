package utilities

import (
	"github.com/degreane/octopus/internal/database"
	lua "github.com/yuin/gopher-lua"
)

// GetDataFromCollectionLua is a Lua binding function that retrieves data from a MongoDB collection
// based on the provided URI, database name, collection name, and filter. It converts the Lua
// filter table to a Go interface and returns the result as a string or an error.
// The function takes a Lua state and returns the number of return values (1 for success, 2 for error).
func GetDataFromCollectionLua(L *lua.LState) int {
	uri := L.CheckString(1)
	dbName := L.CheckString(2)
	collectionName := L.CheckString(3)
	filter := L.CheckTable(4)

	// Convert Lua table to Go interface
	var filterData interface{}
	filterData = tableToInterface(filter)

	result, err := database.GetDataFromCollection(uri, dbName, collectionName, filterData)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(result))
	return 1
}

// SetDataToCollectionLua is a Lua binding function that updates data in a MongoDB collection
// based on the provided URI, database name, collection name, filter, and update data.
// It converts the Lua filter and data tables to Go interfaces and returns the result
// as a string or an error. The function takes a Lua state and returns the number of
// return values (1 for success, 2 for error).
func SetDataToCollectionLua(L *lua.LState) int {
	uri := L.CheckString(1)
	dbName := L.CheckString(2)
	collectionName := L.CheckString(3)
	filter := L.CheckTable(4)
	data := L.CheckTable(5)

	filterData := tableToInterface(filter)
	updateData := tableToInterface(data)

	result, err := database.SetDataToCollection(uri, dbName, collectionName, filterData, updateData)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(result))
	return 1
}

// DelDataFromCollectionLua is a Lua binding function that deletes data from a MongoDB collection
// based on the provided URI, database name, collection name, and filter.
// It converts the Lua filter table to a Go interface and returns the result
// as a string or an error. The function takes a Lua state and returns the number of
// return values (1 for success, 2 for error).
func DelDataFromCollectionLua(L *lua.LState) int {
	uri := L.CheckString(1)
	dbName := L.CheckString(2)
	collectionName := L.CheckString(3)
	filter := L.CheckTable(4)

	filterData := tableToInterface(filter)

	result, err := database.DelDataFromCollection(uri, dbName, collectionName, filterData)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(result))
	return 1
}

// InsertDataToCollectionLua is a Lua binding function that inserts data into a MongoDB collection
// based on the provided URI, database name, collection name, and data.
// It converts the Lua data table to a Go interface and returns the result
// as a string or an error. The function takes a Lua state and returns the number of
// return values (1 for success, 2 for error).
func InsertDataToCollectionLua(L *lua.LState) int {
	uri := L.CheckString(1)
	dbName := L.CheckString(2)
	collectionName := L.CheckString(3)
	data := L.CheckTable(4)

	insertData := tableToInterface(data)

	result, err := database.InsertDataToCollection(uri, dbName, collectionName, insertData)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(result))
	return 1
}

// Helper function to convert Lua tables to Go interface{}
// tableToInterface converts a Lua table to a Go interface{} by recursively
// transforming Lua table values into their corresponding Go types. It supports
// converting string, number, boolean, and nested table values, creating a
// map[string]interface{} representation of the original Lua table.
func tableToInterface(table *lua.LTable) interface{} {
	result := make(map[string]interface{})

	table.ForEach(func(k, v lua.LValue) {
		switch v.Type() {
		case lua.LTString:
			result[k.String()] = v.String()
		case lua.LTNumber:
			result[k.String()] = float64(v.(lua.LNumber))
		case lua.LTBool:
			result[k.String()] = bool(v.(lua.LBool))
		case lua.LTTable:
			result[k.String()] = tableToInterface(v.(*lua.LTable))
		}
	})

	return result
}
