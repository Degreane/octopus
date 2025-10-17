// Package routes provides routing functionality and request handling for the application.
//
// This package manages route configuration, Lua script integration, and dynamic
// route setup based on module configurations.
package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/degreane/octopus/config"
	"github.com/degreane/octopus/internal/middleware"
	"github.com/degreane/octopus/internal/utilities"
	"github.com/degreane/octopus/internal/utilities/debug"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	lua "github.com/yuin/gopher-lua"
)

// User represents application user data structure
type User struct {
	Name     string
	ID       string
	Username string
	Email    string
}

// RouteInfo stores metadata about registered routes
type RouteInfo struct {
	Method    string // HTTP method (GET, POST, etc.)
	Path      string // URL path
	Group     string // Group/module name for route organization
	View      string // Associated view template name
	WebSocket bool   // Indicates if the route is a WebSocket route
}

// MessageObject represents the envelope for messages exchanged over Socket.IO
type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

// per-connection storage for middlewares and context
type connBundle struct {
	ctx         *socketio.Websocket
	middlewares []func(*socketio.Websocket) error
}

var (
	connStore    sync.Map // key: mwKey(string) -> *connBundle
	registerOnce sync.Once
)

func WsScript(luaFile string, settings config.ModulesConfig, moduleBasePath ...string) func(*socketio.Websocket) error {
	return func(c *socketio.Websocket) error {
		luaState, ok := c.Locals("luaState").(*lua.LState)
		var L *lua.LState

		if !ok {
			L = lua.NewState()
			eoctoTable := luaState.NewTable()
			// debugging output
			eoctoTable.RawSetString("debug", L.NewFunction(utilities.Debug))
			// locals
			// L.SetGlobal("lua_getLocal", L.NewFunction(utilities.GetLocal(c)))
			eoctoTable.RawSetString("getLocal", L.NewFunction(utilities.GetWsLocal(c)))
			// L.SetGlobal("lua_setLocal", L.NewFunction(utilities.SetLocal(c)))
			eoctoTable.RawSetString("setLocal", L.NewFunction(utilities.SetWsLocal(c)))
			// L.SetGlobal("lua_deleteLocal", L.NewFunction(utilities.DeleteLocal(c)))
			eoctoTable.RawSetString("deleteLocal", L.NewFunction(utilities.DeleteWsLocal(c)))
			// L.SetGlobal("lua_getLocals", L.NewFunction(utilities.GetLocals(c)))
			eoctoTable.RawSetString("getLocals", L.NewFunction(utilities.GetWsLocals(c)))

			eoctoTable.RawSetString("getCookie", L.NewFunction(utilities.GetWsCookie(c)))
			eoctoTable.RawSetString("setCookie", L.NewFunction(utilities.SetWsCookie(c)))
			eoctoTable.RawSetString("getAllCookies", L.NewFunction(utilities.GetWsAllCookies(c)))
			eoctoTable.RawSetString("deleteCookie", L.NewFunction(utilities.DeleteWsCookie(c)))
			eoctoTable.RawSetString("clearAllCookies", L.NewFunction(utilities.ClearWsAllCookies(c)))

			eoctoTable.RawSetString("getUUID", L.NewFunction(func(L *lua.LState) int {
				uuid := utils.UUIDv4()
				L.Push(lua.LString(uuid))
				return 1
			}))
			L.SetGlobal("eocto", eoctoTable)
			c.Conn.Locals("luaState", L)
		} else {
			L = luaState
		}
		scriptPath := luaFile
		if len(moduleBasePath) > 0 && moduleBasePath[0] != "" {
			scriptPath = filepath.Join(moduleBasePath[0], luaFile)
		} else {
			scriptPath = filepath.Join("modules", luaFile)
		}
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			debug.Debug(debug.Error, fmt.Sprintf("Script file %s does not exist ", scriptPath))

			return err
		}
		debug.Debug(debug.Info, fmt.Sprintf("script file %s,\n\t ScriptFilePath %s", luaFile, scriptPath))
		if err := L.DoFile(scriptPath); err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error executing Lua script: %v", err))
			return err
		}
		if L.GetTop() > 0 {
			returnValue := L.Get(-1).String()
			L.Pop(1)

			c.Conn.Locals(scriptPath, returnValue)
		}
		debug.Debug(debug.Info, fmt.Sprintf("Script %s executed successfully and c.Locals contain %+v\n", scriptPath, c.Locals(scriptPath)))
		return nil
	}
}

