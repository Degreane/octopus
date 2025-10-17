package utilities

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	lua "github.com/yuin/gopher-lua"
)

// readCsvAsJson reads a CSV file and returns its content as a JSON string.
// Assumes the first row contains headers. Each subsequent row is converted
// into an object keyed by those headers.
//
// Parameters:
//   - filePath: path to the CSV file (relative or absolute)
//
// Returns:
//   - string: JSON representation of the CSV content (array of objects)
//   - error: error if file reading or conversion fails
func readCsvAsJson(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open CSV file %s: %w", filePath, err)
	}
	defer f.Close()

	// Wrap with a reader to allow BOM stripping
	br := bufio.NewReader(f)
	peek, _ := br.Peek(3)
	if len(peek) >= 3 && bytes.Equal(peek, []byte{0xEF, 0xBB, 0xBF}) {
		// discard BOM
		_, _ = br.Discard(3)
	}

	r := csv.NewReader(br)
	// Let csv.Reader handle variable records; we'll normalize later
	r.FieldsPerRecord = -1

	records, err := r.ReadAll()
	if err != nil {
		if err == io.EOF {
			return "[]", nil
		}
		return "", fmt.Errorf("failed to read CSV file %s: %w", filePath, err)
	}

	if len(records) == 0 {
		return "[]", nil
	}

	headers := records[0]
	rows := records[1:]

	// Build slice of maps
	out := make([]map[string]string, 0, len(rows))
	for _, rec := range rows {
		rowMap := make(map[string]string, len(headers))
		max := len(headers)
		for i := 0; i < max; i++ {
			var val string
			if i < len(rec) {
				val = rec[i]
			} else {
				val = ""
			}
			h := headers[i]
			rowMap[h] = val
		}
		out = append(out, rowMap)
	}

	jsonBytes, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("failed to convert CSV to JSON for %s: %w", filePath, err)
	}
	return string(jsonBytes), nil
}

// ReadCsvFileLua exposes eocto.readCsvFile(path) to Lua.
// It reads a CSV file and returns its contents as a JSON string (array of objects).
// On error or invalid input, it returns nil.
// Usage in Lua:
//   local jsonStr = eocto.readCsvFile("/path/to/file.csv")
func ReadCsvFileLua(L *lua.LState) int {
	path := L.OptString(1, "")
	if path == "" {
		L.Push(lua.LNil)
		return 1
	}
	jsonStr, err := readCsvAsJson(path)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(jsonStr))
	return 1
}
