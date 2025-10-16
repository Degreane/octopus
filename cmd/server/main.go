// Package main is the entry point for the Octopus server application.
// It initializes the server, sets up middleware, and starts listening for requests.
package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

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
	"github.com/gofiber/fiber/v2/utils"
	"github.com/joho/godotenv"
)

// func getAllowedOrigins() []string {
// 	originsEnv := os.Getenv("ALLOWED_ORIGINS")
// 	if originsEnv == "" {
// 		// Default origins for development
// 		return []string{
// 			"http://localhost:3000",
// 			"http://localhost:3001",
// 			"http://localhost:8080",
// 			"http://127.0.0.1:3000",
// 			"http://127.0.0.1:8080",
// 		}
// 	}

// 	origins := strings.Split(originsEnv, ",")
// 	var cleanOrigins []string
// 	for _, origin := range origins {
// 		cleanOrigins = append(cleanOrigins, strings.TrimSpace(origin))
// 	}
// 	return cleanOrigins
// }

// func isOriginAllowed(origin string) bool {
// 	if origin == "" {
// 		return true // Allow same-origin requests
// 	}

// 	allowedOrigins := getAllowedOrigins()
// 	for _, allowed := range allowedOrigins {
// 		if allowed == origin {
// 			return true
// 		}
// 	}
// 	return false
// }

// MessageObject Basic chat message object
type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

func FormatFloatWithCommas(fn interface{}, precision ...int) string {
	f, ok := convertToFloat64(fn)
	if !ok {
		return "0"
	}

	// Determine precision (default 2 if not specified)
	prec := 2
	if len(precision) > 0 {
		prec = precision[0]
	}

	// Convert to string with specified precision
	s := strconv.FormatFloat(f, 'f', prec, 64)

	// Split into integer and fractional parts
	parts := strings.Split(s, ".")
	integerPart := parts[0]
	fractionalPart := ""
	if len(parts) > 1 {
		fractionalPart = parts[1]
	}

	// Format integer part with commas
	formattedInteger := formatIntegerWithCommas(integerPart)

	// Recombine with fractional part
	if fractionalPart != "" {
		return formattedInteger + "." + fractionalPart
	}
	return formattedInteger
}
func isFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
func isInt(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}
func isBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}
func isDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}
func isDatetime(s string) bool {
	_, err := time.Parse("2006-01-02T15:04:05", s)
	return err == nil
}
func isTime(s string) bool {
	_, err := time.Parse("15:04:05", s)
	return err == nil
}
func isTimestamp(s string) bool {
	_, err := time.Parse("2006-01-02T15:04:05.000000", s)
	return err == nil
}
func isTimestampWithTimezone(s string) bool {
	_, err := time.Parse("2006-01-02T15:04:05.000000Z", s)
	return err == nil
}

// check if string is a float or integer
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return true
	}
	_, err = strconv.ParseInt(s, 10, 64)
	if err == nil {
		return true
	}
	return false
}

// check if string is int return int if float return float
func checkType(s interface{}) interface{} {
	// Check if input is string type
	str, ok := s.(string)
	if ok {
		if _, err := strconv.ParseInt(str, 10, 64); err == nil {
			return "strInt"
		}
		if _, err := strconv.ParseFloat(str, 64); err == nil {
			// Check if float string has only zeros after decimal point
			parts := strings.Split(str, ".")
			if len(parts) == 2 {
				decimals := strings.TrimRight(parts[1], "0")
				if decimals == "" {
					return "strInt"
				}
			}
			return "strFloat"
		}
		return "string"
	}
	_, ok = s.(int)
	if ok {
		return "int"
	}

	_, ok = s.(float64)
	if ok {
		// Convert float to string and check decimal part
		str := fmt.Sprintf("%v", s)
		parts := strings.Split(str, ".")
		if len(parts) == 2 {
			decimals := strings.TrimRight(parts[1], "0")
			if decimals == "" {
				return "int"
			}
		}
		if len(parts) == 1 {
			return "floatInt"
		}
		return "float"
	}

	return reflect.TypeOf(s).String()
}

func formatIntegerWithCommas(s string) string {
	if len(s) <= 3 {
		return s
	}

	// Handle negative numbers
	prefix := ""
	if s[0] == '-' {
		prefix = "-"
		s = s[1:]
	}

	// Insert commas every 3 digits from right
	var b strings.Builder
	start := len(s) % 3
	if start == 0 {
		start = 3
	}

	b.WriteString(s[:start])
	for i := start; i < len(s); i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}

	return prefix + b.String()
}

