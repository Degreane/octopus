package utilities

import (
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

func GetCookie(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.ToString(1)
		cookie := c.Cookies(name)
		if cookie == "" {
			L.Push(lua.LNil)
			return 1
		}
		//log.Printf("Cookie value for key %s: %+v", name, cookie)
		L.Push(lua.LString(cookie))
		return 1
	}
}

func GetWsCookie(c *socketio.Websocket) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.ToString(1)
		cookie := c.Cookies(name)
		if cookie == "" {
			L.Push(lua.LNil)
			return 1
		}
		// log.Printf("Cookie value for key %s: %+v", name, cookie)
		L.Push(lua.LString(cookie))
		return 1
	}
}

// SetCookie creates a new cookie with the given name, value, and optional session configuration.
// If sessionOnly is false, the cookie will have an expiration time set based on sessionExpirey (in hours).
// The cookie is set to be secure and HTTP-only by default.
// Parameters:
//   - name: The name of the cookie
//   - value: The value of the cookie
//   - sessionOnly (optional): Whether the cookie is a session cookie (defaults to true)
//   - sessionExpirey (optional): Number of hours until the cookie expires (defaults to 744 hours / 31 days)
//   - sessionPath (optional): Path for the cookie (defaults to "/")
func SetCookie(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.ToString(1)
		value := L.ToString(2)
		sessionOnly := true
		sessionExpirey := 744
		sessionPath := "/"
		if L.GetTop() >= 3 {
			sessionOnly = L.ToBool(3)
			if L.GetTop() >= 4 {
				sessionExpirey = L.ToInt(4)
				if L.GetTop() >= 5 {
					sessionPath = L.ToString(5)
				}
			}
		}

		cookie := &fiber.Cookie{
			Name:        name,
			Value:       value,
			Secure:      true,
			HTTPOnly:    true,
			Path:        sessionPath,
			SessionOnly: sessionOnly,
		}

		if !sessionOnly {
			cookie.Expires = time.Now().Add(time.Duration(sessionExpirey) * time.Hour)
		}

		c.Cookie(cookie)
		return 0
	}
}

func SetWsCookie(c *socketio.Websocket) lua.LGFunction {
	// basically its here just as a placeholder to stay confined with the other functions
	return func(L *lua.LState) int {
		//name := L.ToString(1)
		//value := L.ToString(2)
		//sessionOnly := true
		//sessionExpirey := 744
		//sessionPath := "/"
		//if L.GetTop() >= 3 {
		//	sessionOnly = L.ToBool(3)
		//	if L.GetTop() >= 4 {
		//		sessionExpirey = L.ToInt(4)
		//		if L.GetTop() >= 5 {
		//			sessionPath = L.ToString(5)
		//		}
		//	}
		//}
		//
		//cookie := &fiber.Cookie{
		//	Name:        name,
		//	Value:       value,
		//	Secure:      true,
		//	HTTPOnly:    true,
		//	Path:        sessionPath,
		//	SessionOnly: sessionOnly,
		//}
		//
		//if !sessionOnly {
		//	cookie.Expires = time.Now().Add(time.Duration(sessionExpirey) * time.Hour)
		//}

		//c.Cookie(cookie)
		return 0
	}
}

func GetAllCookies(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		cookiesTable := L.NewTable()
		c.Request().Header.VisitAllCookie(func(key, value []byte) {
			cookiesTable.RawSetString(string(key), lua.LString(string(value)))
		})
		L.Push(cookiesTable)
		return 1
	}
}
func GetWsAllCookies(c *socketio.Websocket) lua.LGFunction {
	return func(L *lua.LState) int {
		//cookiesTable := L.NewTable()
		//c.Request().Header.VisitAllCookie(func(key, value []byte) {
		//	cookiesTable.RawSetString(string(key), lua.LString(string(value)))
		//})
		//L.Push(cookiesTable)
		return 0
	}
}

func DeleteCookie(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.ToString(1)
		c.Cookie(&fiber.Cookie{
			Name:     name,
			Value:    "",
			Expires:  time.Now().Add(-time.Hour * 24),
			HTTPOnly: true,
			Path:     "/",
		})
		return 0
	}
}

func ClearAllCookies(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		c.Request().Header.VisitAllCookie(func(key, _ []byte) {
			c.Cookie(&fiber.Cookie{
				Name:     string(key),
				Value:    "",
				Expires:  time.Now().Add(-time.Hour * 24),
				HTTPOnly: true,
				Path:     "/",
			})
		})
		return 0
	}
}

func DeleteWsCookie(c *socketio.Websocket) lua.LGFunction {
	return func(L *lua.LState) int {
		//name := L.ToString(1)
		//c.Cookie(&fiber.Cookie{
		//	Name:     name,
		//	Value:    "",
		//	Expires:  time.Now().Add(-time.Hour * 24),
		//	HTTPOnly: true,
		//	Path:     "/",
		//})
		return 0
	}
}

func ClearWsAllCookies(c *socketio.Websocket) lua.LGFunction {
	return func(L *lua.LState) int {
		//c.Request().Header.VisitAllCookie(func(key, _ []byte) {
		//	c.Cookie(&fiber.Cookie{
		//		Name:     string(key),
		//		Value:    "",
		//		Expires:  time.Now().Add(-time.Hour * 24),
		//		HTTPOnly: true,
		//		Path:     "/",
		//	})
		//})
		return 0
	}
}
