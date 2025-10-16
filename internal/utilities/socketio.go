package utilities

import (
	"sync"
	"time"

	"github.com/gofiber/contrib/socketio"
)

// Core connection lifecycle events
const (
	// Connection Events
	EventConnect    = "connect"
	EventDisconnect = "disconnect"
	EventReconnect  = "reconnect"
	EventConnecting = "connecting"
	EventError      = "connect_error"
	EventTimeout    = "connect_timeout"
)

const (
	// Room Management
	EventJoinRoom    = "join_room"
	EventLeaveRoom   = "leave_room"
	EventCreateRoom  = "create_room"
	EventDeleteRoom  = "delete_room"
	EventRoomJoined  = "room_joined"
	EventRoomLeft    = "room_left"
	EventRoomCreated = "room_created"
	EventRoomDeleted = "room_deleted"
	EventRoomError   = "room_error"

	// Room Information
	EventRoomUsers = "room_users"
	EventRoomList  = "room_list"
	EventRoomInfo  = "room_info"
)

const (
	// User Presence
	EventUserJoined        = "user_joined"
	EventUserLeft          = "user_left"
	EventUserOnline        = "user_online"
	EventUserOffline       = "user_offline"
	EventUserTyping        = "user_typing"
	EventUserStoppedTyping = "user_stopped_typing"

	// User Actions
	EventUserUpdate = "user_update"
	EventUserStatus = "user_status"
	EventUserList   = "user_list"
)

const (
	// Direct Messaging
	EventMessage         = "message"
	EventPrivateMessage  = "private_message"
	EventBroadcast       = "broadcast"
	EventMessageSent     = "message_sent"
	EventMessageReceived = "message_received"
	EventMessageRead     = "message_read"
	EventMessageError    = "message_error"

	// Message Types
	EventTextMessage  = "text_message"
	EventImageMessage = "image_message"
	EventFileMessage  = "file_message"
	EventVoiceMessage = "voice_message"
	EventVideoMessage = "video_message"
)
const (
	// Namespace Management
	EventJoinNamespace   = "join_namespace"
	EventLeaveNamespace  = "leave_namespace"
	EventNamespaceJoined = "namespace_joined"
	EventNamespaceLeft   = "namespace_left"
	EventNamespaceError  = "namespace_error"
	EventNamespaceList   = "namespace_list"
)

const (
	// System Events
	EventPing          = "ping"
	EventPong          = "pong"
	EventHeartbeat     = "heartbeat"
	EventServerMessage = "server_message"
	EventSystemAlert   = "system_alert"
	EventMaintenance   = "maintenance"
	EventServerRestart = "server_restart"
)

const (
	// Notification Events
	EventNotification = "notification"
	EventAlert        = "alert"
	EventUpdate       = "update"
	EventInfo         = "info"
	EventSuccess      = "success"
	EventWarning      = "warning"
)

// SocketClients manages WebSocket client connections across the entire application
type SocketClients struct {
	clients map[string]*ClientInfo
	mutex   sync.RWMutex
}

// ClientInfo holds detailed information about each connected client
type ClientInfo struct {
	UUID        string                 `json:"uuid"`
	UserID      string                 `json:"userId"`
	ConnectedAt time.Time              `json:"connectedAt"`
	LastSeen    time.Time              `json:"lastSeen"`
	Attributes  map[string]interface{} `json:"attributes"`
	Socket      *socketio.Websocket    `json:"-"`
}

// Global singleton instance
var (
	socketClientsInstance *SocketClients
	once                  sync.Once
)

// GetSocketClients returns the singleton instance of SocketClients
func GetSocketClients() *SocketClients {
	once.Do(func() {
		//log.Print("Calling Once DO ")
		socketClientsInstance = &SocketClients{
			clients: make(map[string]*ClientInfo),
		}
		go socketClientsInstance.CleanupStaleConnections(10 * time.Minute)
	})

	return socketClientsInstance
}