// dict creates a map[string]interface{} from key-value pairs
func dict(values ...interface{}) map[string]interface{} {
	if len(values)%2 != 0 {
		// If odd number of arguments, ignore the last one
		values = values[:len(values)-1]
	}

	result := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key := fmt.Sprintf("%v", values[i])
		value := values[i+1]
		result[key] = value
	}

	return result
}
func timestamp(dateStr string) int64 {
	// Try multiple datetime formats
	formats := []string{
		"2006-01-02T15:04:05.000000",  // Your format: 2025-06-04T09:41:56.973000
		"2006-01-02T15:04:05.000000Z", // With Z timezone
		"2006-01-02T15:04:05Z",        // ISO 8601 with Z
		"2006-01-02T15:04:05",         // ISO 8601 without timezone
		"2006-01-02 15:04:05",         // SQL datetime format
		"2006-01-02T15:04:05.000Z",    // With milliseconds and Z
		"2006-01-02T15:04:05-07:00",   // With timezone offset
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Unix()
		}
	}

	// If all parsing fails, return 0
	return 0
}
func intEq(a, b interface{}) bool {
	// Convert both values to int64 for comparison
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)

	// Return true only if both are integers and equal
	return aOk && bOk && aInt == bInt
}

// intGt checks if first integer is greater than second integer
func intGt(a, b interface{}) bool {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	return aOk && bOk && aInt > bInt
}

// intGte checks if first integer is greater than or equal to second integer
func intGte(a, b interface{}) bool {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	return aOk && bOk && aInt >= bInt
}

// intLt checks if first integer is less than second integer
func intLt(a, b interface{}) bool {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	return aOk && bOk && aInt < bInt
}

// intLte checks if first integer is less than or equal to second integer
func intLte(a, b interface{}) bool {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	return aOk && bOk && aInt <= bInt
}

// intNe checks if first integer is not equal to second integer
func intNe(a, b interface{}) bool {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	return aOk && bOk && aInt != bInt
}

// addInt adds two integers
func addInt(a, b interface{}) int64 {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	if !aOk || !bOk {
		return 0
	}
	return aInt + bInt
}

// subtractInt subtracts second integer from first integer
func subtractInt(a, b interface{}) int64 {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	if !aOk || !bOk {
		return 0
	}
	return aInt - bInt
}

// multiplyInt multiplies two integers
func multiplyInt(a, b interface{}) int64 {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	if !aOk || !bOk {
		return 0
	}
	return aInt * bInt
}

// divideInt divides first integer by second integer (returns 0 if division by zero)
func divideInt(a, b interface{}) int64 {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	if !aOk || !bOk {
		return 0
	}
	return aInt / bInt
}

// modInt returns remainder of first integer divided by second integer
func modInt(a, b interface{}) int64 {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	if !aOk || !bOk {
		return 0
	}
	return aInt % bInt
}

// absInt returns absolute value of integer
func absInt(a interface{}) int64 {
	aInt, aOk := convertToInt64(a)
	if !aOk {
		return 0
	}
	if aInt < 0 {
		return -aInt
	}
	return aInt
}

// maxInt returns the larger of two integers
func maxInt(a, b interface{}) int64 {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	if !aOk || !bOk {
		return 0
	}
	if aInt > bInt {
		return aInt
	}
	return bInt
}

// minInt returns the smaller of two integers
func minInt(a, b interface{}) int64 {
	aInt, aOk := convertToInt64(a)
	bInt, bOk := convertToInt64(b)
	if !aOk || !bOk {
		return 0
	}
	if aInt < bInt {
		return aInt
	}
	return bInt
}

