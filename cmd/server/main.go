// Package main is the entry point for the Octopus server application.
// It initializes the server, sets up middleware, and starts listening for requests.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/degreane/octopus/config"
	"github.com/degreane/octopus/internal/routes"
	lgr "github.com/degreane/octopus/internal/service/logger"
	"github.com/degreane/octopus/internal/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

// MessageObject Basic chat message object
type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

var clients *utilities.SocketClients

// main is the entry point for the Octopus server application.
// It initializes logging, loads configuration, sets up a Fiber web server with various middleware,
// configures HTML template engine with custom template functions, and starts the server listening on a specified port.
// The server supports features like compression, health checks, metrics monitoring, and static file serving.
func main() {

	// Note: rand.Seed() is deprecated in Go 1.20+ and no longer needed
	// The global random number generator is automatically seeded
	// Initialize the socket client map
	clients = utilities.GetSocketClients()
	// Initialize the logger with the "Server" component tag for easier log filtering and identification
	logr := lgr.GetLogger().WithField("component", "Server")
	//logr.Info("Starting Octopus Server")

	// Load environment variables from .env file
	// This allows for configuration via environment variables without modifying code
	err := godotenv.Load()
	if err != nil {
		// Log error but continue execution as .env file might be optional
		logr.Error("Error loading .env file", err)
	}

	// Initialize and parse server configuration from environment variables or config files
	// This centralizes all configuration management in the config package
	appConfig, err := config.ParseServerConfig()
	if err != nil {
		// Fatal error stops execution as the server cannot run without proper configuration
		logr.Fatal(fmt.Sprintf("Error initializing config %+v", err))
	}

	// Initialize the HTML template engine with the views directory and .html extension
	// This engine will be used to render HTML templates for web pages
	engine := config.SetupTemplateEngine("./views", ".html", true)
	engine.Reload(true)
	engine.AddFunc("dict", dictHelper)
	engine.AddFunc("checkType", checkType)
	engine.AddFunc("length", length)
	engine.AddFunc("formatNum", formatNum)
	engine.AddFunc("FormatNumWithComma", formatNumWithComma)
	engine.AddFunc("iterate", iterate)
	engine.AddFunc("iterateFilter", iterateFilter)
	engine.AddFunc("iterateRange", iterateRange)
	engine.AddFunc("iterateEven", iterateEven)
	engine.AddFunc("iterateOdd", iterateOdd)
	engine.AddFunc("iterateMultiple", iterateMultiple)
	engine.AddFunc("rand", randHelper)
	engine.AddFunc("randRange", randRange)
	engine.AddFunc("randFloat", randFloat)
	engine.AddFunc("randFloatRange", randFloatRange)
	engine.AddFunc("intEq", intEq)
	engine.AddFunc("intGt", intGt)
	engine.AddFunc("intGte", intGte)
	engine.AddFunc("intLt", intLt)
	engine.AddFunc("intLte", intLte)
	engine.AddFunc("intNe", intNe)
	engine.AddFunc("addInt", addInt)
	engine.AddFunc("subtractInt", subtractInt)
	engine.AddFunc("multiplyInt", multiplyInt)
	engine.AddFunc("divideInt", divideInt)
	engine.AddFunc("modInt", modInt)
	engine.AddFunc("absInt", absInt)
	engine.AddFunc("maxInt", maxInt)
	engine.AddFunc("minInt", minInt)
	engine.AddFunc("timestamp", timestamp)
	engine.AddFunc("dict", dict)

	engine.AddFunc("multiply", multiply)
	engine.AddFunc("add", addFloat)           // Add this line
	engine.AddFunc("subtract", subtractFloat) // Add this line
	engine.AddFunc("divide", divideFloat)     // Add this line
	engine.AddFunc("floatEq", floatEq)
	engine.AddFunc("floatGt", floatGt)
	engine.AddFunc("floatGte", floatGte)
	engine.AddFunc("floatLt", floatLt)
	engine.AddFunc("floatLte", floatLte)
	engine.AddFunc("floatNe", floatNe)
	engine.AddFunc("addFloat", addFloat)
	engine.AddFunc("subtractFloat", subtractFloat)
	engine.AddFunc("multiplyFloat", multiplyFloat)
	engine.AddFunc("divideFloat", divideFloat)
	engine.AddFunc("modFloat", modFloat)
	engine.AddFunc("absFloat", absFloat)
	engine.AddFunc("maxFloat", maxFloat)
	engine.AddFunc("minFloat", minFloat)
	engine.AddFunc("uuid", uuidHelper)
	engine.AddFunc("asInt", asInt)
	engine.AddFunc("asFloat", asFloat)
	engine.AddFunc("asString", asString)

	// Type checking helpers
	engine.AddFunc("isFloat", isFloat)
	engine.AddFunc("isInt", isInt)
	engine.AddFunc("isBool", isBool)
	engine.AddFunc("isDate", isDate)
	engine.AddFunc("isDatetime", isDatetime)
	engine.AddFunc("isTime", isTime)
	engine.AddFunc("isTimestamp", isTimestamp)
	engine.AddFunc("isTimestampWithTimezone", isTimestampWithTimezone)
	engine.AddFunc("isNumeric", isNumeric)

	// Number formatting helpers
	engine.AddFunc("formatIntegerWithCommas", formatIntegerWithCommas)

	// Initialize the Fiber application with configuration options
	// Fiber is a fast, Express-inspired web framework for Go
	app := fiber.New(fiber.Config{
		Views:             engine,                 // Set the template engine
		PassLocalsToViews: true,                   // Pass local variables to views
		Prefork:           appConfig.Prefork,      // Enable/disable prefork based on config
		ServerHeader:      appConfig.ServerHeader, // Custom server header
		StrictRouting:     false,                  // Disable strict routing
		AppName:           "Eocto",                // Application name
		EnablePrintRoutes: false,                  // Print routes in debug mode
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Handle errors gracefully
			code := fiber.StatusInternalServerError
			e, ok := err.(*fiber.Error)
			if ok {
				code = e.Code
			}
			if code == fiber.StatusNotFound {
				return c.Status(fiber.StatusNotFound).Render("404", fiber.Map{
					"Title": "Page Not Found",
					"Path":  c.Path(),
				})
			}
			return c.Status(code).SendString("Internal Server Error")
		},
	})

	// Setup the middleware to retrieve the data sent in first GET request
	//app.Use(func(c *fiber.Ctx) error {
	//	// IsWebSocketUpgrade returns true if the client
	//	// requested upgrade to the WebSocket protocol.
	//	if websocket.IsWebSocketUpgrade(c) {
	//		c.Locals("allowed", true)
	//		return c.Next()
	//	}
	//	//return c.Next()
	//	return fiber.ErrUpgradeRequired
	//})
	// Multiple event handling supported
	//socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
	//	//session_id:=ep.Kws.Cookies("session_id")
	//	session_id := ep.Kws.Locals("session_id")
	//	fmt.Printf("Connection event 1 - Session_id: % +v\n", session_id)
	//	fmt.Printf("Connection event 1 - User: %s\n", ep.Kws.GetStringAttribute("user_id"))
	//})
	//// On message event
	//socketio.On(socketio.EventMessage, func(ep *socketio.EventPayload) {
	//	session_id := ep.Kws.Locals("session_id")
	//	fmt.Printf("Message event 1 - Session_id: % +v\n", session_id)
	//	userID := ep.Kws.GetStringAttribute("user_id")
	//	logr.Info(fmt.Sprintf("UserID %s", userID))
	//	clients.UpdateLastSeen(userID)
	//
	//	logr.Info(fmt.Sprintf("📨 Message event - User: %s - Message: %s\n", userID, string(ep.Data)))
	//
	//	message := MessageObject{}
	//	err := json.Unmarshal(ep.Data, &message)
	//	if err != nil {
	//		logr.Error(fmt.Sprintf("❌ Error unmarshaling message: %v\n", err))
	//		return
	//	}
	//
	//	// Fire custom event
	//	if message.Event != "" {
	//		ep.Kws.Fire(message.Event, []byte(message.Data))
	//	}
	//
	//	// Emit to target user using shared clients
	//	if targetUUID, exists := clients.GetClientUUID(message.To); exists {
	//		err = ep.Kws.EmitTo(targetUUID, ep.Data, socketio.TextMessage)
	//		if err != nil {
	//			logr.Error(fmt.Sprintf("❌ Error emitting to user %s: %v\n", message.To, err))
	//		}
	//	} else {
	//		logr.Warn(fmt.Sprintf("⚠️ Target user %s not found in clients\n", message.To))
	//	}
	//})
	//
	//// On disconnect event
	//socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
	//	userID := ep.Kws.GetStringAttribute("user_id")
	//	clients.RemoveClient(userID)
	//	logr.Info(fmt.Sprintf("❌ Disconnection event - User: %s, Remaining clients: %d\n", userID, clients.GetConnectedCount()))
	//})
	//
	//// On close event
	//// This event is called when the server disconnects the user actively with .Close() method
	//socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
	//	userID := ep.Kws.GetStringAttribute("user_id")
	//	clients.RemoveClient(userID)
	//	logr.Info(fmt.Sprintf("🔒 Close event - User: %s\n", userID))
	//})
	//
	//// On error event
	//socketio.On(socketio.EventError, func(ep *socketio.EventPayload) {
	//	fmt.Printf("Error event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	//})

	//app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
	//	for {
	//		mt, msg, err := c.ReadMessage()
	//		if err != nil {
	//			log.Println("read:", err)
	//			break
	//		}
	//		log.Printf("recv: %s", msg)
	//		err = c.WriteMessage(mt, msg)
	//		if err != nil {
	//			log.Println("write:", err)
	//			break
	//		}
	//	}
	//}))

	// Add recovery middleware to handle panics gracefully
	// This prevents the server from crashing when a panic occurs in a handler
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true, // Enable stack trace for better debugging
	}))

	// Add logger middleware to log HTTP requests
	// This logs information about each request such as method, path, status code, and response time
	app.Use(logger.New())

	// Add favicon middleware to serve favicon.ico efficiently
	// This prevents unnecessary processing of favicon requests
	app.Use(favicon.New())

	// Add health check middleware for monitoring server health
	// This adds a /healthcheck endpoint that returns server status
	app.Use(healthcheck.New())

	// Add metrics endpoint for monitoring server performance
	// This adds a /metrics endpoint with real-time server statistics
	app.Get("/metrics", monitor.New())

	// Add compression middleware to reduce response size
	// This compresses responses using gzip or other algorithms to save bandwidth
	app.Use(compress.New(
		compress.Config{
			Level: compress.LevelDefault, // Maximum compression level
		},
	))

	// Parse the modules configuration from the configuration source
	// This loads all available modules that should be initialized and registered with the application
	modules, err := config.ParseModulesConfig()
	if err != nil {
		// If module configuration cannot be parsed, the application cannot continue
		// as modules are essential components of the system architecture
		log.Fatal(err)
	}
	// Iterate through all loaded modules and set up their respective routes
	// Each module's routes are configured using the SetupRoutes function, which maps endpoints
	// to the application. If route setup fails for any module, the server will terminate
	for _, module := range modules {
		// logr.Info("Setting up routes for module: % +v", module)
		// Set up the routes for the current module by registering them with the Fiber application
		// This connects the module's handlers to specific HTTP endpoints
		err := routes.SetupRoutes(app, module)
		if err != nil {
			// If routes cannot be set up for a module, the application cannot function correctly
			// This is a fatal error that requires immediate attention
			log.Fatal(err)
		}
	}

	// Configure static file serving for public assets
	// This serves files from the ./public directory at the /public URL path
	app.Static("/public", "./public")

	// Configure static file serving for images
	// This serves image files from the ./public/img directory at the /images URL path
	app.Static("/images", "./public/img")

	// Determine the port to listen on from environment variables or configuration
	// This allows for flexible deployment in different environments
	port := os.Getenv("PORT")
	if port == "" {
		if appConfig.Port == "" {
			port = "3000" // Default port if not specified
		} else {
			port = appConfig.Port
		}
	}

	// Start the server and listen for incoming connections
	// This blocks until the server is shut down
	log.Fatal(app.Listen(":" + port))
}
