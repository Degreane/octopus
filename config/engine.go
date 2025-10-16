// Package config provides configuration utilities for the Octopus application.
// This file contains the template engine configuration and custom template functions.
package config

import (
	"html/template"
	"strings"

	"github.com/gofiber/template/html/v2"
)

// SetupTemplateEngine initializes and configures the HTML template engine.
// It sets up the template directory, file extension, and adds custom template functions.
// The engine is configured for development with template reloading enabled.
//
// Parameters:
//   - viewsDir: Directory containing the HTML templates
//   - extension: File extension for the templates (e.g., ".html")
//   - reload: Whether to reload templates on each request (useful for development)
//
// Returns:
//   - A configured HTML template engine ready to be used with Fiber
func SetupTemplateEngine(viewsDir string, extension string, reload bool) *html.Engine {
	// Initialize the HTML template engine with the specified directory and extension
	engine := html.New(viewsDir, extension)

	// Enable or disable template reloading based on the reload parameter
	engine.Reload(reload)

	// Register all custom template functions
	registerTemplateFunctions(engine)

	return engine
}

// registerTemplateFunctions adds all custom template functions to the template engine.
// These functions extend the template engine's capabilities for various operations
// like HTML rendering, arithmetic, string manipulation, and data structure creation.
func registerTemplateFunctions(engine *html.Engine) {
	// HTMLSafe allows inserting HTML content that won't be escaped by the template engine
	engine.AddFunc("HTMLSafe", func(s string) template.HTML {
		return template.HTML(s)
	})

	// intAdd takes a variable number of integers and returns their sum
	engine.AddFunc("intAdd", func(i ...int) int {
		var k int = 0
		for _, val := range i {
			k += val
		}
		return k
	})

	// intSub takes an initial value and subtracts all subsequent values from it
	engine.AddFunc("intSub", func(ini int, i ...int) int {
		var k int = ini
		for _, val := range i {
			k -= val
		}
		return k
	})

	// floatAdd takes a variable number of float64 values and returns their sum
	engine.AddFunc("floatAdd", func(i ...float64) float64 {
		var k float64 = 0
		for _, val := range i {
			k += val
		}
		return k
	})

	// floatSub takes an initial float64 value and subtracts all subsequent values from it
	engine.AddFunc("floatSub", func(ini float64, i ...float64) float64 {
		var k float64 = ini
		for _, val := range i {
			k -= val
		}
		return k
	})

	// pct calculates what percentage one value is of another
	engine.AddFunc("pct", func(one, two float64) float64 {
		if two == 0 {
			return 100 // Avoid division by zero
		}
		return (one / two) * 100
	})

	// gt returns true if the first value is greater than the second
	engine.AddFunc("gt", func(one, two float64) bool {
		return one > two
	})

	// fraction divides one value by another, handling division by zero
	engine.AddFunc("fraction", func(one, two float64) float64 {
		if two == 0 {
			return 1 // Avoid division by zero
		}
		return one / two
	})

	// dict creates a map from key-value pairs for use in templates
	engine.AddFunc("dict", func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, &template.Error{Name: "dict"} // Must have even number of arguments
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, &template.Error{Name: "dict"} // Keys must be strings
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	})

	// trim removes specified characters from the beginning and end of a string
	engine.AddFunc("trim", func(str1, str2 string) string {
		return strings.Trim(str1, str2)
	})
}