// multiply multiplies two numbers of any type (int, float, etc.)
func multiply(a, b interface{}) float64 {
	aFloat, aOk := convertToFloat64(a)
	bFloat, bOk := convertToFloat64(b)
	if !aOk || !bOk {
		return 0
	}
	return aFloat * bFloat
}
func floatEq(a, b interface{}) bool {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	return aOk && bOk && fa == fb
}
func floatGt(a, b interface{}) bool {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	return aOk && bOk && fa > fb
}
func floatGte(a, b interface{}) bool {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	return aOk && bOk && fa >= fb
}
func floatLt(a, b interface{}) bool {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	return aOk && bOk && fa < fb
}
func floatLte(a, b interface{}) bool {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	return aOk && bOk && fa <= fb
}
func floatNe(a, b interface{}) bool {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	return aOk && bOk && fa != fb
}
func addFloat(a, b interface{}) float64 {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	if !aOk || !bOk {
		return 0
	}
	return fa + fb
}
func subtractFloat(a, b interface{}) float64 {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	if !aOk || !bOk {
		return 0
	}
	return fa - fb
}
func multiplyFloat(a, b interface{}) float64 {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	if !aOk || !bOk {
		return 0
	}
	return fa * fb
}
func divideFloat(a, b interface{}) float64 {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	if !aOk || !bOk || fb == 0 {
		return 0
	}
	return fa / fb
}
func modFloat(a, b interface{}) float64 {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	if !aOk || !bOk || fb == 0 {
		return 0
	}
	return math.Mod(fa, fb)
}
func absFloat(a interface{}) float64 {
	fa, aOk := convertToFloat64(a)
	if !aOk {
		return 0
	}
	return math.Abs(fa)
}
func maxFloat(a, b interface{}) float64 {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	if !aOk || !bOk {
		return 0
	}
	return math.Max(fa, fb)
}
func minFloat(a, b interface{}) float64 {
	fa, aOk := convertToFloat64(a)
	fb, bOk := convertToFloat64(b)
	if !aOk || !bOk {
		return 0
	}
	return math.Min(fa, fb)
}

func convertToInt64(v interface{}) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int8:
		return int64(val), true
	case int16:
		return int64(val), true
	case int32:
		return int64(val), true
	case int64:
		return val, true
	case uint:
		return int64(val), true
	case uint8:
		return int64(val), true
	case uint16:
		return int64(val), true
	case uint32:
		return int64(val), true
	case uint64:
		if val <= math.MaxInt64 {
			return int64(val), true
		}
		return 0, false
	case float32:
		if val == float32(int64(val)) {
			return int64(val), true
		}
		return 0, false
	case float64:
		if val == float64(int64(val)) {
			return int64(val), true
		}
		return 0, false
	default:
		return 0, false
	}
}

