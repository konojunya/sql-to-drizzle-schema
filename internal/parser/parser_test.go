package parser

import (
	"testing"
)

func TestDefaultParseOptions(t *testing.T) {
	options := DefaultParseOptions()

	if options.Dialect != PostgreSQL {
		t.Errorf("DefaultParseOptions() Dialect = %v, want %v", options.Dialect, PostgreSQL)
	}
	if options.StrictMode != false {
		t.Errorf("DefaultParseOptions() StrictMode = %v, want %v", options.StrictMode, false)
	}
	if options.IgnoreUnsupported != true {
		t.Errorf("DefaultParseOptions() IgnoreUnsupported = %v, want %v", options.IgnoreUnsupported, true)
	}
}

func TestNewParser(t *testing.T) {
	tests := []struct {
		name         string
		dialect      DatabaseDialect
		expectedType string
		expectError  bool
	}{
		{
			name:         "PostgreSQL parser",
			dialect:      PostgreSQL,
			expectedType: "*parser.PostgreSQLParser",
			expectError:  false,
		},
		{
			name:         "MySQL parser (unsupported)",
			dialect:      MySQL,
			expectedType: "",
			expectError:  true,
		},
		{
			name:         "Spanner parser (unsupported)",
			dialect:      Spanner,
			expectedType: "",
			expectError:  true,
		},
		{
			name:         "Invalid dialect",
			dialect:      DatabaseDialect("invalid"),
			expectedType: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.dialect)

			if tt.expectError && err == nil {
				t.Errorf("NewParser() expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("NewParser() unexpected error: %v", err)
				return
			}
			if tt.expectError {
				return
			}

			if parser == nil {
				t.Errorf("NewParser() returned nil parser")
				return
			}

			if parser.SupportedDialect() != tt.dialect {
				t.Errorf("NewParser() SupportedDialect() = %v, want %v", parser.SupportedDialect(), tt.dialect)
			}
		})
	}
}

func TestParseSQLContent(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		dialect        DatabaseDialect
		expectedTables int
		expectedErrors int
		expectError    bool
	}{
		{
			name: "Valid PostgreSQL content",
			content: `CREATE TABLE users (
				id BIGSERIAL NOT NULL,
				name VARCHAR(255) NOT NULL
			);`,
			dialect:        PostgreSQL,
			expectedTables: 1,
			expectedErrors: 0,
			expectError:    false,
		},
		{
			name:           "Empty content",
			content:        "",
			dialect:        PostgreSQL,
			expectedTables: 0,
			expectedErrors: 0,
			expectError:    false,
		},
		{
			name:        "Unsupported dialect",
			content:     "CREATE TABLE test (id INT);",
			dialect:     MySQL,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := DefaultParseOptions()
			options.Dialect = tt.dialect

			result, err := ParseSQLContent(tt.content, tt.dialect, options)

			if tt.expectError && err == nil {
				t.Errorf("ParseSQLContent() expected error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("ParseSQLContent() unexpected error: %v", err)
				return
			}
			if tt.expectError {
				return
			}

			if len(result.Tables) != tt.expectedTables {
				t.Errorf("ParseSQLContent() tables count = %v, want %v", len(result.Tables), tt.expectedTables)
			}
			if len(result.Errors) != tt.expectedErrors {
				t.Errorf("ParseSQLContent() errors count = %v, want %v", len(result.Errors), tt.expectedErrors)
			}
			if result.Dialect != tt.dialect {
				t.Errorf("ParseSQLContent() dialect = %v, want %v", result.Dialect, tt.dialect)
			}
		})
	}
}

func TestDatabaseDialectString(t *testing.T) {
	tests := []struct {
		dialect  DatabaseDialect
		expected string
	}{
		{PostgreSQL, "postgresql"},
		{MySQL, "mysql"},
		{Spanner, "spanner"},
	}

	for _, tt := range tests {
		t.Run(string(tt.dialect), func(t *testing.T) {
			if string(tt.dialect) != tt.expected {
				t.Errorf("DatabaseDialect string = %v, want %v", string(tt.dialect), tt.expected)
			}
		})
	}
}
