package utilities

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ProxyRequest proxies the request to the target URL
// ProxyRequest proxies an incoming HTTP request to a target URL, forwarding the original request's method, headers, and body.
// It handles header filtering, sets X-Forwarded headers, and copies the response back to the original client.
// Returns an error if the target URL is invalid, request creation fails, or proxying encounters an error.
func ProxyRequest(c *fiber.Ctx, targetURL string) error {
	// Parse target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		return c.Status(500).SendString("Invalid target URL")
	}

	// Create new request
	proxyURL := target.Scheme + "://" + target.Host + c.OriginalURL()

	req, err := http.NewRequest(c.Method(), proxyURL, strings.NewReader(string(c.Body())))
	if err != nil {
		return c.Status(500).SendString("Failed to create request")
	}

	// Copy headers from original request
	c.Request().Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		// Skip certain headers that shouldn't be forwarded
		if !shouldSkipHeader(keyStr) {
			req.Header.Set(keyStr, string(value))
		}
	})

	// Set X-Forwarded headers
	req.Header.Set("X-Forwarded-For", c.IP())
	req.Header.Set("X-Forwarded-Proto", c.Protocol())
	req.Header.Set("X-Forwarded-Host", c.Hostname())

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(502).SendString("Bad Gateway")
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Response().Header.Add(key, value)
		}
	}

	// Set status code
	c.Status(resp.StatusCode)

	// Copy response body
	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	return err
}

// shouldSkipHeader determines if a header should be skipped when proxying
// shouldSkipHeader determines whether a specific HTTP header should be skipped during proxying.
// It checks the header against a predefined list of headers that should not be forwarded.
// Returns true if the header should be skipped, false otherwise.
func shouldSkipHeader(header string) bool {
	skipHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	header = strings.ToLower(header)
	for _, skip := range skipHeaders {
		if strings.ToLower(skip) == header {
			return true
		}
	}
	return false
}
