package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/konojunya/sql-to-drizzle-schema/internal/parser"
)

func TestDefaultGeneratorOptions(t *testing.T) {
	options := DefaultGeneratorOptions()

	if options.TableNameCase != CamelCase {
		t.Errorf("DefaultGeneratorOptions() TableNameCase = %v, want %v", options.TableNameCase, CamelCase)
	}
	if options.ColumnNameCase != CamelCase {
		t.Errorf("DefaultGeneratorOptions() ColumnNameCase = %v, want %v", options.ColumnNameCase, CamelCase)
	}
	if options.IncludeComments != true {
		t.Errorf("DefaultGeneratorOptions() IncludeComments = %v, want %v", options.IncludeComments, true)
	}
	if options.ExportPrefix != "" {
		t.Errorf("DefaultGeneratorOptions() ExportPrefix = %v, want %v", options.ExportPrefix, "")
	}
	if options.IndentSize != 2 {
		t.Errorf("DefaultGeneratorOptions() IndentSize = %v, want %v", options.IndentSize, 2)
	}
}

func TestNewSchemaGenerator(t *testing.T) {
	tests := []struct {
		name        string
		dialect     parser.DatabaseDialect
		expectError bool
	}{
		{
			name:        "PostgreSQL generator",
			dialect:     parser.PostgreSQL,
			expectError: false,
		},
		{
			name:        "MySQL generator (unsupported)",
			dialect:     parser.MySQL,
			expectError: true,
		},
		{
			name:        "Spanner generator (unsupported)",
			dialect:     parser.Spanner,
			expectError: true,
		},
		{
			name:        "Invalid dialect",
			dialect:     parser.DatabaseDialect("invalid"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator, err := NewSchemaGenerator(tt.dialect)

			if tt.expectError && err == nil {
				t.Errorf("NewSchemaGenerator() expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("NewSchemaGenerator() unexpected error: %v", err)
				return
			}
			if tt.expectError {
				return
			}

			if generator == nil {
				t.Errorf("NewSchemaGenerator() returned nil generator")
				return
			}

			if generator.SupportedDialect() != tt.dialect {
				t.Errorf("NewSchemaGenerator() SupportedDialect() = %v, want %v", generator.SupportedDialect(), tt.dialect)
			}
		})
	}
}

func TestWriteSchemaToFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "generator_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		content     string
		filename    string
		expectError bool
	}{
		{
			name:        "Valid schema write",
			content:     "export const usersTable = pgTable('users', {});",
			filename:    filepath.Join(tempDir, "test.ts"),
			expectError: false,
		},
		{
			name:        "Empty content",
			content:     "",
			filename:    filepath.Join(tempDir, "empty.ts"),
			expectError: false,
		},
		{
			name:        "Invalid directory",
			content:     "content",
			filename:    "/nonexistent/dir/file.ts",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteSchemaToFile(tt.content, tt.filename)

			if tt.expectError && err == nil {
				t.Errorf("WriteSchemaToFile() expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("WriteSchemaToFile() unexpected error: %v", err)
				return
			}
			if tt.expectError {
				return
			}

			// Verify file was created and has correct content
			if _, err := os.Stat(tt.filename); os.IsNotExist(err) {
				t.Errorf("WriteSchemaToFile() file was not created: %s", tt.filename)
				return
			}

			content, err := os.ReadFile(tt.filename)
			if err != nil {
				t.Errorf("WriteSchemaToFile() failed to read written file: %v", err)
				return
			}

			if string(content) != tt.content {
				t.Errorf("WriteSchemaToFile() content = %v, want %v", string(content), tt.content)
			}
		})
	}
}

func TestGenerateSchemaToFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "generator_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test table data
	tables := []parser.Table{
		{
			Name: "users",
			Columns: []parser.Column{
				{
					Name:    "id",
					Type:    "BIGSERIAL",
					NotNull: true,
				},
				{
					Name:    "name",
					Type:    "VARCHAR",
					Length:  intPtr(255),
					NotNull: true,
				},
			},
			PrimaryKey: []string{"id"},
		},
	}

	outputFile := filepath.Join(tempDir, "schema.ts")
	options := DefaultGeneratorOptions()

	tests := []struct {
		name        string
		tables      []parser.Table
		dialect     parser.DatabaseDialect
		outputFile  string
		expectError bool
	}{
		{
			name:        "Valid PostgreSQL generation",
			tables:      tables,
			dialect:     parser.PostgreSQL,
			outputFile:  outputFile,
			expectError: false,
		},
		{
			name:        "Unsupported dialect",
			tables:      tables,
			dialect:     parser.MySQL,
			outputFile:  outputFile,
			expectError: true,
		},
		{
			name:        "Invalid output file",
			tables:      tables,
			dialect:     parser.PostgreSQL,
			outputFile:  "/nonexistent/dir/schema.ts",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateSchemaToFile(tt.tables, tt.dialect, tt.outputFile, options)

			if tt.expectError && err == nil {
				t.Errorf("GenerateSchemaToFile() expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("GenerateSchemaToFile() unexpected error: %v", err)
				return
			}
			if tt.expectError {
				return
			}

			// Verify file was created
			if _, err := os.Stat(tt.outputFile); os.IsNotExist(err) {
				t.Errorf("GenerateSchemaToFile() file was not created: %s", tt.outputFile)
				return
			}

			// Verify file has content
			content, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Errorf("GenerateSchemaToFile() failed to read generated file: %v", err)
				return
			}

			if len(content) == 0 {
				t.Errorf("GenerateSchemaToFile() generated empty file")
			}

			// Basic validation of generated content
			contentStr := string(content)
			if !containsString(contentStr, "import") {
				t.Errorf("GenerateSchemaToFile() generated content missing import statement")
			}
			if !containsString(contentStr, "pgTable") {
				t.Errorf("GenerateSchemaToFile() generated content missing pgTable")
			}
			if !containsString(contentStr, "users") {
				t.Errorf("GenerateSchemaToFile() generated content missing users table")
			}
		})
	}
}

func TestNamingCase(t *testing.T) {
	tests := []struct {
		caseType NamingCase
		expected string
	}{
		{CamelCase, "camel"},
		{PascalCase, "pascal"},
		{SnakeCase, "snake"},
		{KebabCase, "kebab"},
	}

	for _, tt := range tests {
		t.Run(string(tt.caseType), func(t *testing.T) {
			if string(tt.caseType) != tt.expected {
				t.Errorf("NamingCase string = %v, want %v", string(tt.caseType), tt.expected)
			}
		})
	}
}

// Helper functions for tests
func intPtr(i int) *int {
	return &i
}

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
