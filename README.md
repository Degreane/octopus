# Octopus

**High-Performance Web Server Built with GoLang, Fiber, and Embedded Lua**

Octopus is a modern, high-performance web server designed for speed, scalability, and ease of development. Built using Go (Golang) and the Fiber framework, Octopus leverages the efficiency of Go's concurrency model and Fiber's Express.js-like simplicity to deliver a robust backend solution for web applications, APIs, and microservices. What sets Octopus apart is its deep integration of Lua scripting, allowing dynamic server logic customization without recompilation.

---

## Table of Contents

- [Why Octopus?](#why-octopus)
- [Requirements](#requirements)
- [Key Dependencies](#key-dependencies)
- [Architecture](#architecture)
- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Documentation](#documentation)
- [Development](#development)
- [License](#license)

---

## Why Octopus?

### Go + Fiber + Lua: The Perfect Stack

By embedding Lua as a first-class scripting language, Octopus enables developers to:

- **Modify request/response handling on the fly** - Change server behavior without recompilation
- **Implement custom middleware logic without restarting** - Update business rules dynamically
- **Execute data transformations, validations, and routing rules via Lua scripts** - Flexible data processing
- **Seamlessly map Go functions to Lua** - Hybrid Go-Lua workflows for optimal performance

### Performance Foundation

1. **Go + Fiber Efficiency**
   - Go's compiled nature eliminates interpreter overhead
   - Fiber's Fasthttp engine handles 100K+ requests/sec with minimal latency
   - Significantly faster than interpreted languages (Python, Ruby, PHP)

2. **Built for Concurrency**
   - Goroutines enable thousands of concurrent connections
   - Non-blocking I/O ensures smooth handling of high-traffic loads
   - Lightweight architecture with minimal resource footprint

3. **Developer-Friendly**
   - Fiber's Express-like API for familiar development patterns
   - Minimal boilerplate code
   - YAML-based configuration for human-readable routing

4. **Scalability & Microservices-Ready**
   - Stateless design for horizontal scaling
   - Lower RAM consumption than Node.js or Java-based servers
   - Container-friendly (Docker, Kubernetes)

5. **Embedded Lua: The Dynamic Layer**
   - Runtime scripting without recompilation
   - Go-Lua interoperability with bidirectional function calls
   - Sandboxed execution for security

---

## Requirements

- **Go Version**: 1.24.2 or higher
- **Air** (optional, for hot reloading during development)

---

## Key Dependencies

### Core Framework & Scripting
```
github.com/gofiber/fiber/v2 v2.52.9              # High-performance web framework
github.com/yuin/gopher-lua v1.1.1                # Lua VM for Go
```

### Database & Caching
```
go.mongodb.org/mongo-driver v1.17.4              # MongoDB driver
github.com/redis/go-redis/v9 v9.12.1             # Redis client
github.com/gofiber/storage/redis/v3 v3.4.1       # Redis storage adapter
github.com/gofiber/storage/memory/v2 v2.1.0      # In-memory storage
```

### Real-Time Communication
```
github.com/gofiber/contrib/socketio v1.1.6       # Socket.IO support
github.com/gofiber/contrib/websocket v1.3.4      # WebSocket integration
```

### Template Engine & Configuration
```
github.com/gofiber/template/html/v2 v2.1.3       # HTML template engine
gopkg.in/yaml.v3 v3.0.1                          # YAML configuration parser
github.com/joho/godotenv v1.5.1                  # Environment configuration
```

### External Services
```
github.com/twilio/twilio-go v1.27.2              # Twilio API integration
```

### Utilities
```
github.com/google/uuid v1.6.0                    # UUID generation
github.com/golang-jwt/jwt/v5 v5.3.0              # JWT authentication
```

For a complete list of dependencies including transitive packages, see `go.mod`.

---

## Architecture

### YAML-Based Modular Routing

Octopus uses a configuration-driven architecture where routes, middleware, and module settings are defined in YAML files. This approach provides:

- **Human-Readable Configuration** - Clean, intuitive syntax for developers and system administrators
- **Hierarchical Structure** - Perfect for nested configurations like routes and middleware chains
- **Version Control Friendly** - Easy to track changes and merge configurations
- **No Code Deployment** - Modify routes and middleware without recompiling
- **Environment Agnostic** - Same format across development, staging, and production
- **Validation Ready** - Schema validation ensures configuration integrity

#### Module Configuration Structure

```yaml
- Name: "EzeKod Octopus"
  Description: "EzeKod Octopus is a powerful and versatile module..."
  Version: 25.10.16
  Author: Faysal AL-Banna
  Email: degreane@gmail.com
  Website: https://www.ezekod.com
  License: GPL3
  Dependencies: [redis, mongodb, lua, tailwindcss]
  Settings: []
  BasePath: /OC
  db: octopus
  Routes:
    - method: GET
      path: /
      preCheck:
        - script: auth/isLoggedIn.lua
        - script: parseHeaders.lua
        - script: currentPage.lua
        - script: session.lua
      view: pages/home
    - method: GET
      path: /ws/:id
      websocket: true
      preCheck:
        - script: auth/isLoggedIn.lua
      view: pages/ws
```

#### Request Processing Pipeline

1. **HTTP Request** → Incoming request from client
2. **Route Matching** → Fiber router matches against configured routes
3. **Lua PreCheck** → Sequential execution of middleware scripts
4. **View Render** → Template rendering with accumulated context data

#### PreCheck Scripts (Embedded Lua Middleware)

PreCheck scripts form a powerful middleware chain that executes before the main route handler. Each script can:

- Modify request context
- Validate permissions
- Transform data
- Terminate request early

**Example Scripts:**

- `auth/isLoggedIn.lua` - Session validation and authentication checks
- `parseHeaders.lua` - HTTP header extraction and processing
- `currentPage.lua` - Page context and navigation state setup
- `session.lua` - User session data and preferences management

### URL Structure

Routes are namespaced using the `BasePath` property:

```
BasePath: /OC  +  path: /         = https://example.com/OC/
BasePath: /OC  +  path: /lua      = https://example.com/OC/lua
BasePath: /OC  +  path: /api/users = https://example.com/OC/api/users
```

### Template Resolution

```
BasePath: /OC
View: pages/home
Template Path: views/OC/pages/home.html
```

---

## Features

### 🚀 High-Performance Web Framework
- Built on Fiber and Fasthttp for maximum throughput
- 100K+ requests/sec capability
- Non-blocking I/O and goroutine-based concurrency
- Efficient memory pooling and script caching

### 🔧 50+ Template Helper Functions
Template helpers available in HTML templates for:
- **Type Checking**: `isInt`, `isFloat`, `isNumeric`, `isBool`, `isDate`, `isDatetime`, `isTime`, `isTimestamp`
- **Formatting**: `FormatFloatWithCommas`, `formatIntegerWithCommas`, `formatNum`
- **Arithmetic**: Integer operations (`addInt`, `subtractInt`, `multiplyInt`, `divideInt`, `modInt`)
- **Arithmetic**: Float operations (`addFloat`, `subtractFloat`, `multiplyFloat`, `divideFloat`)
- **Utilities**: `dict`, `length`, `checkType`, `timestamp`, `uuid`
- **Iteration**: `iterate`, `iterateFilter`, `iterateRange`, `iterateEven`, `iterateOdd`
- **Random**: `rand`, `randRange`, `randFloat`, `randFloatRange`
- **Conversion**: `asInt`, `asFloat`, `asString`

### 🌙 Comprehensive Lua API (50+ Functions)

**Debug Functions**
- `eocto.debug(level, message)` - Structured logging with severity levels

**Session Management**
- `eocto.getSession(key)`, `eocto.setSession(key, value)`, `eocto.deleteSession(key)`
- `eocto.setSessionExpiry(seconds)` - Configure session lifetime

**Cookie Management**
- `eocto.getCookie(name)`, `eocto.setCookie(name, value, options)`
- `eocto.getAllCookies()`, `eocto.deleteCookie(name)`, `eocto.clearAllCookies()`

**Security & Encryption**
- `eocto.getCsrfToken()` - CSRF protection tokens
- `eocto.encryptData(data)`, `eocto.decryptData(encryptedData)` - Server-side encryption

**HTTP & Request Handling**
- `eocto.getMethod()`, `eocto.getPath()`, `eocto.getHost()`, `eocto.getSchema()`
- `eocto.getQueryParams()`, `eocto.getPathParams()`, `eocto.getPathParam(key)`
- `eocto.getPostBody()` - Request body parsing
- `eocto.makeRequest(options)` - External HTTP calls
- `eocto.proxy(options)` - Request proxying

**Headers Management**
- `eocto.getHeaders()`, `eocto.getHeader(name)`
- `eocto.setHeader(name, value)`, `eocto.deleteHeader(name)`

**Local Variables (Request Context)**
- `eocto.getLocal(key)`, `eocto.setLocal(key, value)`
- `eocto.deleteLocal(key)`, `eocto.getLocals()`

**Data Encoding & JSON**
- `eocto.decodeJSON(jsonString)`, `eocto.encodeJSON(table)`
- `eocto.encodeBase32(data)`, `eocto.decodeBase32(encodedData)`

**MongoDB Operations**
- `eocto.getDataFromCollection(collection, query)` - Query documents
- `eocto.setDataToCollection(collection, filter, update)` - Update documents
- `eocto.delDataFromCollection(collection, filter)` - Delete documents
- `eocto.insertDataToCollection(collection, document)` - Insert documents

**Redis Cache Operations**
- `eocto.getRedis(key)`, `eocto.setRedis(key, value, expiration)`
- `eocto.deleteRedis(key)` - Cache management

**Response & Rendering**
- `eocto.setResponse(statusCode, data)` - JSON responses
- `eocto.render(template, data)` - HTML template rendering
- `eocto.renderJson(data)` - Direct JSON output

**Utility Functions**
- `eocto.getUUID()` - UUID v4 generation
- `eocto.timeStamp()`, `eocto.timeStampMilli()`, `eocto.timeStampNano()` - Timestamps
- `eocto.getSettings()` - Module configuration access (BasePath, LocalPath)

**Communication Functions**
- `eocto.sendWhatsAppMessage(options)` - Twilio WhatsApp integration

**WebSocket Functions**
- `eocto.wsAddRoom(roomId)`, `eocto.wsRemoveRoom(roomId)` - Room management
- `eocto.wsGetUserRooms()`, `eocto.wsIsUserInRoom(roomId)` - Room queries
- `eocto.wsEmitToRoom(roomId, event, data)` - Broadcast to rooms

### 🗄️ Database Integration
- **MongoDB**: Full document database support with CRUD operations
- **Redis**: High-speed caching with expiration management
- Session storage with Redis or in-memory backends

### 🔌 WebSocket Support
- Real-time bidirectional communication
- Room-based broadcasting
- User connection tracking
- Socket.IO integration

### 🎨 HTML Template Engine
- Go html/template with custom helpers
- Partial template support
- Dynamic data binding
- CSRF token integration

### 🔐 Security Features
- Session management with expiration
- CSRF protection
- Data encryption/decryption
- Sandboxed Lua execution
- Secure cookie handling

---

## Installation

### Clone the Repository
```bash
git clone https://github.com/degreane/octopus.git
cd octopus
```

### Install Dependencies
```bash
go mod download
```

### Build the Server
```bash
go build ./cmd/server
```

### Run the Server
```bash
./server
```

---

## Configuration

### Environment Variables
Create a `.env` file in the project root:

```env
# Server Configuration
PORT=3000
HOST=localhost

# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DB=octopus

# Redis Configuration
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Twilio Configuration (optional)
TWILIO_ACCOUNT_SID=your_account_sid
TWILIO_AUTH_TOKEN=your_auth_token
TWILIO_PHONE_NUMBER=your_twilio_number
```

### Module Configuration
Module routes and settings are defined in YAML files. See the [YAML Configuration Documentation](#documentation) for detailed structure and examples.

---

## Documentation

Octopus includes comprehensive built-in documentation accessible via the web interface:

### Main Documentation (`/OC/`)
- Project overview and architecture
- Go + Fiber + Lua integration benefits
- Performance characteristics
- Scalability features

### Lua API Reference (`/OC/lua`)
- Complete function reference with examples
- 50+ functions across 14 categories
- Code snippets for common use cases
- Parameter descriptions and return values

### YAML Configuration Guide
- Module configuration structure
- Route definition syntax
- PreCheck middleware chains
- Best practices and examples

### Template Helpers Reference
- 50+ helper functions with examples
- Type checking and validation
- Formatting and arithmetic operations
- Iteration and random number generation
- Real-world usage scenarios

---

## Development

### Hot Reloading with Air

For development with automatic reloading on file changes:

#### Install Air
```bash
go install github.com/cosmtrek/air@latest
```

#### Run with Air
```bash
air
```

Air will watch for changes in `.go`, `.lua`, `.yaml`, and `.html` files and automatically rebuild and restart the server.

### Project Structure

```
octopus/
├── cmd/
│   └── server/
│       ├── main.go              # Main server entry point
│       └── templateHelpers.go   # Template helper functions
├── views/
│   └── OC/
│       ├── pages/               # Page templates
│       └── partials/            # Partial templates
├── config/
│   └── modules.yaml             # Module configuration
├── scripts/                     # Lua scripts
│   └── auth/                    # Authentication scripts
├── go.mod                       # Go dependencies
└── .env                         # Environment configuration
```

### Adding Custom Lua Functions

1. Define the function in Go (in `main.go` or a separate file)
2. Register it with the Lua VM in the initialization code
3. Use it in Lua scripts via the `eocto` namespace

### Adding Template Helpers

1. Add the function to `cmd/server/templateHelpers.go`
2. Register it in `main.go` using `engine.AddFunc("functionName", functionReference)`
3. Use it in templates with `{{ functionName args }}`

### Creating New Routes

1. Edit your module's YAML configuration file
2. Add a new route entry with method, path, and view
3. Optionally add PreCheck scripts for middleware
4. Create the corresponding view template

---

## Performance & Optimization

### Script Caching
- Lua scripts are compiled once and reused across requests
- Minimizes parsing overhead for high-throughput scenarios

### Memory Pooling
- Lua state objects are pooled and reused
- Reduces garbage collection pressure

### Async Operations
- Database and external API calls handled asynchronously
- Maintains responsiveness under load

---

## Real-World Use Cases

- **Authentication & Authorization**: Dynamic permission checks, role-based access control
- **Data Validation & Processing**: Input sanitization, business rule validation, data transformation pipelines
- **API Rate Limiting**: Dynamic throttling, quota management, abuse prevention
- **Content Personalization**: User-specific filtering, A/B testing, feature flags
- **Integration Workflows**: Third-party API calls, webhook processing, external service integration
- **Real-Time Systems**: WebSocket-based chat, notifications, live dashboards
- **Microservices Gateways**: Request routing, load balancing, service orchestration

---

## License

**GPL3**

Copyright © 2025 Faysal AL-Banna

[//]: # (Website: https://www.ezekod.com)
Email: degreane@gmail.com

---

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

---

## Support

For questions, issues, or feature requests, please open an issue on the GitHub repository or contact the author directly.

---

**Octopus Server • Eocto 0.23.1.25**
*Configuration-First • Script-Powered • Performance-Optimized • Developer-Friendly*

