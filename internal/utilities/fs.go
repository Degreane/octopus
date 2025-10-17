package utilities

import (
	"log"
	"os"
	"path/filepath"

	"github.com/degreane/octopus/internal/middleware"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

// getCWD returns the current working directory as a string. It returns an empty
// string if the retrieval fails.
func getCWD() string {
	dir, err := os.Getwd()
	if err == nil {
		return dir
	}
	return ""
}

// GetCWD exposes a Lua function that returns the current working directory for
// the session associated with the provided Fiber context. It caches the value
// in the session under the key "eocto_cWd". On failure to access the session,
// it returns nil to Lua.
func GetCWD(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		sess, err := middleware.Store.Get(c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			L.Push(lua.LNil)
			return 1
		}
		if val := sess.Get("eocto_cWd"); val != nil {
			if s, ok := val.(string); ok && s != "" {
				L.Push(lua.LString(s))
				return 1
			}
		}
		cwd := getCWD()
		if cwd != "" {
			sess.Set("eocto_cWd", cwd)
		}
		sess.Fresh()
		_ = sess.Save()
		L.Push(lua.LString(cwd))
		return 1
	}
}

// ResetWD exposes a Lua function that resets the working directory stored in
// the session to the current process working directory. It returns no values on
// success. On session access error, it returns nil.
func ResetWD(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		sess, err := middleware.Store.Get(c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			L.Push(lua.LNil)
			return 1
		}
		sess.Set("eocto_cWd", getCWD())
		sess.Fresh()
		_ = sess.Save()
		return 0
	}
}

// ListFiles exposes a Lua function that lists file and directory names in the
// working directory. The search path is determined by, in order of precedence:
// 1) the first Lua argument (string path), 2) the session value "eocto_cWd",
// 3) the current process working directory. On error it returns nil.
func ListFiles(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		var wd string
		sess, err := middleware.Store.Get(c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			L.Push(lua.LNil)
			return 1
		}
		if v := sess.Get("eocto_cWd"); v != nil {
			if s, ok := v.(string); ok {
				wd = s
			}
		}
		if wd == "" {
			wd = getCWD()
		}
		if L.GetTop() > 0 {
			wd = L.ToString(1)
		}
		files, err := os.ReadDir(wd)
		if err != nil {
			log.Printf("error listing files for path %q: %v", wd, err)
			L.Push(lua.LNil)
			return 1
		}
		table := L.NewTable()
		for _, file := range files {
			table.Append(lua.LString(file.Name()))
		}
		L.Push(table)
		return 1
	}
}

// SetWD exposes a Lua function that sets the working directory in the session
// to the provided path (first Lua argument). The path must exist and be a
// directory. Returns no values on success, or nil on error.
func SetWD(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		sess, err := middleware.Store.Get(c)
		if err != nil {
			log.Printf("Error getting session: %v", err)
			L.Push(lua.LNil)
			return 1
		}
		wd := L.ToString(1)
		if wd == "" {
			L.Push(lua.LNil)
			return 1
		}
		fInfo, err := os.Stat(wd)
		if err != nil {
			log.Printf("error changing dir to %q: %v", wd, err)
			L.Push(lua.LNil)
			return 1
		}
		if !fInfo.IsDir() {
			log.Printf("path is not a directory: %q", wd)
			L.Push(lua.LNil)
			return 1
		}

		if abs, err := filepath.Abs(wd); err == nil {
			sess.Set("eocto_cWd", abs)
			sess.Fresh()
			_ = sess.Save()
			log.Printf("Setting working directory to %s", abs)

		} else {
			sess.Set("eocto_cWd", abs)
			sess.Fresh()
			_ = sess.Save()
			log.Printf("Setting working directory to %s", wd)
		}
		return 0

	}
}
