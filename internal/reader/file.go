// Package reader provides file reading utilities for SQL files.
//
// This package handles the reading and basic validation of SQL files
// that will be converted to Drizzle ORM schema definitions.
package reader

import (
	"fmt"
	"io"
	"os"
)

// ReadSQLFile reads the content of a SQL file and returns it as a string.
//
// This function opens the specified file, reads its entire content into memory,
// and returns it as a string. It includes proper error handling for file
// operations and uses wrapped errors for better error reporting.
//
// Parameters:
//   - filename: The path to the SQL file to read. Can be relative or absolute.
//
// Returns:
//   - string: The complete content of the SQL file
//   - error: An error if the file cannot be opened or read
//
// Example usage:
//
//	content, err := reader.ReadSQLFile("./schema.sql")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(content)
//
// Error handling:
//   - Returns wrapped errors for better debugging
//   - Distinguishes between file opening errors and reading errors
//   - Automatically closes the file using defer
func ReadSQLFile(filename string) (string, error) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		// Wrap the error with context about which file failed to open
		return "", fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	// Ensure the file is closed when the function returns
	defer file.Close()

	// Read the entire file content into memory
	content, err := io.ReadAll(file)
	if err != nil {
		// Wrap the error with context about which file failed to read
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Convert byte slice to string and return
	return string(content), nil
}
