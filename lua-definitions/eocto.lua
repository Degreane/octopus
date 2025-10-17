---@meta

---Global eocto table containing all exposed functions
---@class eocto
eocto = {}

---Debug function for logging messages
---@param level "info"|"warning"|"error"|"important" Log level
---@param message string|number|boolean|table Message to log
function eocto.debug(level, message) end

---Get session value by key
---@param key string Session key
---@return string|nil value Session value or nil if not found
function eocto.getSession(key) end

---Set session value
---@param key string Session key
---@param value string Session value
---@param expiry? number Optional expiry time in seconds
function eocto.setSession(key, value, expiry) end

---Delete session key
---@param key string Session key to delete
function eocto.deleteSession(key) end

---Set session expiry time
---@param expiry number Expiry time in seconds
function eocto.setSessionExpiry(expiry) end

---Get cookie value
---@param name string Cookie name
---@return string|nil value Cookie value or nil if not found
function eocto.getCookie(name) end

---Set cookie
---@param name string Cookie name
---@param value string Cookie value
---@param options? table Optional cookie options (expires, path, domain, etc.)
function eocto.setCookie(name, value, options) end

---Get all cookies
---@return table cookies Table of all cookies
function eocto.getAllCookies() end

---Delete cookie
---@param name string Cookie name to delete
function eocto.deleteCookie(name) end

---Clear all cookies
function eocto.clearAllCookies() end

---Add user to WebSocket room
---@param userId string User ID
---@param roomId string Room ID
---@return boolean success Success status
---@return string message Info message
function eocto.wsAddRoom(userId, roomId) end

---Remove user from WebSocket room
---@param userId string User ID
---@param roomId string Room ID
---@return boolean success Success status
---@return string message Info message
function eocto.wsRemoveRoom(userId, roomId) end

---Get all rooms for a user
---@param userId string User ID
---@return table rooms List of room names
---@return string message Info message
function eocto.wsGetUserRooms(userId) end

---Check if user is in a specific room
---@param userId string User ID
---@param roomId string Room ID
---@return boolean inRoom True if user is in the room
---@return string message Info message
function eocto.wsIsUserInRoom(userId, roomId) end

---Emit a message to all users in a specific room
---@param roomId string Room ID
---@param event string Event name
---@param data any Data to send (string|number|boolean|table|nil)
---@param excludeUsers? table Optional array of user IDs to exclude
---@return number deliveredCount Number of clients that received the message
---@return string message Info message
function eocto.wsEmitToRoom(roomId, event, data, excludeUsers) end

---Get CSRF token
---@return string token CSRF token
function eocto.getCsrfToken() end

---Get local context value
---@param key string Local key
---@return any value Local value
function eocto.getLocal(key) end

---Set local context value
---@param key string Local key
---@param value any Local value
function eocto.setLocal(key, value) end

---Delete local context value
---@param key string Local key to delete
function eocto.deleteLocal(key) end

---Get all local context values
---@return table locals Table of all locals
function eocto.getLocals() end

---Get all request headers
---@return table headers Table of all headers
function eocto.getHeaders() end

---Get specific request header
---@param name string Header name
---@return string|nil value Header value or nil if not found
function eocto.getHeader(name) end

---Set response header
---@param name string Header name
---@param value string Header value
function eocto.setHeader(name, value) end

---Delete response header
---@param name string Header name to delete
function eocto.deleteHeader(name) end

---Get request path
---@return string path Request path
function eocto.getPath() end

---Get request host
---@return string host Request host
function eocto.getHost() end

---Get request schema (http/https)
---@return string schema Request schema
function eocto.getSchema() end

---Get query parameters
---@return table params Table of query parameters
function eocto.getQueryParams() end

---Get path parameters
---@return table params Table of path parameters
function eocto.getPathParams() end

---Get specific path parameter
---@param name string Parameter name
---@return string|nil value Parameter value or nil if not found
function eocto.getPathParam(name) end

---Get POST request body
---@return string body Request body
function eocto.getPostBody() end

---Get HTTP method
---@return string method HTTP method (GET, POST, etc.)
function eocto.getMethod() end

---Decrypt data
---@param data string Encrypted data
---@return string|nil decrypted Decrypted data or nil on error
function eocto.decryptData(data) end

