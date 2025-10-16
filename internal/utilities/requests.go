package utilities

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/degreane/octopus/internal/utilities/debug"
	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

// MakeHTTPRequest performs an HTTP request and returns the complete response details.
// It handles both successful and failed requests, returning appropriate response data.
//
// Parameters:
//   - method: string - HTTP method (GET, POST, PUT, DELETE, etc.)
//   - url: string - The target URL for the request
//   - headers: map[string]string - Request headers to be set
//   - body: io.Reader - Optional request body (nil for no body)
//
// Returns:
//   - map[string]interface{} - Response containing:
//   - status: int - HTTP status code
//   - headers: map[string][]string - Response headers
//   - cookies: []string - Response cookies
//   - body: string - Response body
//   - error: string - Error message (if any)
//
// Example:
//
//	response := MakeHTTPRequest("GET", "https://api.example.com",
//	    map[string]string{"Authorization": "Bearer token"}, nil)
func MakeHTTPRequest(method, url string, headers map[string]string, body io.Reader) map[string]interface{} {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return map[string]interface{}{
			"status": 500,
			"error":  err.Error(),
		}
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"status": 500,
			"error":  err.Error(),
		}
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	var bodyContent interface{}

	// Handle different content types
	switch {
	case strings.HasPrefix(contentType, "audio/"),
		strings.HasPrefix(contentType, "video/"),
		strings.HasPrefix(contentType, "application/octet-stream"),
		strings.HasPrefix(contentType, "application/pdf"),
		strings.HasPrefix(contentType, "image/"):
		// For large files and streams, read in chunks
		chunks := []byte{}
		buffer := make([]byte, 32*1024) // 32KB chunks
		for {
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				chunks = append(chunks, buffer[:n]...)
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return map[string]interface{}{
					"status": 500,
					"error":  err.Error(),
				}
			}
		}
		bodyContent = base64.StdEncoding.EncodeToString(chunks)
	default:
		// Text content
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return map[string]interface{}{
				"status": 500,
				"error":  err.Error(),
			}
		}
		bodyContent = string(responseBody)
	}

	cookies := []string{}
	for _, cookie := range resp.Cookies() {
		cookies = append(cookies, cookie.String())
	}

	return map[string]interface{}{
		"status":      resp.StatusCode,
		"headers":     resp.Header,
		"cookies":     cookies,
		"body":        bodyContent,
		"contentType": contentType,
		"error":       "",
	}
}

// GetRequest returns a Lua function that performs HTTP requests.
// The function is exposed to Lua scripts and provides access to the HTTP client functionality.
//
// Parameters passed from Lua:
//   - method: string - HTTP method (GET, POST, etc.)
//   - url: string - Target URL for the request
//   - headers: table - Request headers (optional)
//   - body: string - Request body (optional)
//
// Returns to Lua:
//   - table containing:
//   - status: number - HTTP status code
//   - headers: table - Response headers
//   - cookies: table - Response cookies
//   - body: string - Response body (base64 encoded for binary data)
//   - contentType: string - Response content type
//   - error: string - Error message if any
//
// Usage in Lua:
//
//	local response = request("GET", "https://api.example.com", {
//	    ["Authorization"] = "Bearer token"
//	})
func GetRequest(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		method := L.ToString(1)
		url := L.ToString(2)

		// Convert Lua headers table to Go map
		headers := make(map[string]string)
		if L.GetTop() >= 3 && L.Get(3) != lua.LNil {
			luaHeaders := L.ToTable(3)
			luaHeaders.ForEach(func(k, v lua.LValue) {
				headers[k.String()] = v.String()
			})
		}
		debug.Debug(debug.Important, fmt.Sprintf("Headers: %+v", headers))

		// Handle optional body
		var body io.Reader
		if L.GetTop() >= 4 && L.GetTop() >= 5 {
			log.Printf("Body: Buff %s", L.ToString(4))
			body = bytes.NewBufferString(L.ToString(4))
		} else if L.GetTop() >= 4 {
			log.Printf("Body: Str %s", L.ToString(4))
			body = strings.NewReader(L.ToString(4))
		}

		response := MakeHTTPRequest(method, url, headers, body)

		// Convert response to Lua table
		responseTable := L.NewTable()
		for k, v := range response {
			switch val := v.(type) {
			case int:
				responseTable.RawSetString(k, lua.LNumber(val))
			case string:
				responseTable.RawSetString(k, lua.LString(val))
			case http.Header:
				headerTable := L.NewTable()
				for hk, hv := range val {
					headerTable.RawSetString(hk, lua.LString(strings.Join(hv, ", ")))
				}
				responseTable.RawSetString(k, headerTable)
			case []string:
				cookieTable := L.NewTable()
				for i, cookie := range val {
					cookieTable.RawSetInt(i+1, lua.LString(cookie))
				}
				responseTable.RawSetString(k, cookieTable)
			}
		}

		L.Push(responseTable)
		return 1
	}
}

