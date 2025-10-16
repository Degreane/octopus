package utilities

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

// WsAddRoom adds a room to a user's WebSocket client attributes
func WsAddRoom(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		// Get parameters from Lua
		userId := L.CheckString(1) // First parameter: userId
		roomId := L.CheckString(2) // Second parameter: roomId

		// Validate parameters
		if userId == "" {
			log.Printf("WsAddRoom: userId cannot be empty")
			L.Push(lua.LBool(false))
			L.Push(lua.LString("userId cannot be empty"))
			return 2
		}

		if roomId == "" {
			log.Printf("WsAddRoom: roomId cannot be empty")
			L.Push(lua.LBool(false))
			L.Push(lua.LString("roomId cannot be empty"))
			return 2
		}

		// Get global socket clients instance
		clients := GetSocketClients()

		// Check if user exists in clients
		_, exists := clients.GetClient(userId)
		if !exists {
			log.Printf("WsAddRoom: User %s not found in clients", userId)
			L.Push(lua.LBool(false))
			L.Push(lua.LString("user not found"))
			return 2
		}

		// Get existing rooms or create new slice
		var rooms []string
		if existingRoomsInterface, ok := clients.GetClientAttribute(userId, "rooms"); ok {
			if existingRooms, ok := existingRoomsInterface.([]string); ok {
				rooms = existingRooms
			} else {
				// Handle case where rooms metadata exists but is not []string
				rooms = []string{}
			}
		} else {
			// No rooms metadata exists, create new slice
			rooms = []string{}
		}

		// Check if room already exists for this user
		for _, existingRoom := range rooms {
			if existingRoom == roomId {
				log.Printf("WsAddRoom: User %s already in room %s", userId, roomId)
				L.Push(lua.LBool(true))
				L.Push(lua.LString("user already in room"))
				return 2
			}
		}

		// Add the room to the user's rooms
		rooms = append(rooms, roomId)

		// Set the updated rooms metadata
		clients.SetClientAttribute(userId, "rooms", rooms)

		// Also set individual room metadata for easier lookup
		roomKey := "room_" + roomId
		clients.SetClientAttribute(userId, roomKey, true)

		log.Printf("WsAddRoom: Successfully added user %s to room %s", userId, roomId)

		// Return success
		L.Push(lua.LBool(true))
		L.Push(lua.LString("user added to room successfully"))
		return 2
	}
}

// WsRemoveRoom removes a room from a user's WebSocket client attributes
func WsRemoveRoom(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		userId := L.CheckString(1)
		roomId := L.CheckString(2)

		if userId == "" || roomId == "" {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("userId and roomId cannot be empty"))
			return 2
		}

		clients := GetSocketClients()

		// Check if user exists
		_, exists := clients.GetClient(userId)
		if !exists {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("user not found"))
			return 2
		}

		// Get existing rooms
		var rooms []string
		if existingRoomsInterface, ok := clients.GetClientAttribute(userId, "rooms"); ok {
			if existingRooms, ok := existingRoomsInterface.([]string); ok {
				rooms = existingRooms
			}
		}

		// Remove the room from the slice
		var newRooms []string
		roomFound := false
		for _, existingRoom := range rooms {
			if existingRoom != roomId {
				newRooms = append(newRooms, existingRoom)
			} else {
				roomFound = true
			}
		}

		if !roomFound {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("user not in room"))
			return 2
		}

		// Update rooms metadata
		clients.SetClientAttribute(userId, "rooms", newRooms)

		// Remove individual room metadata
		roomKey := "room_" + roomId
		clients.SetClientAttribute(userId, roomKey, nil)

		log.Printf("WsRemoveRoom: Successfully removed user %s from room %s", userId, roomId)

		L.Push(lua.LBool(true))
		L.Push(lua.LString("user removed from room successfully"))
		return 2
	}
}

// WsGetUserRooms gets all rooms for a user
func WsGetUserRooms(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		userId := L.CheckString(1)

		if userId == "" {
			L.Push(lua.LNil)
			L.Push(lua.LString("userId cannot be empty"))
			return 2
		}

		clients := GetSocketClients()

		// Check if user exists
		_, exists := clients.GetClient(userId)
		if !exists {
			L.Push(lua.LNil)
			L.Push(lua.LString("user not found"))
			return 2
		}

		// Get rooms metadata
		var rooms []string
		if existingRoomsInterface, ok := clients.GetClientAttribute(userId, "rooms"); ok {
			if existingRooms, ok := existingRoomsInterface.([]string); ok {
				rooms = existingRooms
			}
		}

		// Convert to Lua table
		roomsTable := L.NewTable()
		for i, room := range rooms {
			roomsTable.RawSetInt(i+1, lua.LString(room))
		}

		L.Push(roomsTable)
		L.Push(lua.LString("success"))
		return 2
	}
}

