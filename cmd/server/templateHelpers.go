package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/utils"
)

// FormatFloatWithCommas formats a float with commas and specified precision
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

// isNumeric checks if string is a float or integer
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

// checkType checks if string is int return int if float return float
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

// convertToFloat64 converts various numeric types to float64
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

// dictHelper creates a map[string]interface{} from key-value pairs for template use
func dictHelper(values ...interface{}) map[string]interface{} {
	dict := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		dict[key] = values[i+1]
	}
	return dict
}

// length returns the length of various types of collections and strings
func length(value interface{}) int {
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
}

// formatNum formats a numeric value to a specified decimal precision
func formatNum(number interface{}, precision int) string {
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
}

// formatNumWithComma formats a number with commas and optional precision
func formatNumWithComma(f interface{}, precision ...int) string {
	if len(precision) > 0 {
		return FormatFloatWithCommas(f, precision[0])
	}
	return FormatFloatWithCommas(f)
}

// iterate creates a sequence of integers for template loops
// Usage: iterate n - generates 0 to n-1
// Usage: iterate n step - generates 0 to n-1 with step size
// Usage: iterate start end - generates start to end-1
// Usage: iterate start end step - generates start to end-1 with step
func iterate(args ...int) []int {
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
}

// iterateFilter creates a filtered sequence of integers
func iterateFilter(start, end, step int) []int {
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
}

// iterateRange creates a range with start, end, and step
func iterateRange(start, end, step int) []int {
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
}

// iterateEven generates even numbers in a range
func iterateEven(start, end int) []int {
	var result []int
	for i := start; i < end; i++ {
		if i%2 == 0 {
			result = append(result, i)
		}
	}
	return result
}

// iterateOdd generates odd numbers in a range
func iterateOdd(start, end int) []int {
	var result []int
	for i := start; i < end; i++ {
		if i%2 != 0 {
			result = append(result, i)
		}
	}
	return result
}

// iterateMultiple generates multiples of a number
func iterateMultiple(multiple, count int) []int {
	var result []int
	for i := 1; i <= count; i++ {
		result = append(result, multiple*i)
	}
	return result
}

// randHelper returns a random number between 0-99
func randHelper() int {
	return rand.Intn(100)
}

// randRange returns a random number between min and max
func randRange(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min) + min
}

// randFloat returns a random float between 0.0-1.0
func randFloat() float64 {
	return rand.Float64()
}

// randFloatRange returns a random float between min and max
func randFloatRange(min, max float64) float64 {
	if min >= max {
		return min
	}
	return rand.Float64()*(max-min) + min
}

// uuidHelper generates a UUID v4
func uuidHelper() string {
	return utils.UUIDv4()
}

// asInt converts a value to int64
func asInt(v interface{}) int64 {
	if result, ok := convertToInt64(v); ok {
		return result
	}
	return 0
}

// asFloat converts a value to float64
func asFloat(v interface{}) float64 {
	if result, ok := convertToFloat64(v); ok {
		return result
	}
	return 0.0
}

// asString converts a value to string
func asString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
