package middleware

import (
	"errors"
	// "fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/degreane/octopus/config"
	"github.com/degreane/octopus/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
)

func CsrfFromHeader(headerName string) func(*fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Get(headerName)
		if token == "" {
			return "", fiber.ErrForbidden
		}
		return token, nil
	}
}

func NewCSRFMiddleware(store fiber.Storage) fiber.Handler {
	const HeaderName = "X-CSRF-Token"

	return func(c *fiber.Ctx) error {
		var isSecure bool
		log.Printf("Ise Secure %s", c.Protocol())
		if c.Protocol() == "http" {
			isSecure = false
		} else {
			isSecure = true
		}
		return csrf.New(csrf.Config{
			KeyLookup:         "header:" + HeaderName,
			CookieName:        "__Host-csrf_",
			CookieSameSite:    "Lax",
			CookieSecure:      isSecure,
			CookieSessionOnly: true,
			CookieHTTPOnly:    true,
			// Expiration:        1 * time.Hour,
			KeyGenerator: utils.UUIDv4,
			// Storage:           store,
			SessionKey:        "fiber.csrf.token",
			ContextKey:        "fiber.csrf.token",
			HandlerContextKey: "fiber.csrf.handler",
			// SingleUseToken:    true,
		})(c)
	}
}

type EoctoCsrfStore struct {
	*session.Store
	HandlerContextKey string
}

// var CsrfStore = session.New(session.Config{
// 	KeyLookup:         "cookie:__Host-csrf_",
// 	CookieSecure:      true,
// 	CookieHTTPOnly:    true,
// 	CookieSessionOnly: true,
// 	CookieSameSite:    "lax",

// })
var CsrfStore = NewCsrfStore()

func NewCsrfStore() *EoctoCsrfStore {
	appConfig, err := config.ParseServerConfig()
	if err != nil {
		//debug.Debug(debug.Error, fmt.Sprintf("<<MW>,sessions.go> Error parsing config file: %v", err))
		os.Exit(-1)
		return nil

	} else if appConfig.Storage == config.Redis {
		return &EoctoCsrfStore{
			Store: session.New(session.Config{
				Storage:           database.NewRedisStorage(),
				KeyLookup:         "cookie:__Host-csrf_",
				CookieSecure:      true,
				CookieHTTPOnly:    true,
				CookieSessionOnly: true,
				CookieSameSite:    "lax",
				Expiration:        15 * time.Minute,
			}),
			HandlerContextKey: "fiber.csrf.handler",
		}
	} else {
		return &EoctoCsrfStore{
			Store: session.New(session.Config{
				KeyLookup:         "cookie:__Host-csrf_",
				CookieSecure:      true,
				CookieHTTPOnly:    true,
				CookieSessionOnly: true,
				CookieSameSite:    "lax",
				Expiration:        15 * time.Minute,
			}),
			HandlerContextKey: "fiber.csrf.handler",
		}
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
func CreateCsrfSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := CsrfStore.Get(c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			return err
		}

		// Store current session data
		keys := sess.Keys()
		sessionData := make(map[string]interface{})
		for _, key := range keys {
			sessionData[key] = sess.Get(key)
		}

		// Destroy current session
		if err := sess.Destroy(); err != nil {
			log.Printf("Error destroying session: %v", err)
			return err
		}

		// Create new session with same data
		newSess, err := CsrfStore.Get(c)
		if err != nil {
			log.Printf("Error creating new session: %v", err)
			return err
		}

		// Restore data to new session
		for key, value := range sessionData {
			newSess.Set(key, value)
		}

		if err := newSess.Save(); err != nil {
			log.Printf("Error saving new session: %v", err)
			return err
		}

		log.Printf("Rotated CSRF session: %s", newSess.ID())
		return c.Next()
	}
}

func CreateEoctoCSRFMiddleware() fiber.Handler {
	const HeaderName = "X-CSRF-Token"
	return func(c *fiber.Ctx) error {
		sess, err := CsrfStore.Get(c)
		if err != nil {
			log.Printf("Error Getting CSRF  Session: %v", err)
		}
		keys := sess.Keys()
		sessionData := make(map[string]interface{})
		for _, key := range keys {
			sessionData[key] = sess.Get(key)
		}
		log.Printf("Keys From CRSF STore %+v", sessionData)
		if err := sess.Destroy(); err != nil {
			log.Printf("Error destroying CSRF session: %v", err)
			return err
		}
		log.Printf("Session CSRF ID = %s", sess.ID())
		sess.Set("csrf_token", utils.UUIDv4())
		log.Printf("CSRF Token: %s", sess.Get("csrf_token"))
		if err := sess.Save(); err != nil {
			log.Printf("Error saving CSRF session: %v", err)
			return err
		}
		sess, err = CsrfStore.Get(c)
		if err != nil {
			log.Printf("Error getting Second CSRF session: %v", err)
			return err
		}
		if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut || c.Method() == fiber.MethodDelete {
			// next we get the HostName
			hostName := c.Hostname()
			// also we get the referer
			referer := strings.ToLower(c.Get(fiber.HeaderReferer))
			if referer == "" {
				return errors.New("referer not supplied")
			}
			refererURL, err := url.Parse(referer)
			if err != nil {
				return errors.New("referer invalid")
			}
			if refererURL.Host != hostName {
				return errors.New("referer not allowed cross site")
			}
			log.Printf("Referer: %s", referer)
			log.Printf("HostName: %s", hostName)
			log.Printf("Referer HostName: %s", refererURL.Host)

		}
		if err := sess.Regenerate(); err != nil {
			log.Printf("Error regenerating CSRF session: %v", err)
			return err
		}

		log.Printf("Regenerated CSRF session:( %s )", sess.ID())
		log.Printf("CSRF Token: %s", sess.Get("csrf_token"))
		c.Locals("csrf_token", sess.Get("csrf_token"))
		return c.Next()

	}
}