// Add this to your existing requests.go file

func ProxyRequestLua(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		targetURL := L.ToString(1)

		// Custom headers
		customHeaders := make(map[string]string)
		if L.GetTop() >= 2 && L.Get(2) != lua.LNil {
			luaHeaders := L.ToTable(2)
			luaHeaders.ForEach(func(k, v lua.LValue) {
				customHeaders[k.String()] = v.String()
			})
		}

		// Path rewriting
		var rewritePath string
		if L.GetTop() >= 3 && L.Get(3) != lua.LNil {
			rewritePath = L.ToString(3)
		}

		// Skip TLS verification option
		skipTLS := false
		if L.GetTop() >= 4 && L.Get(4) != lua.LNil {
			skipTLS = bool(L.ToBool(4))
		}

		err := ProxyRequestWithTLS(c, targetURL, customHeaders, rewritePath, skipTLS)

		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LBool(true))
		return 1
	}
}

func ProxyRequestWithTLS(c *fiber.Ctx, targetURL string, customHeaders map[string]string, rewritePath string, skipTLS bool) error {
	target, err := url.Parse(targetURL)
	if err != nil {
		return c.Status(500).SendString("Invalid target URL")
	}

	// Build the proxy URL
	var proxyURL string
	if rewritePath != "" {
		proxyURL = target.Scheme + "://" + target.Host + rewritePath
	} else {
		originalPath := string(c.Request().URI().Path())
		proxyURL = target.Scheme + "://" + target.Host + originalPath
	}

	if c.Request().URI().QueryString() != nil {
		proxyURL += "?" + string(c.Request().URI().QueryString())
	}

	req, err := http.NewRequest(c.Method(), proxyURL, strings.NewReader(string(c.Body())))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return c.Status(500).SendString("Failed to create request")
	}

	// Copy headers
	c.Request().Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		if !shouldSkipHeader(keyStr) {
			req.Header.Set(keyStr, string(value))
		}
	})

	// Add custom headers
	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}

	// Set forwarding headers
	req.Header.Set("X-Forwarded-For", c.IP())
	req.Header.Set("X-Forwarded-Proto", c.Protocol())
	req.Header.Set("X-Forwarded-Host", c.Hostname())
	req.Header.Set("X-Real-IP", c.IP())

	// Create client with optional TLS skip
	transport := &http.Transport{}
	if skipTLS {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return c.Status(502).SendString("Bad Gateway: " + err.Error())
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		if !shouldSkipResponseHeader(key) {
			for _, value := range values {
				c.Response().Header.Add(key, value)
			}
		}
	}

	c.Status(resp.StatusCode)
	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	return err
}
func shouldSkipResponseHeader(header string) bool {
	skipHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Transfer-Encoding",
		"Upgrade",
		"Content-Length", // Let Go handle this
	}

	header = strings.ToLower(header)
	for _, skip := range skipHeaders {
		if strings.ToLower(skip) == header {
			return true
		}
	}
	return false
}