// WsIsUserInRoom checks if a user is in a specific room
func WsIsUserInRoom(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		userId := L.CheckString(1)
		roomId := L.CheckString(2)

		if userId == "" || roomId == "" {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("userId and roomId cannot be empty"))
			return 2
		}

		clients := GetSocketClients()

		// Check if user exists
		_, exists := clients.GetClient(userId)
		if !exists {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("user not found"))
			return 2
		}

		// Check individual room metadata for faster lookup
		roomKey := "room_" + roomId
		if roomValue, ok := clients.GetClientAttribute(userId, roomKey); ok {
			if inRoom, ok := roomValue.(bool); ok && inRoom {
				L.Push(lua.LBool(true))
				L.Push(lua.LString("user is in room"))
				return 2
			}
		}

		L.Push(lua.LBool(false))
		L.Push(lua.LString("user not in room"))
		return 2
	}
}

// WSEmitToRoom emits a message to all users in a specific room
func WSEmitToRoom(roomId string, event string, data interface{}, excludeUsers ...string) int {
	clients := GetSocketClients()
	allClients := clients.GetClients()

	excludeMap := make(map[string]bool)
	for _, userID := range excludeUsers {
		excludeMap[userID] = true
	}

	successCount := 0
	roomKey := "room_" + roomId

	for userID, client := range allClients {
		// Skip excluded users
		if excludeMap[userID] {
			continue
		}

		// Check if user is in the room
		if roomValue, ok := clients.GetClientAttribute(userID, roomKey); ok {
			if inRoom, ok := roomValue.(bool); ok && inRoom && client.Socket != nil {
				if sendDirectMessage(client.Socket, event, data, userID) {
					successCount++
				}
			}
		}
	}

	log.Printf("Message emitted to room %s: %d clients received message", roomId, successCount)
	return successCount
}

func sendDirectMessage(socket *socketio.Websocket, event string, data interface{}, userID string) bool {
	// Convert data to string if it's not already
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	case []byte:
		dataStr = string(v)
	default:
		// Marshal to JSON string if it's a complex type
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			log.Printf("Error marshaling data for user %s: %v", userID, err)
			return false
		}
		dataStr = string(jsonBytes)
	}

	// Create the message using your MessageObject struct
	message := MessageObject{
		Data:  dataStr,
		From:  "server",
		Event: event,
		To:    userID,
	}

	// Marshal the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling MessageObject for user %s: %v", userID, err)
		return false
	}

	// Check if socket is valid
	if socket == nil {
		log.Printf("Socket is nil for user %s", userID)
		return false
	}

	// Send the message via WebSocket
	socket.Emit(jsonData, socketio.TextMessage)

	log.Printf("Message successfully sent to user %s: %s", userID, event)
	return true

}

// WsEmitToRoom Lua function wrapper for WSEmitToRoom
func WsEmitToRoom(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		roomId := L.CheckString(1)
		event := L.CheckString(2)

		// Get data parameter (can be string, number, table, etc.)
		var data interface{}
		dataValue := L.Get(3)
		switch dataValue.Type() {
		case lua.LTString:
			data = dataValue.String()
		case lua.LTNumber:
			data = float64(dataValue.(lua.LNumber))
		case lua.LTBool:
			data = bool(dataValue.(lua.LBool))
		case lua.LTTable:
			// Convert Lua table to Go map
			data = luaTableToMap2(L, dataValue.(*lua.LTable))
		case lua.LTNil:
			data = nil
		default:
			data = dataValue.String()
		}

		if roomId == "" || event == "" {
			L.Push(lua.LNumber(0))
			L.Push(lua.LString("roomId and event cannot be empty"))
			return 2
		}

		// Optional exclude users parameter
		var excludeUsers []string
		if L.GetTop() >= 4 {
			excludeTable := L.Get(4)
			if excludeTable.Type() == lua.LTTable {
				table := excludeTable.(*lua.LTable)
				table.ForEach(func(key, value lua.LValue) {
					if value.Type() == lua.LTString {
						excludeUsers = append(excludeUsers, value.String())
					}
				})
			}
		}

		count := WSEmitToRoom(roomId, event, data, excludeUsers...)

		L.Push(lua.LNumber(count))
		L.Push(lua.LString("success"))
		return 2
	}
}

// Helper function to convert Lua table to Go map
func luaTableToMap2(L *lua.LState, table *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	table.ForEach(func(key, value lua.LValue) {
		keyStr := key.String()
		switch value.Type() {
		case lua.LTString:
			result[keyStr] = value.String()
		case lua.LTNumber:
			result[keyStr] = float64(value.(lua.LNumber))
		case lua.LTBool:
			result[keyStr] = bool(value.(lua.LBool))
		case lua.LTTable:
			result[keyStr] = luaTableToMap2(L, value.(*lua.LTable))
		case lua.LTNil:
			result[keyStr] = nil
		default:
			result[keyStr] = value.String()
		}
	})
	return result
}
