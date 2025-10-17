package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// ReadYamlAsJson reads a YAML file and returns its content as a JSON string.
//
// Parameters:
//   - filePath: path to the YAML file (relative or absolute)
//
// Returns:
//   - string: JSON representation of the YAML content
//   - error: error if file reading or conversion fails
func ReadYamlAsJson(filePath string) (string, error) {
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