// Script returns a Fiber middleware handler that executes Lua scripts.
// It maintains Lua state and provides session access to Lua scripts.
//
// Parameters:
//   - luaFile: path to the Lua script file
//   - moduleBasePath: optional base path for module-specific scripts
//
// The handler sets up a Lua environment with:
//   - Session management functions (getSession, setSession)
//   - Script execution results in context locals
func Script(luaFile string, settings config.ModulesConfig, moduleBasePath ...string) fiber.Handler {
	//debug.Debug(debug.Info, fmt.Sprintf("Script for lua file %s", luaFile))
	return func(c *fiber.Ctx) error {
		luaState, ok := c.Locals("luaState").(*lua.LState)
		var L *lua.LState

		if !ok {
			// log.Printf("Lua state not found in context locals Script Handler")
			L = lua.NewState()
			eoctoTable := luaState.NewTable()
			// debug Messages
			eoctoTable.RawSetString("debug", L.NewFunction(utilities.Debug))

			// sessions
			// L.SetGlobal("lua_getSession", L.NewFunction(utilities.GetSession(c)))
			eoctoTable.RawSetString("getSession", L.NewFunction(utilities.GetSession(c)))
			// L.SetGlobal("lua_setSession", L.NewFunction(utilities.SetSession(c)))
			eoctoTable.RawSetString("setSession", L.NewFunction(utilities.SetSession(c)))
			// L.SetGlobal("lua_deleteSession", L.NewFunction(utilities.DeleteSession(c)))
			eoctoTable.RawSetString("deleteSession", L.NewFunction(utilities.DeleteSession(c)))
			eoctoTable.RawSetString("setSessionExpiry", L.NewFunction(utilities.SetSessionExpiry(c)))
			// Cookie Handlers
			eoctoTable.RawSetString("getCookie", L.NewFunction(utilities.GetCookie(c)))
			eoctoTable.RawSetString("setCookie", L.NewFunction(utilities.SetCookie(c)))
			eoctoTable.RawSetString("getAllCookies", L.NewFunction(utilities.GetAllCookies(c)))
			eoctoTable.RawSetString("deleteCookie", L.NewFunction(utilities.DeleteCookie(c)))
			eoctoTable.RawSetString("clearAllCookies", L.NewFunction(utilities.ClearAllCookies(c)))

			// webSockets
			eoctoTable.RawSetString("wsAddRoom", L.NewFunction(utilities.WsAddRoom(c)))
			eoctoTable.RawSetString("wsRemoveRoom", L.NewFunction(utilities.WsRemoveRoom(c)))
			eoctoTable.RawSetString("wsGetUserRooms", L.NewFunction(utilities.WsGetUserRooms(c)))
			eoctoTable.RawSetString("wsIsUserInRoom", L.NewFunction(utilities.WsIsUserInRoom(c)))
			eoctoTable.RawSetString("wsEmitToRoom", L.NewFunction(utilities.WsEmitToRoom(c)))
			//eoctoTable.RawSetString()

			// csrf
			eoctoTable.RawSetString("getCsrfToken", L.NewFunction(utilities.GetCsrfToken(c)))
			// locals
			// L.SetGlobal("lua_getLocal", L.NewFunction(utilities.GetLocal(c)))
			eoctoTable.RawSetString("getLocal", L.NewFunction(utilities.GetLocal(c)))
			// L.SetGlobal("lua_setLocal", L.NewFunction(utilities.SetLocal(c)))
			eoctoTable.RawSetString("setLocal", L.NewFunction(utilities.SetLocal(c)))
			// L.SetGlobal("lua_deleteLocal", L.NewFunction(utilities.DeleteLocal(c)))
			eoctoTable.RawSetString("deleteLocal", L.NewFunction(utilities.DeleteLocal(c)))
			// L.SetGlobal("lua_getLocals", L.NewFunction(utilities.GetLocals(c)))
			eoctoTable.RawSetString("getLocals", L.NewFunction(utilities.GetLocals(c)))
			// headers
			// L.SetGlobal("lua_getHeaders", L.NewFunction(utilities.GetHeaders(c)))
			eoctoTable.RawSetString("getHeaders", L.NewFunction(utilities.GetHeaders(c)))
			// L.SetGlobal("lua_getHeader", L.NewFunction(utilities.GetHeader(c)))
			eoctoTable.RawSetString("getHeader", L.NewFunction(utilities.GetHeader(c)))
			// L.SetGlobal("lua_setHeader", L.NewFunction(utilities.SetHeader(c)))
			eoctoTable.RawSetString("setHeader", L.NewFunction(utilities.SetHeader(c)))
			// L.SetGlobal("lua_deleteHeader", L.NewFunction(utilities.DeleteHeader(c)))
			eoctoTable.RawSetString("deleteHeader", L.NewFunction(utilities.DeleteHeader(c)))
			// L.SetGlobal("lua_getPath", L.NewFunction(utilities.GetPath(c)))
			eoctoTable.RawSetString("getPath", L.NewFunction(utilities.GetPath(c)))
			// L.SetGlobal("lua_getHost", L.NewFunction(utilities.GetHost(c)))
			eoctoTable.RawSetString("getHost", L.NewFunction(utilities.GetHost(c)))
			// L.SetGlobal("lua_getSchema", L.NewFunction(utilities.GetSchema(c)))
			eoctoTable.RawSetString("getSchema", L.NewFunction(utilities.GetSchema(c)))
			// query parameters
			// L.SetGlobal("lua_getQueryParams", L.NewFunction(utilities.GetQueryParams(c)))
			eoctoTable.RawSetString("getQueryParams", L.NewFunction(utilities.GetQueryParams(c)))

			// path parameters
			eoctoTable.RawSetString("getPathParams", L.NewFunction(utilities.GetPathParams(c)))
			// get specific path param
			eoctoTable.RawSetString("getPathParam", L.NewFunction(utilities.GetPathParam(c)))
			// posted data
			// L.SetGlobal("lua_getPostBody", L.NewFunction(utilities.GetPostBody(c)))
			eoctoTable.RawSetString("getPostBody", L.NewFunction(utilities.GetPostBody(c)))
			// get the method
			// L.SetGlobal("lua_getMethod", L.NewFunction(utilities.GetMethod(c)))
			eoctoTable.RawSetString("getMethod", L.NewFunction(utilities.GetMethod(c)))
			// encryption/decryption
			// L.SetGlobal("lua_decryptData", L.NewFunction(utilities.GetDecryptData(c)))
			eoctoTable.RawSetString("decryptData", L.NewFunction(utilities.GetDecryptData(c)))
			// L.SetGlobal("lua_encryptData", L.NewFunction(utilities.GetEncryptData(c)))
			eoctoTable.RawSetString("encryptData", L.NewFunction(utilities.GetEncryptData(c)))
			// L.SetGlobal("lua_decodeJSON", L.NewFunction(utilities.GetDecodeJSON(c)))
			eoctoTable.RawSetString("decodeJSON", L.NewFunction(utilities.GetDecodeJSON(c)))
			eoctoTable.RawSetString("encodeJSON", L.NewFunction(utilities.GetEncodeJSON(c)))
			// base 32
			eoctoTable.RawSetString("encodeBase32", L.NewFunction(utilities.GetEncodeBase32(c)))
			eoctoTable.RawSetString("decodeBase32", L.NewFunction(utilities.GetDecodeBase32(c)))

			// mongodbDatabase functionalities
			eoctoTable.RawSetString("getDataFromCollection", L.NewFunction(utilities.GetDataFromCollectionLua))
			eoctoTable.RawSetString("setDataToCollection", L.NewFunction(utilities.SetDataToCollectionLua))
			eoctoTable.RawSetString("delDataFromCollection", L.NewFunction(utilities.DelDataFromCollectionLua))
			eoctoTable.RawSetString("insertDataToCollection", L.NewFunction(utilities.InsertDataToCollectionLua))

			// make http Requests to other servers
			// L.SetGlobal("lua_makeRequest", L.NewFunction(utilities.GetRequest(c)))
			eoctoTable.RawSetString("makeRequest", L.NewFunction(utilities.GetRequest(c)))
			eoctoTable.RawSetString("proxy", L.NewFunction(utilities.ProxyRequestLua(c)))
			// Register Twilio account
			// Register the WhatsApp function
			eoctoTable.RawSetString("sendWhatsAppMessage", L.NewFunction(utilities.SendWhatsAppMessageLua))

			// set lua response
			// L.SetGlobal("lua_setResponse", L.NewFunction(utilities.GetResponse(c)))
			eoctoTable.RawSetString("setResponse", L.NewFunction(utilities.GetResponse(c)))
			// expose c.Render to lua
			eoctoTable.RawSetString("render", L.NewFunction(utilities.GetRender(c)))
			eoctoTable.RawSetString("renderJson", L.NewFunction(utilities.GetRenderJson(c)))

			// Expose getUUID function
			eoctoTable.RawSetString("getUUID", L.NewFunction(func(L *lua.LState) int {
				uuid := utils.UUIDv4()
				L.Push(lua.LString(uuid))
				return 1
			}))
			eoctoTable.RawSetString("getSettings", L.NewFunction(func(l *lua.LState) int {
				tbl := L.NewTable()

				tbl.RawSetString("BasePath", lua.LString(settings.BasePath))
				tbl.RawSetString("LocalPath", lua.LString(settings.LocalPath))
				L.Push(tbl)
				return 1
			}))
			// Redis functionalities
			eoctoTable.RawSetString("getRedis", L.NewFunction(utilities.GetRedisValueLua))
			eoctoTable.RawSetString("setRedis", L.NewFunction(utilities.SetRedisValueLua))
			eoctoTable.RawSetString("deleteRedis", L.NewFunction(utilities.DeleteRedisKeyLua))
			eoctoTable.RawSetString("timeStampNano", L.NewFunction(func(l *lua.LState) int {
				tstamp := time.Now().UnixNano()
				L.Push(lua.LNumber(tstamp))
				return 1
			}))
			eoctoTable.RawSetString("timeStampMilli", L.NewFunction(func(l *lua.LState) int {
				tstamp := time.Now().UnixMilli()
				L.Push(lua.LNumber(tstamp))
				return 1
			}))
			eoctoTable.RawSetString("timeStamp", L.NewFunction(func(l *lua.LState) int {
				tstamp := time.Now().Unix()
				L.Push(lua.LNumber(tstamp))
				return 1
			}))
			// FileSystem functionalities
			eoctoTable.RawSetString("getCWD", L.NewFunction(utilities.GetCWD(c)))
			eoctoTable.RawSetString("resetWD", L.NewFunction(utilities.ResetWD(c)))
			eoctoTable.RawSetString("setWD", L.NewFunction(utilities.SetWD(c)))
			eoctoTable.RawSetString("listFiles", L.NewFunction(utilities.ListFiles(c)))
			// YAML utilities
			eoctoTable.RawSetString("readYamlFile", L.NewFunction(utilities.ReadYamlFileLua))
			// Set the eocto table a a global
			L.SetGlobal("eocto", eoctoTable)
			c.Locals("luaState", L)
		} else {
			L = luaState
		}
		scriptPath := luaFile
		if len(moduleBasePath) > 0 && moduleBasePath[0] != "" {
			scriptPath = filepath.Join(moduleBasePath[0], luaFile)
		} else {
			scriptPath = filepath.Join("modules", luaFile)
		}
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			debug.Debug(debug.Error, fmt.Sprintf("Script file %s does not exist ", scriptPath))

			return c.Next()
		}
		debug.Debug(debug.Info, fmt.Sprintf("script file %s,\n\t ScriptFilePath %s", luaFile, scriptPath))
		if err := L.DoFile(scriptPath); err != nil {
			debug.Debug(debug.Error, fmt.Sprintf("Error executing Lua script: %v", err))
			return c.Next()
		}
		if L.GetTop() > 0 {
			returnValue := L.Get(-1).String()
			L.Pop(1)

			c.Locals(scriptPath, returnValue)
		}
		debug.Debug(debug.Info, fmt.Sprintf("Script %s executed successfully and c.Locals contain %+v\n", scriptPath, c.Locals(scriptPath)))
		return c.Next()
	}
}