// Helper function to convert various numeric types to float64
func convertToFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case string:
		if parsed, err := strconv.ParseFloat(val, 64); err == nil {
			return parsed, true
		}
		return 0, false
	default:
		return 0, false
	}
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
	// dict is a template function that creates a dictionary (map) from a variadic list of key-value pairs.
	// It takes an even number of arguments, where even-indexed arguments are keys (strings) and
	// odd-indexed arguments are their corresponding values. This allows for easy dictionary creation
	// within Go HTML templates.
	//
	// Example usage in a template: {{ dict "name" "John" "age" 30 }}
	// Returns: map[string]interface{}{"name": "John", "age": 30}
	engine.AddFunc("dict", func(values ...interface{}) map[string]interface{} {
		dict := make(map[string]interface{})
		for i := 0; i < len(values); i += 2 {
			key := values[i].(string)
			dict[key] = values[i+1]
		}
		return dict
	})

	engine.AddFunc("checkType", checkType)
	// length is a template function that returns the length of various types of collections and strings.
	// It supports strings, slices, arrays, maps, and uses reflection for additional types.
	// If the input is nil or cannot be measured, it returns 0.
	// Supports types like string, []interface{}, []string, []int, map[string]interface{}, map[string]string,
	// and uses reflection to handle other slice, array, map, and string types.
	engine.AddFunc("length", func(value interface{}) int {
		if value == nil {
			return 0
		}

		switch v := value.(type) {
		case string:
			return len(v)
		case []interface{}:
			return len(v)
		case []string:
			return len(v)
		case []int:
			return len(v)
		case map[string]interface{}:
			return len(v)
		case map[string]string:
			return len(v)
		default:
			// Use reflection for other slice/map types
			rv := reflect.ValueOf(value)
			switch rv.Kind() {
			case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
				return rv.Len()
			default:
				return 0
			}
		}
	})
	// formatNum is a template function that formats a numeric value to a specified decimal precision.
	// It supports various numeric types (float64, float32, int, int64, int32) and attempts to parse strings.
	// If conversion fails, it returns a zero value with the specified number of decimal places.
	// The precision parameter determines the number of decimal places in the formatted output.
	engine.AddFunc("formatNum", func(number interface{}, precision int) string {
		var num float64

		switch v := number.(type) {
		case float64:
			num = v
		case float32:
			num = float64(v)
		case int:
			num = float64(v)
		case int64:
			num = float64(v)
		case int32:
			num = float64(v)
		default:
			// Try to convert string to float
			if str, ok := v.(string); ok {
				if parsed, err := strconv.ParseFloat(str, 64); err == nil {
					num = parsed
				} else {
					return "0" + strings.Repeat(".0", precision)
				}
			} else {
				return "0" + strings.Repeat(".0", precision)
			}
		}

		format := fmt.Sprintf("%%.%df", precision)
		return fmt.Sprintf(format, num)
	})

	engine.AddFunc("FormatNumWithComma", func(f interface{}, precision ...int) string {
		if len(precision) > 0 {
			return FormatFloatWithCommas(f, precision[0])
		}
		return FormatFloatWithCommas(f)
	})
	// iterate creates a sequence of integers for template loops
	// Usage: iterate n - generates 0 to n-1
	// Usage: iterate n step - generates 0 to n-1 with step size
	// Usage: iterate start end - generates start to end-1
	// Usage: iterate start end step - generates start to end-1 with step
	engine.AddFunc("iterate", func(args ...int) []int {
		if len(args) == 0 {
			return []int{}
		}

		var start, end, step int

		switch len(args) {
		case 1:
			// iterate n - generate 0 to n-1
			start = 0
			end = args[0]
			step = 1
		case 2:
			// iterate n step - generate 0 to n-1 with step
			// OR iterate start end - generate start to end-1
			if args[1] <= args[0] {
				// Assume it's iterate start end
				start = args[0]
				end = args[1]
				step = 1
			} else {
				// Assume it's iterate n step
				start = 0
				end = args[0]
				step = args[1]
			}
		case 3:
			// iterate start end step
			start = args[0]
			end = args[1]
			step = args[2]
		default:
			return []int{}
		}

		if step == 0 {
			step = 1
		}

		var result []int
		if step > 0 {
			for i := start; i < end; i += step {
				result = append(result, i)
			}
		} else {
			for i := start; i > end; i += step {
				result = append(result, i)
			}
		}

		return result
	})

	// iterateFilter creates a filtered sequence of integers
	// Usage in templates requires combining with other functions
	engine.AddFunc("iterateFilter", func(start, end, step int) []int {
		if step == 0 {
			step = 1
		}

		var result []int
		if step > 0 {
			for i := start; i < end; i += step {
				result = append(result, i)
			}
		} else {
			for i := start; i > end; i += step {
				result = append(result, i)
			}
		}

		return result
	})

	// iterateRange creates a range with start, end, and step
	engine.AddFunc("iterateRange", func(start, end, step int) []int {
		if step == 0 {
			return []int{}
		}

		var result []int
		if step > 0 {
			for i := start; i < end; i += step {
				result = append(result, i)
			}
		} else {
			for i := start; i > end; i += step {
				result = append(result, i)
			}
		}

		return result
	})

	// iterateEven generates even numbers in a range
	engine.AddFunc("iterateEven", func(start, end int) []int {
		var result []int
		for i := start; i < end; i++ {
			if i%2 == 0 {
				result = append(result, i)
			}
		}
		return result
	})

	// iterateOdd generates odd numbers in a range
	engine.AddFunc("iterateOdd", func(start, end int) []int {
		var result []int
		for i := start; i < end; i++ {
			if i%2 != 0 {
				result = append(result, i)
			}
		}
		return result
	})

	// iterateMultiple generates multiples of a number
	engine.AddFunc("iterateMultiple", func(multiple, count int) []int {
		var result []int
		for i := 1; i <= count; i++ {
			result = append(result, multiple*i)
		}
		return result
	})

	// Add rand functionality for generating random numbers in templates
	engine.AddFunc("rand", func() int {
		return rand.Intn(100) // Returns random number 0-99
	})
	engine.AddFunc("randRange", func(min, max int) int {
		if min >= max {
			return min
		}
		return rand.Intn(max-min) + min
	})
	engine.AddFunc("randFloat", func() float64 {
		return rand.Float64() // Returns random float 0.0-1.0
	})
	engine.AddFunc("randFloatRange", func(min, max float64) float64 {
		if min >= max {
			return min
		}
		return rand.Float64()*(max-min) + min
	})
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
	engine.AddFunc("uuid", func() string {
		return utils.UUIDv4()
	})
	engine.AddFunc("asInt", func(v interface{}) int64 {
		if result, ok := convertToInt64(v); ok {
			return result
		}
		return 0
	})
	engine.AddFunc("asFloat", func(v interface{}) float64 {
		if result, ok := convertToFloat64(v); ok {
			return result
		}
		return 0.0
	})
	engine.AddFunc("asString", func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	})
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
