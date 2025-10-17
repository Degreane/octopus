package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v3"
)

// readYamlAsJson reads a YAML file and returns its content as a JSON string.
//
// Parameters:
//   - filePath: path to the YAML file (relative or absolute)
//
// Returns:
//   - string: JSON representation of the YAML content
//   - error: error if file reading or conversion fails
func readYamlAsJson(filePath string) (string, error) {
	// Open the YAML file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open YAML file %s: %w", filePath, err)
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read YAML file %s: %w", filePath, err)
	}

	// Parse YAML content
	var yamlData interface{}
	err = yaml.Unmarshal(content, &yamlData)
	if err != nil {
		return "", fmt.Errorf("failed to parse YAML content from %s: %w", filePath, err)
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(yamlData)
	if err != nil {
		return "", fmt.Errorf("failed to convert YAML to JSON for %s: %w", filePath, err)
	}

	return string(jsonBytes), nil
}

// ReadYamlFileLua exposes eocto.readYamlFile(path) to Lua.
// It reads a YAML file from the given path and returns its contents as a JSON string.
// On error or invalid input, it returns nil.
// Usage in Lua:
//   local jsonStr = eocto.readYamlFile("/path/to/config.yaml")
func ReadYamlFileLua(L *lua.LState) int {
	path := L.OptString(1, "")
	if path == "" {
		L.Push(lua.LNil)
		return 1
	}
	jsonStr, err := readYamlAsJson(path)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(jsonStr))
	return 1
}