func CreateSocketIOWIthMessageMiddlewares(c *fiber.Ctx, middlewares ...func(*socketio.Websocket) error) fiber.Handler {
	registerOnce.Do(registerGlobalHandlers)
	//registerGlobalHandlers()
	debug.Debug(debug.Warning, fmt.Sprintf("CreateSocketIOWIthMessageMiddlewares  % +v", c))
	return socketio.New(func(kws *socketio.Websocket) {
		// Get the original Fiber context
		//raw := kws.Locals("fiber.ctx")
		//ctx, _ := raw.(*fiber.Ctx)
		//if ctx == nil {
		//	debug.Debug(debug.Error, "socketio_connection: ctx is nil")
		//	return
		//}

		// Generate a unique key for this connection
		mwKey := fmt.Sprintf("mw-%p", kws)

		// Stash the key on the connection so global handlers can retrieve its bundle
		debug.Debug(debug.Info, "mw_key: %s", mwKey)
		kws.SetAttribute("mw_key", mwKey)
		kws.SetAttribute(fmt.Sprintf("http_connection-%p", kws), &connBundle{
			ctx:         kws,
			middlewares: middlewares,
		})

		// Store the bundle for this connection
		connStore.Store(mwKey, &connBundle{
			ctx:         kws,
			middlewares: middlewares,
		})

		// Optionally expose the connection on the context
		//if ctx != nil {
		//	c.Locals("socketio_connection", kws)
		//}

	})
}

