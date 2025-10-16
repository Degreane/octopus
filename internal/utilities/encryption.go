package utilities

import (
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	lua "github.com/yuin/gopher-lua"
)

// getChars returns a slice of runes containing ASCII letters, digits, and punctuation
func GetChars() []rune {
	var chars []rune
	for c := 'A'; c <= 'Z'; c++ {
		chars = append(chars, c)
	}
	for c := 'a'; c <= 'z'; c++ {
		chars = append(chars, c)
	}
	for c := '0'; c <= '9'; c++ {
		chars = append(chars, c)
	}
	chars = append(chars, []rune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ ")...)
	return chars
}

/*
// getChars returns a slice of runes containing ASCII letters, digits, and punctuation
func getChars() []rune {
	var chars []rune
	for c := 'A'; c <= 'Z'; c++ {
		chars = append(chars, c)
	}
	for c := 'a'; c <= 'z'; c++ {
		chars = append(chars, c)
	}
	for c := '0'; c <= '9'; c++ {
		chars = append(chars, c)
	}
	chars = append(chars, []rune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ ")...)
	return chars
}
*/

// StringToRuneSlice converts a string to a slice of runes
func StringToRuneSlice(s string) []rune {
	return []rune(s)
}

// RuneSliceToString converts a slice of runes back to a string
func RuneSliceToString(rs []rune) string {
	return string(rs)
}

// FindRuneIndex finds the index of a rune in a rune slice
func FindRuneIndex(rs []rune, r rune) int {
	for i, v := range rs {
		if v == r {
			return i
		}
	}
	return -1
}

// Encrypt encrypts a string using character substitution
func Encrypt(s1 string, r1 []rune, r2 []rune) []rune {
	// println("s1:", s1)
	// println("r1:", string(r1))
	// println("r2:", string(r2))
	indices := make([]rune, len(s1))
	for i, char := range s1 {
		index := FindRuneIndex(r1, rune(char))
		if index != -1 {
			indices[i] = r2[index]
		} else {
			indices[i] = rune(char)
		}
	}
	// println("indices:", string(indices))
	return indices
}

// Decrypt decrypts a string encrypted with Encrypt
func Decrypt(s1 string, r1 []rune, r2 []rune) []rune {
	// println("Decrypt s1 :", s1)
	// println("Decrypt r1 :", string(r1))
	// println("Decrypt r2 :", string(r2))
	indices := make([]rune, len(s1))
	for i, char := range s1 {
		index := FindRuneIndex(r2, rune(char))
		if index != -1 {
			indices[i] = r1[index]
		} else {
			indices[i] = rune(char)
		}
	}
	// println("Decrypt indices :", string(indices))
	return indices
}

/*
func Decrypt(s1 string, r1 []rune, r2 []rune) []rune {
	indices := make([]rune, len(s1))

	for i, char := range s1 {
		index := FindRuneIndex(r2, rune(char))
		if index != -1 {
			indices[i] = r1[index]
		} else {
			indices[i] = rune(char)
		}
	}
	return indices
}
*/

// ShuffleRuneSlice shuffles a rune slice based on a seed
func ShuffleRuneSlice(rs []rune, seed int64) []rune {
	print(seed)
	shuffled := make([]rune, len(rs))
	copy(shuffled, rs)
	source := rand.NewSource(seed)
	r := rand.New(source)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

/*
// ShuffleRuneSlice shuffles a rune slice based on a given seed
func ShuffleRuneSlice(rs []rune, seed int64) []rune {
	// Create a copy of the original slice to avoid modifying it
	shuffled := make([]rune, len(rs))
	copy(shuffled, rs)

	// Create a new random source with the given seed
	source := rand.NewSource(seed)
	r := rand.New(source)

	// Perform Fisher-Yates shuffle
	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

*/

// GetUnixGMTTimestampNearestMinute gets the current Unix timestamp rounded to the nearest minute in GMT
func GetUnixGMTTimestampNearestMinute() int64 {
	// return 1750260960
	now := time.Now().UTC()
	rounded := now.Truncate(time.Minute)
	log.Println(rounded.Unix())
	return rounded.Unix()
}

// GetUnixGMTTimestampPreviousMinute returns the Unix timestamp in GMT for the previous minute.
// It subtracts one minute from the current time, truncates it to remove seconds and nanoseconds,
// and then returns the Unix timestamp.
//
// Returns:
//   - int64: The Unix timestamp in seconds for the previous minute.
func GetUnixGMTTimestampPreviousMinute() int64 {
	now := time.Now().UTC()
	previousMinute := now.Add(-time.Minute).Truncate(time.Minute)
	return previousMinute.Unix()
}

// GetUnixGMTTimeStamp returns a Unix timestamp in GMT/UTC based on the provided minute parameter.
// The function handles different time intervals and returns timestamps aligned to specific boundaries.
//
// Parameters:
//   - minute: optional int - Number of minutes to process
//   - If no minute provided or minute <= 0: returns current timestamp truncated to minute
//   - If minute = 60: returns timestamp for start of current hour
//   - If minute = 120: returns timestamp for start of previous hour
//   - For other positive values: returns timestamp that many minutes in the past
//
// Returns:
//   - int64 - Unix timestamp in seconds
//
// Example:
//
//	current := GetUnixGMTTimeStamp()           // Current minute
//	fiveMinutesAgo := GetUnixGMTTimeStamp(5)   // 5 minutes ago
//	currentHour := GetUnixGMTTimeStamp(60)     // Start of current hour
//	previousHour := GetUnixGMTTimeStamp(120)   // Start of previous hour
func GetUnixGMTTimeStamp(minute ...int) int64 {
	if len(minute) > 0 {
		switch minute[0] {
		case -2:
			log.Println("GetUnixGMTTimeStamp: minute = -2")
			return GetUnixGMTTimestampPreviousMinute()
		case -1:
			log.Println("GetUnixGMTTimeStamp: minute = -1")
			return GetUnixGMTTimestampNearestMinute()
		case 60:
			log.Print("GetUnixGMTTimeStamp: minute = 60")
			now := time.Now().UTC()
			currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
			return currentHour.Unix()
		case 120:
			log.Print("GetUnixGMTTimeStamp: minute = 120")
			now := time.Now().UTC()
			previousHour := now.Add(-time.Hour).Truncate(time.Hour)
			return previousHour.Unix()
		case 1440:
			now := time.Now().UTC()
			startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			return startOfDay.Unix()
		case 10080:
			now := time.Now().UTC()
			startOfWeek := now.AddDate(0, 0, -int(now.Weekday())).Truncate(24 * time.Hour)
			return startOfWeek.Unix()
		case 43200:
			now := time.Now().UTC()
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
			return startOfMonth.Unix()
		case 525600:
			now := time.Now().UTC()
			startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			return startOfYear.Unix()
		default:
			if minute[0] > 0 {
				now := time.Now().UTC()
				previousMinute := now.Add(-time.Duration(minute[0]) * time.Minute).Truncate(time.Minute)
				return previousMinute.Unix()
			}
		}
	}
	return GetUnixGMTTimestampNearestMinute()
}

/*
chars := utilities.GetChars()
shuffledChars := utilities.ShuffleRuneSlice(chars, utilities.GetUnixGMTTimestampNearestMinute())
decryptedData := string(utilities.Decrypt(encryptedData, chars, shuffledChars)
*/

// DecryptData decrypts an encrypted string using character shuffling based on GMT timestamp.
// It uses the standard character set and shuffles it according to the nearest minute timestamp
// to maintain consistency across decryption attempts within the same minute.
//
// Parameters:
//   - dataToDecrypt: string - The encrypted data to be decrypted
//
// Returns:
//   - string - The decrypted data
//   - minute: optional int - Minutes to look back for timestamp
//   - 60: current hour
//   - 120: previous hour
//   - 1440: start of day
//   - 10080: start of week
//   - 43200: start of month
//   - 525600: start of year
//   - other positive values: that many minutes ago
//
// Example:
//
//	decrypted := DecryptData("encrypted_string_here")
//	fmt.Println(decrypted) // prints the decrypted result
func DecryptData(dataToDecrypt string, minute ...int) string {
	chars := GetChars()
	shuffledChars := ShuffleRuneSlice(chars, GetUnixGMTTimeStamp(minute...))
	decryptedData := string(Decrypt(dataToDecrypt, shuffledChars, chars))
	return decryptedData
}

// EncryptData encrypts a string using character shuffling based on GMT timestamp.
// It uses the standard character set and shuffles it according to the nearest minute timestamp
// to maintain consistency across encryption attempts within the same minute.
//
// Parameters:
//   - dataToEncrypt: string - The data to be encrypted
//   - minute: optional int - Minutes to look back for timestamp
//   - 60: current hour
//   - 120: previous hour
//   - 1440: start of day
//   - 10080: start of week
//   - 43200: start of month
//   - 525600: start of year
//   - other positive values: that many minutes ago
//
// Returns:
//   - string - The encrypted data
//
// Example:
//
//	encrypted := EncryptData("secret_data")
//	encryptedHourly := EncryptData("secret_data", 60)
func EncryptData(dataToEncrypt string, minute ...int) string {
	chars := GetChars()
	shuffledChars := ShuffleRuneSlice(chars, GetUnixGMTTimeStamp(minute...))
	encryptedData := string(Encrypt(dataToEncrypt, shuffledChars, chars))
	return encryptedData
}

// GetDecryptData returns a Lua function that decrypts encrypted data.
// The function is exposed to Lua scripts and provides access to the decryption functionality.
//
// Parameters passed from Lua:
//   - data: string - The encrypted data to decrypt
//   - minute: number (optional) - Minutes to look back for timestamp
//
// Returns to Lua:
//   - string - The decrypted data
//
// Usage in Lua:
//
//	local decrypted = decryptData("encrypted_string_here")
//	-- or with minute parameter
//	local decrypted = decryptData("encrypted_string_here", 5)
func GetDecryptData(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		data := L.ToString(1)
		minute := L.OptInt(2, 0) // Optional minute parameter, defaults to 0
		decrypted := DecryptData(data, minute)
		L.Push(lua.LString(decrypted))
		return 1
	}
}

// GetEncryptData returns a Lua function that encrypts data.
// The function is exposed to Lua scripts and provides access to the encryption functionality.
//
// Parameters passed from Lua:
//   - data: string - The data to encrypt
//   - minute: number (optional) - Minutes to look back for timestamp
//
// Returns to Lua:
//   - string - The encrypted data
//
// Usage in Lua:
//
//	local encrypted = encryptData("data_to_encrypt")
//	-- or with minute parameter
//	local encrypted = encryptData("data_to_encrypt", 60)
func GetEncryptData(c *fiber.Ctx) lua.LGFunction {
	return func(L *lua.LState) int {
		data := L.ToString(1)
		// println("Encryping Data 1 ", data)
		minute := L.OptInt(2, 0) // Optional minute parameter, defaults to 0
		encrypted := EncryptData(data, minute)
		// println("Encryping Data 2 ", encrypted)
		L.Push(lua.LString(encrypted))
		return 1
	}
}
