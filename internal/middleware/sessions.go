// Package middleware provides HTTP middleware functions for the application.
//
// This package contains middleware components that can be used to process
// HTTP requests and responses in a Fiber web application.
package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/degreane/octopus/config"
	"github.com/degreane/octopus/internal/database"
	"github.com/degreane/octopus/internal/utilities/debug"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory/v2"
	"github.com/gofiber/storage/redis/v3"
)

var Store = newSessionStore()

func GetAllSessions() []string {

	if store, ok := Store.Storage.(*redis.Storage); ok {
		//debug.Debug(debug.Info, "<<MW>, Redis sessions.go> GetAllSessions()")
		// For Redis storage
		byteKeys, err := store.Keys()
		if err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error getting session keys: %+v", err))
			return nil
		}
		keys := make([]string, len(byteKeys))
		for i, key := range byteKeys {
			// val, _ := store.Get(string(key))
			// debug.Debug(debug.Warning, fmt.Sprintf("Session key: %s => % +v", key, string(val)))
			keys[i] = string(key)
		}
		return keys
	}

	// For memory storage
	if store, ok := Store.Storage.(*memory.Storage); ok {
		//debug.Debug(debug.Info, "<<MW>, Memory sessions.go> GetAllSessions()")
		byteKeys, err := store.Keys()
		if err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error getting memory store keys: %+v", err))
			return nil
		}
		// Convert [][]byte to []string for memory store as well
		keys := make([]string, len(byteKeys))
		for i, byteKey := range byteKeys {
			keys[i] = string(byteKey)
		}
		return keys
	}

	return nil
}

func newSessionStore() *session.Store {
	//debug.Debug(debug.Warning, fmt.Sprintf("<<MW>,sessions.go> Calling New Session Store : %v", "123"))
	appConfig, err := config.ParseServerConfig()
	if err != nil {
		//debug.Debug(debug.Error, fmt.Sprintf("<<MW>,sessions.go> Error parsing config file: %v", err))
		os.Exit(-1)
		return nil

	} else if appConfig.Storage == config.Redis {
		//debug.Debug(debug.Warning, "<<MW>,sessions.go> Config file parsed successfully, Using Redis as storage")
		return session.New(session.Config{
			Storage:           database.NewRedisStorage(),
			Expiration:        24 * time.Hour,
			KeyLookup:         "cookie:session_id",
			CookieSecure:      false,
			CookieHTTPOnly:    true,
			CookieSessionOnly: true,
			CookieSameSite:    "lax",
		})
	} else {
		// for now we just use the memory store
		//debug.Debug(debug.Warning, "<<MW>,sessions.go> Config file parsed successfully, Using Memory as storage")
		return session.New(session.Config{
			KeyLookup:         "cookie:session_id",
			Expiration:        24 * time.Hour,
			CookieSecure:      false,
			CookieHTTPOnly:    true,
			CookieSessionOnly: true,
			CookieSameSite:    "lax",
		})
	}

}

// CreateSession returns a Fiber middleware handler that initializes and manages
// HTTP sessions. It performs the following operations:
//
// - Retrieves or creates a new session
//
// - Sets initial session data if empty
//
// - Logs session state for debugging
//
// - Ensures proper session persistence
//
// Usage:
//
//	app := fiber.New()
//
//	app.Use(middleware.CreateSession())
func CreateSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := Store.Get(c)
		if err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error getting session: %v", err))
			return err
		}

		// Store current session data
		keys := sess.Keys()
		// debug.Debug(debug.Error, fmt.Sprintf("Session Old = %+v", sess.Keys()))
		sessionData := make(map[string]interface{})
		for _, key := range keys {
			sessionData[key] = sess.Get(key)
		}

		// Destroy current session
		if err := sess.Destroy(); err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error destroying session: %v", err))
			// log.Printf("Error destroying session: %v", err)
			return err
		}
		sess.SetExpiry(2 * time.Second)

		// Create new session with same data
		newSess, err := Store.Get(c)
		if err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error creating new session: %v", err))
			// log.Printf("Error creating new session: %v", err)
			return err
		}

		// Restore data to new session
		for key, value := range sessionData {
			// debug.Debug(debug.Error, fmt.Sprintf("Key: %s, Value: %v", key, value))
			newSess.Set(key, value)
		}

		// debug.Debug(debug.Important, fmt.Sprintf("Rotating session: %s", newSess.ID()))
		if err := newSess.Save(); err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error saving new session: %v", err))
			// log.Printf("Error saving new session: %v", err)
			return err
		}
		c.Locals("c_session", newSess)
		// log.Printf("Rotated session: %s", newSess.ID())

		return c.Next()
	}
}