// registerGlobalHandlers sets up the global Socket.IO handlers for message and cleanup
func registerGlobalHandlers() {
	// Handle all incoming messages globally
	socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
		debug.Debug(debug.Warning, fmt.Sprintf("socketio_message received, data: % +v", ep.Kws.Conn))
		debug.Debug(debug.Info, "socketio_message received")
		mwKey := ep.Kws.GetStringAttribute("mw_key")
		if mwKey == "" {
			debug.Debug(debug.Error, "mw_key not set, skipping message")
			// Not a connection created by our handler
			return
		}
		debug.Debug(debug.Info, "Recieved mw_key: %s", mwKey)
		// Lookup the bundle for this connection

		//value, ok := connStore.Load(mwKey)
		//if !ok {
		//	debug.Debug(debug.Error, "mw_key not found in connStore, skipping message")
		//	return
		//}
		//debug.Debug(debug.Info, fmt.Sprintf("mw_key found in connStore % +v", value))
		bundlee := ep.Kws.GetAttribute(fmt.Sprintf("http_connection-%p", ep.Kws))
		if bundlee == "" {
			debug.Debug(debug.Error, "mw_key not found in connStore, skipping message")
			return
		}
		bundle := bundlee.(*connBundle)
		//bundle := value.(*connBundle)

		if bundle.ctx == nil {
			// If no ctx, just echo
			debug.Debug(debug.Info, "socketio_message received, but no ctx, echoing message")
			_ = ep.Kws.Conn.WriteMessage(socketio.TextMessage, ep.Data)
			return
		}

		// Parse message into MessageObject (fallback to plain text)
		var msgObj MessageObject
		if err := json.Unmarshal(ep.Data, &msgObj); err != nil {
			msgObj = MessageObject{
				Data:  string(ep.Data),
				Event: "message",
			}
		}
		//_ctx := bundle.ctx
		// Inject message values into ctx locals for middlewares
		ep.Kws.Conn.Locals("socketio_message_raw", ep.Data)
		ep.Kws.Conn.Locals("socketio_message_data", string(ep.Data))
		ep.Kws.Conn.Locals("socketio_message_from", msgObj.From)
		ep.Kws.Conn.Locals("socketio_message_event", msgObj.Event)
		ep.Kws.Conn.Locals("socketio_message_to", msgObj.To)
		ep.Kws.Conn.Locals("socketio_message_object", msgObj)

		// Run middlewares in order
		var mwErr error
		for _, mw := range bundle.middlewares {
			if err := mw(bundle.ctx); err != nil {
				mwErr = err
				break
			}
		}

		if mwErr != nil {
			// Reply with middleware error payload
			_ = sendSocketIOJSON(ep.Kws, fiber.Map{
				"type":    "middleware_error",
				"error":   mwErr.Error(),
				"message": msgObj,
			})
			cleanupMessageLocals(bundle.ctx)
			return
		}

		// Check for custom response set by middlewares
		if resp := bundle.ctx.Locals("socketio_response"); resp != nil {
			switch r := resp.(type) {
			case fiber.Map, map[string]interface{}:
				_ = sendSocketIOJSON(ep.Kws, r)
			case string:
				_ = ep.Kws.Conn.WriteMessage(socketio.TextMessage, []byte(r))
			case []byte:
				_ = ep.Kws.Conn.WriteMessage(socketio.TextMessage, r)
			default:
				_ = sendSocketIOJSON(ep.Kws, r)
			}
			bundle.ctx.Conn.Locals("socketio_response", nil)
		} else {
			// Default echo behavior
			debug.Debug(debug.Info, "socketio_response not set, echoing message")
			_ = ep.Kws.Conn.WriteMessage(socketio.TextMessage, ep.Data)
		}

		// Cleanup per-message locals
		cleanupMessageLocals(bundle.ctx)
	})

	// Cleanup when client disconnects
	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		mwKey := ep.Kws.GetStringAttribute("mw_key")
		if mwKey != "" {
			connStore.Delete(mwKey)
		}
	})

	// Ensure cleanup on close as well
	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		mwKey := ep.Kws.GetStringAttribute("mw_key")
		if mwKey != "" {
			connStore.Delete(mwKey)
		}
	})
}