---Encrypt data
---@param data string Data to encrypt
---@return string|nil encrypted Encrypted data or nil on error
function eocto.encryptData(data) end

---Decode JSON string
---@param json string JSON string
---@return table|nil decoded Decoded table or nil on error
function eocto.decodeJSON(json) end

---Encode table to JSON
---@param data table Table to encode
---@return string|nil json JSON string or nil on error
function eocto.encodeJSON(data) end

---Encode string to Base32
---@param data string Data to encode
---@return string encoded Base32 encoded string
function eocto.encodeBase32(data) end

---Decode Base32 string
---@param data string Base32 encoded string
---@return string|nil decoded Decoded string or nil on error
function eocto.decodeBase32(data) end

---Get data from MongoDB collection
---@param collection string Collection name
---@param query table Query parameters
---@return table|nil result Query result or nil on error
function eocto.getDataFromCollection(collection, query) end

---Set/Update data in MongoDB collection
---@param collection string Collection name
---@param query table Query parameters
---@param data table Data to set
---@return boolean success Success status
function eocto.setDataToCollection(collection, query, data) end

---Delete data from MongoDB collection
---@param collection string Collection name
---@param query table Query parameters
---@return boolean success Success status
function eocto.delDataFromCollection(collection, query) end

---Insert data to MongoDB collection
---@param collection string Collection name
---@param data table Data to insert
---@return string|nil id Inserted document ID or nil on error
function eocto.insertDataToCollection(collection, data) end

---Make HTTP request
---@param method string HTTP method
---@param url string Request URL
---@param headers? table Optional headers
---@param body? string Optional request body
---@return table response Response table with status, headers, body
function eocto.makeRequest(method, url, headers, body) end

---Proxy request to another server
---@param url string Target URL
---@param options? table Optional proxy options
---@return boolean success Success status
function eocto.proxy(url, options) end

---Send WhatsApp message
---@param to string Recipient phone number
---@param message string Message content
---@return boolean success Success status
function eocto.sendWhatsAppMessage(to, message) end

---Set HTTP response
---@param status number HTTP status code
---@param body string|table Response body
---@param headers? table Optional response headers
function eocto.setResponse(status, body, headers) end

---Render template
---@param template string Template name
---@param data? table Optional template data
---@return boolean success Success status
function eocto.render(template, data) end

---Render JSON response
---@param data table Data to render as JSON
---@param status number HTTP status code
---@return boolean success Success status
function eocto.renderJson(data, status) end

---Generate UUID v4
---@return string uuid Generated UUID
function eocto.getUUID() end

---Get module settings
---@return table settings Table containing BasePath and LocalPath
function eocto.getSettings() end

---Get Redis value
---@param key string Redis key
---@return string|nil value Redis value or nil if not found
function eocto.getRedis(key) end

---Set Redis value
---@param key string Redis key
---@param value string Redis value
---@param expiry? number Optional expiry time in seconds
---@return boolean success Success status
function eocto.setRedis(key, value, expiry) end

---Delete Redis key
---@param key string Redis key to delete
---@return boolean success Success status
function eocto.deleteRedis(key) end

---Get current timestamp in nanoseconds
---@return number timestamp Unix timestamp in nanoseconds
function eocto.timeStampNano() end

---Get current timestamp in milliseconds
---@return number timestamp Unix timestamp in milliseconds
function eocto.timeStampMilli() end

---Get current timestamp in seconds
---@return number timestamp Unix timestamp in seconds
function eocto.timeStamp() end

---Get current working directory
---@return string cwd Current working directory path
function eocto.getCWD() end


---Reset the session working directory to the server's current process working directory
function eocto.resetWD() end

---Set the session working directory
---@param path string Absolute or relative directory path
function eocto.setWD(path) end

---List files and directories in a path or the current working directory
---@param path? string Optional path to list; defaults to the session working directory
---@return table|nil items Array of names or nil on error
function eocto.listFiles(path) end

---Read a YAML file and return its content as a JSON string
---@param path string File path
---@return string|nil json JSON string or nil on error
function eocto.readYamlFile(path) end

---Read a CSV file and return its content as a JSON string (array of objects)
---@param path string File path
---@return string|nil json JSON string or nil on error
function eocto.readCsvFile(path) end