// AddClient adds a new client connection
func (s *SocketClients) AddClient(userID, uuid string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.clients[userID] = &ClientInfo{
		UUID:        uuid,
		UserID:      userID,
		ConnectedAt: time.Now(),
		LastSeen:    time.Now(),
		Attributes:  make(map[string]interface{}),
	}
}

// RemoveClient removes a client connection
func (s *SocketClients) RemoveClient(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.clients, userID)
}

func (s *SocketClients) SetSocket(userID string, socket *socketio.Websocket) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if client, exists := s.clients[userID]; exists {
		client.Socket = socket
		client.LastSeen = time.Now()
	}
}

// GetClient retrieves a specific client's information
func (s *SocketClients) GetClient(userID string) (*ClientInfo, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	client, exists := s.clients[userID]
	return client, exists
}

// GetClientUUID returns the UUID for a specific user
func (s *SocketClients) GetClientUUID(userID string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if client, exists := s.clients[userID]; exists {
		return client.UUID, true
	}
	return "", false
}

// GetClients returns a copy of all clients (thread-safe)
func (s *SocketClients) GetClients() map[string]*ClientInfo {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return a copy to prevent external modifications
	clientsCopy := make(map[string]*ClientInfo)
	for k, v := range s.clients {
		clientsCopy[k] = &ClientInfo{
			UUID:        v.UUID,
			UserID:      v.UserID,
			ConnectedAt: v.ConnectedAt,
			LastSeen:    v.LastSeen,
			Attributes:  make(map[string]interface{}),
			Socket:      v.Socket,
		}
		// Copy attributes
		for attrK, attrV := range v.Attributes {
			clientsCopy[k].Attributes[attrK] = attrV
		}
	}
	return clientsCopy
}

// GetClientsList returns a simple map of userID -> UUID for backward compatibility
func (s *SocketClients) GetClientsList() map[string]string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	clientsList := make(map[string]string)
	for userID, client := range s.clients {
		clientsList[userID] = client.UUID
	}
	return clientsList
}

// UpdateLastSeen updates the last seen timestamp for a client
func (s *SocketClients) UpdateLastSeen(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[userID]; exists {
		client.LastSeen = time.Now()
	}
}

// SetClientAttribute sets a custom attribute for a client
func (s *SocketClients) SetClientAttribute(userID, key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[userID]; exists {
		client.Attributes[key] = value
	}
}

// GetClientAttribute gets a custom attribute for a client
func (s *SocketClients) GetClientAttribute(userID, key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if client, exists := s.clients[userID]; exists {
		if value, attrExists := client.Attributes[key]; attrExists {
			return value, true
		}
	}
	return nil, false
}

// GetConnectedCount returns the number of connected clients
func (s *SocketClients) GetConnectedCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.clients)
}

// GetConnectedUsers returns a list of connected user IDs
func (s *SocketClients) GetConnectedUsers() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]string, 0, len(s.clients))
	for userID := range s.clients {
		users = append(users, userID)
	}
	return users
}

// IsUserConnected checks if a user is currently connected
func (s *SocketClients) IsUserConnected(userID string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.clients[userID]
	return exists
}

// GetClientsInRoom returns clients that have a specific room attribute
func (s *SocketClients) GetClientsInRoom(roomID string) map[string]*ClientInfo {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	roomClients := make(map[string]*ClientInfo)
	for userID, client := range s.clients {
		if room, exists := client.Attributes["room"]; exists && room == roomID {
			roomClients[userID] = client
		}
	}
	return roomClients
}

// CleanupStaleConnections removes connections that haven't been seen for a specified duration
func (s *SocketClients) CleanupStaleConnections(maxIdleTime time.Duration) []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var removedUsers []string
	cutoff := time.Now().Add(-maxIdleTime)

	for userID, client := range s.clients {
		if client.LastSeen.Before(cutoff) {
			delete(s.clients, userID)
			removedUsers = append(removedUsers, userID)
		}
	}

	return removedUsers
}