// sendSocketIOJSON marshals v to JSON and sends it as a text message
func sendSocketIOJSON(kws *socketio.Websocket, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		debug.Debug(debug.Error, "sendSocketIOJSON: json.Marshal failed: %v", err)
		return err
	}
	return kws.Conn.WriteMessage(socketio.TextMessage, b)
}

// cleanupMessageLocals clears the per-message locals
func cleanupMessageLocals(c *socketio.Websocket) {
	if c == nil {
		debug.Debug(debug.Error, "cleanupMessageLocals: c is nil")
		return
	}
	c.Conn.Locals("socketio_message_raw", nil)
	c.Conn.Locals("socketio_message_data", nil)
	c.Conn.Locals("socketio_message_from", nil)
	c.Conn.Locals("socketio_message_event", nil)
	c.Conn.Locals("socketio_message_to", nil)
	c.Conn.Locals("socketio_message_object", nil)
}

// SetupRoutes configures application routes based on module configuration.
// It handles:
//   - Static file serving
//   - Route grouping
//   - Middleware chain setup
//   - View rendering
//   - Pre-check script execution
//
// Parameters:
//   - app: Fiber application instance
//   - module: Module configuration containing route definitions
//
// Returns error if route setup fails
func SetupRoutes(app *fiber.App, module config.ModulesConfig) error {
	// log.Printf("Setting up routes for module %+v", module)
	//logr := lgr.GetLogger().WithField("component", "routes")
	// The key for the map is message.to
	//clients := utilities.GetSocketClients()
	if module.Name == "" && module.BasePath == "" {
		return nil
	}
	var routes []RouteInfo
	group := app.Group(module.BasePath)

	// group.Use(createSession())
	/*
		if LocalPath Exsits then we are going to use LocalPath for our modules

	*/

	var staticPath string
	var cssPath string
	var jsPath string
	var imgPath string
	var fontsPath string
	var iconsPath string
	if strings.TrimSpace(module.LocalPath) != "" && strings.TrimSpace(module.LocalPath) != "/" {
		staticPath = filepath.Clean(fmt.Sprintf("views/%s/public", module.LocalPath))
		cssPath = filepath.Clean(fmt.Sprintf("views/%s/public/css", module.LocalPath))
		jsPath = filepath.Clean(fmt.Sprintf("views/%s/public/js", module.LocalPath))
		imgPath = filepath.Clean(fmt.Sprintf("views/%s/public/images", module.LocalPath))
		fontsPath = filepath.Clean(fmt.Sprintf("views/%s/public/fonts", module.LocalPath))
		iconsPath = filepath.Clean(fmt.Sprintf("views/%s/public/icons", module.LocalPath))
	} else if strings.TrimSpace(module.BasePath) != "" && strings.TrimSpace(module.BasePath) != "/" {
		staticPath = filepath.Clean(fmt.Sprintf("views/%s/public", module.BasePath))
		cssPath = filepath.Clean(fmt.Sprintf("views/%s/public/css", module.BasePath))
		jsPath = filepath.Clean(fmt.Sprintf("views/%s/public/js", module.BasePath))
		imgPath = filepath.Clean(fmt.Sprintf("views/%s/public/images", module.BasePath))
		fontsPath = filepath.Clean(fmt.Sprintf("views/%s/public/fonts", module.BasePath))
		iconsPath = filepath.Clean(fmt.Sprintf("views/%s/public/icons", module.BasePath))
	} else if module.LocalPath == "/" {
		staticPath = filepath.Clean("views/public")
		cssPath = filepath.Clean("views/public/css")
		jsPath = filepath.Clean("views/public/js")
		imgPath = filepath.Clean("views/public/images")
		fontsPath = filepath.Clean("views/public/fonts")
		iconsPath = filepath.Clean("views/public/icons")
	} else if module.BasePath == "/" {
		staticPath = filepath.Clean("views/public")
		cssPath = filepath.Clean("views/public/css")
		jsPath = filepath.Clean("views/public/js")
		imgPath = filepath.Clean("views/public/images")
		fontsPath = filepath.Clean("views/public/fonts")
		iconsPath = filepath.Clean("views/public/icons")
	}

	group.Static("/static", staticPath)
	group.Static("/css", cssPath)
	group.Static("/js", jsPath)
	group.Static("/img", imgPath)
	group.Static("/images", imgPath)
	group.Static("/fonts", fontsPath)
	group.Static("/icons", iconsPath)

	for _, route := range module.Routes {
		var middlewares []fiber.Handler
		var wsmiddlewares []func(*socketio.Websocket) error

		// middlewares = append(middlewares, middleware.CreateSession())
		// middlewares = append(middlewares, middleware.CreateEoctoCSRFMiddleware())
		if len(route.PreCheck) > 0 {
			// log.Printf("Pre-check for route %s", route.Path)
			//if route.WebSocket {
			//	middlewares = append(middlewares, websocket.New(func(c *websocket.Conn) {
			//		for {
			//			mt, msg, err := c.ReadMessage()
			//			if err != nil {
			//				logr.Error(fmt.Sprintf("MiddleWare read: %+v", err))
			//				//log.Println("MiddleWare read:", err)
			//				break
			//			}
			//			logr.Info(fmt.Sprintf("MiddleWare recv: %s %d", msg, mt))
			//			err = c.WriteMessage(mt, msg)
			//			if err != nil {
			//				logr.Error(fmt.Sprintf("MiddleWare write: %+v", err))
			//				break
			//			}
			//		}
			//
			//	}))
			//}
			for _, check := range route.PreCheck {
				if check.Script != "" {
					if route.WebSocket {
						//debug.Debug(debug.Important, fmt.Sprintf("Script for route %s=> %s", route.Path, check.Script))
						if strings.TrimSpace(module.LocalPath) != "" {
							wsmiddlewares = append(wsmiddlewares, WsScript(check.Script, module, filepath.Join("views", module.LocalPath, "scripts")))
						} else {
							wsmiddlewares = append(wsmiddlewares, WsScript(check.Script, module, filepath.Join("views", module.BasePath, "scripts")))
						}
					} else {
						//debug.Debug(debug.Important, fmt.Sprintf("Script for route %s=> %s", route.Path, check.Script))
						if strings.TrimSpace(module.LocalPath) != "" {
							middlewares = append(middlewares, Script(check.Script, module, filepath.Join("views", module.LocalPath, "scripts")))
						} else {
							middlewares = append(middlewares, Script(check.Script, module, filepath.Join("views", module.BasePath, "scripts")))
						}
					}

				}

			}
		}
		routeInfo := RouteInfo{
			Method:    route.Method,
			Path:      route.Path,
			Group:     module.BasePath, // Store the base path as the group name
			View:      route.View,      // Set the default view to an empty string
			WebSocket: route.WebSocket,
		}
		routes = append(routes, routeInfo)
		// middleware.NewCSRFMiddleware(store)
		// , middleware.CreateSession()
		if route.WebSocket {
			group.Add(route.Method, route.Path, middleware.CreateSession(), func(c *fiber.Ctx) error {
				return CreateSocketIOWIthMessageMiddlewares(c, wsmiddlewares...)(c)
			})

		} else {
			group.Add(route.Method, route.Path, append(middlewares, middleware.CreateSession(), func(c *fiber.Ctx) error {
				// l1 := c.Locals("l1").(string)
				if luaResp := c.Locals("lua_response"); luaResp != nil {
					if respMap, ok := luaResp.(map[int]fiber.Map); ok {
						for status, jsonMap := range respMap {
							return c.Status(status).JSON(jsonMap)
						}
					}
				}
				if c.Locals("rendered_from_lua") != nil {
					return nil
				}
				log.Println("__ CSRF TOKEN")
				_tkn := c.Locals("csrf_token")
				var tkn string = ""
				if _tkn != nil {
					tkn = _tkn.(string)
				}
				log.Println(tkn)
				log.Println("^^ CSRF TOKEN")
				//debug.Debug(debug.Warning, fmt.Sprintf("GetAllSessions: % +v", middleware.GetAllSessions()))

				basePath := module.BasePath
				if strings.TrimSpace(module.LocalPath) != "" && strings.TrimSpace(module.LocalPath) != "/" {
					basePath = strings.TrimSpace(module.LocalPath)
				} else if strings.TrimSpace(module.LocalPath) != "" && strings.TrimSpace(module.LocalPath) == "/" {
					basePath = ""
				} else if strings.TrimSpace(module.BasePath) != "" && strings.TrimSpace(module.BasePath) == "/" {
					basePath = ""
				}

				thePath := fmt.Sprintf("%s/%s", strings.TrimPrefix(basePath, "/"), route.View)
				thePath = fmt.Sprintf("%s", strings.TrimPrefix(thePath, "/"))
				debug.Debug(debug.Warning, fmt.Sprintf("%s", thePath))
				return c.Render(thePath, fiber.Map{
					"basePath":  strings.TrimSpace(module.BasePath),
					"localPath": strings.TrimSpace(module.LocalPath),
					"csrf":      tkn,
				})
			})...)
		}

	}
	//for _, r := range routes {
	//	debug.Debug(debug.Important, fmt.Sprintf("Route: Method=%s, Path=%s, Group=%s, WebSocket=%v", r.Method, r.Path, r.Group, r.WebSocket))
	//}
	return nil
}
