package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSQLFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "reader_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name            string
		content         string
		filename        string
		expectedContent string
		expectError     bool
	}{
		{
			name:            "Valid SQL file",
			content:         "CREATE TABLE users (id BIGSERIAL, name VARCHAR(255));",
			filename:        "valid.sql",
			expectedContent: "CREATE TABLE users (id BIGSERIAL, name VARCHAR(255));",
			expectError:     false,
		},
		{
			name:            "Empty SQL file",
			content:         "",
			filename:        "empty.sql",
			expectedContent: "",
			expectError:     false,
		},
		{
			name: "Multi-line SQL file",
			content: `-- This is a comment
CREATE TABLE users (
  id BIGSERIAL NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
  id BIGSERIAL NOT NULL,
  user_id BIGINT NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  CONSTRAINT fk_posts_users FOREIGN KEY (user_id) REFERENCES users(id)
);`,
			filename: "multiline.sql",
			expectedContent: `-- This is a comment
CREATE TABLE users (
  id BIGSERIAL NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
  id BIGSERIAL NOT NULL,
  user_id BIGINT NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  CONSTRAINT fk_posts_users FOREIGN KEY (user_id) REFERENCES users(id)
);`,
			expectError: false,
		},
		{
			name:            "SQL file with UTF-8 characters",
			content:         "-- SQL with UTF-8: ñáéíóú\nCREATE TABLE test (name VARCHAR(255));",
			filename:        "utf8.sql",
			expectedContent: "-- SQL with UTF-8: ñáéíóú\nCREATE TABLE test (name VARCHAR(255));",
			expectError:     false,
		},
		{
			name:        "Non-existent file",
			filename:    "nonexistent.sql",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string

			if !tt.expectError {
				// Create the test file
				filePath = filepath.Join(tempDir, tt.filename)
				err := os.WriteFile(filePath, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			} else {
				// Use a non-existent file path
				filePath = filepath.Join(tempDir, tt.filename)
			}

			// Call the function under test
			result, err := ReadSQLFile(filePath)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("ReadSQLFile() expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("ReadSQLFile() unexpected error: %v", err)
				return
			}
			if tt.expectError {
				// For error cases, we expect the error to contain the filename
				if err != nil && !containsString(err.Error(), tt.filename) {
					t.Errorf("ReadSQLFile() error should contain filename %s, got: %v", tt.filename, err)
				}
				return
			}

			// Check content
			if result != tt.expectedContent {
				t.Errorf("ReadSQLFile() content mismatch.\nGot:\n%q\nWant:\n%q", result, tt.expectedContent)
			}
		})
	}
}

func TestReadSQLFile_ErrorHandling(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "reader_error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name          string
		setupFunc     func() string
		expectedError string
		expectError   bool
	}{
		{
			name: "Permission denied (directory instead of file)",
			setupFunc: func() string {
				dirPath := filepath.Join(tempDir, "directory")
				_ = os.Mkdir(dirPath, 0755)
				return dirPath
			},
			expectedError: "failed to read file",
			expectError:   true,
		},
		{
			name: "Non-existent directory",
			setupFunc: func() string {
				return filepath.Join(tempDir, "nonexistent", "file.sql")
			},
			expectedError: "failed to open file",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFunc()

			result, err := ReadSQLFile(filePath)

			if tt.expectError && err == nil {
				t.Errorf("ReadSQLFile() expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("ReadSQLFile() unexpected error: %v", err)
				return
			}
			if tt.expectError {
				if !containsString(err.Error(), tt.expectedError) {
					t.Errorf("ReadSQLFile() error should contain %q, got: %v", tt.expectedError, err)
				}
				if result != "" {
					t.Errorf("ReadSQLFile() should return empty string on error, got: %q", result)
				}
				return
			}
		})
	}
}

func TestReadSQLFile_LargeFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "reader_large_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a large SQL content (simulate a large schema file)
	var largeContent string
	for i := 0; i < 1000; i++ {
		largeContent += "CREATE TABLE table_" + string(rune('a'+i%26)) + " (id BIGSERIAL, name VARCHAR(255));\n"
	}

	filePath := filepath.Join(tempDir, "large.sql")
	err = os.WriteFile(filePath, []byte(largeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	result, err := ReadSQLFile(filePath)
	if err != nil {
		t.Errorf("ReadSQLFile() unexpected error with large file: %v", err)
		return
	}

	if result != largeContent {
		t.Errorf("ReadSQLFile() large file content mismatch. Length got: %d, want: %d", len(result), len(largeContent))
	}
}

func TestReadSQLFile_EmptyFilename(t *testing.T) {
	result, err := ReadSQLFile("")

	if err == nil {
		t.Errorf("ReadSQLFile() with empty filename should return error")
		return
	}

	if result != "" {
		t.Errorf("ReadSQLFile() should return empty string on error, got: %q", result)
	}

	if !containsString(err.Error(), "failed to open file") {
		t.Errorf("ReadSQLFile() error should contain 'failed to open file', got: %v", err)
	}
}

// Helper function for string containment check
func containsString(haystack, needle string) bool {
	return len(haystack) >= len(needle) && haystack != needle &&
		(haystack[:len(needle)] == needle ||
			haystack[len(haystack)-len(needle):] == needle ||
			containsSubstring(haystack, needle))
}

func containsSubstring(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
