package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/konojunya/sql-to-drizzle-schema/internal/generator"
	"github.com/konojunya/sql-to-drizzle-schema/internal/parser"
	"github.com/konojunya/sql-to-drizzle-schema/internal/reader"
)

// TestEndToEndConversion tests the complete conversion pipeline
func TestEndToEndConversion(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name             string
		sqlContent       string
		expectedTables   []string
		expectedImports  []string
		expectedFeatures []string
		expectError      bool
	}{
		{
			name: "Simple table with basic columns",
			sqlContent: `CREATE TABLE users (
				id BIGSERIAL NOT NULL,
				name VARCHAR(255) NOT NULL,
				email VARCHAR(255) UNIQUE,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT pk_users PRIMARY KEY (id)
			);`,
			expectedTables:   []string{"users"},
			expectedImports:  []string{"bigserial", "varchar", "timestamp", "pgTable"},
			expectedFeatures: []string{"notNull()", "unique()", "primaryKey()", "defaultNow()"},
			expectError:      false,
		},
		{
			name: "Multiple tables with foreign keys",
			sqlContent: `CREATE TABLE users (
				id BIGSERIAL NOT NULL,
				name VARCHAR(255) NOT NULL,
				CONSTRAINT pk_users PRIMARY KEY (id)
			);
			
			CREATE TABLE posts (
				id BIGSERIAL NOT NULL,
				title VARCHAR(255) NOT NULL,
				content TEXT,
				user_id BIGINT NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT pk_posts PRIMARY KEY (id),
				CONSTRAINT fk_posts_users FOREIGN KEY (user_id) REFERENCES users(id)
			);
			
			CREATE TABLE comments (
				id BIGSERIAL NOT NULL,
				content TEXT NOT NULL,
				post_id BIGINT NOT NULL,
				user_id BIGINT NOT NULL,
				CONSTRAINT pk_comments PRIMARY KEY (id),
				CONSTRAINT fk_comments_posts FOREIGN KEY (post_id) REFERENCES posts(id),
				CONSTRAINT fk_comments_users FOREIGN KEY (user_id) REFERENCES users(id)
			);`,
			expectedTables:   []string{"users", "posts", "comments"},
			expectedImports:  []string{"bigserial", "varchar", "text", "bigint", "timestamp", "pgTable"},
			expectedFeatures: []string{"references(() => usersTable.id)", "references(() => postsTable.id)"},
			expectError:      false,
		},
		{
			name: "Table with various data types",
			sqlContent: `CREATE TABLE test_types (
				id BIGSERIAL NOT NULL,
				varchar_col VARCHAR(255) NOT NULL,
				text_col TEXT,
				int_col INTEGER,
				bigint_col BIGINT,
				decimal_col DECIMAL(10,2),
				boolean_col BOOLEAN DEFAULT FALSE,
				timestamp_col TIMESTAMP WITH TIME ZONE,
				date_col DATE,
				real_col REAL,
				double_col DOUBLE PRECISION,
				CONSTRAINT pk_test_types PRIMARY KEY (id)
			);`,
			expectedTables:   []string{"testTypes"},
			expectedImports:  []string{"bigserial", "varchar", "text", "integer", "bigint", "decimal", "boolean", "timestamp", "date", "real", "doublePrecision", "pgTable"},
			expectedFeatures: []string{"default(false)", "precision: 10, scale: 2"},
			expectError:      false,
		},
		{
			name: "Empty SQL file",
			sqlContent: `-- Just comments
			-- Another comment`,
			expectedTables: []string{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create SQL input file
			sqlFile := filepath.Join(tempDir, tt.name+"_input.sql")
			err := os.WriteFile(sqlFile, []byte(tt.sqlContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create SQL file: %v", err)
			}

			// Step 1: Read SQL file
			content, err := reader.ReadSQLFile(sqlFile)
			if err != nil {
				if tt.expectError {
					return
				}
				t.Fatalf("Failed to read SQL file: %v", err)
			}

			// Step 2: Parse SQL content
			parseOptions := parser.DefaultParseOptions()
			parseOptions.Dialect = parser.PostgreSQL
			parseResult, err := parser.ParseSQLContent(content, parser.PostgreSQL, parseOptions)
			if err != nil {
				if tt.expectError {
					return
				}
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			// Validate parsing results
			if len(parseResult.Tables) != len(tt.expectedTables) {
				t.Errorf("Expected %d tables, got %d", len(tt.expectedTables), len(parseResult.Tables))
			}

			// Step 3: Generate Drizzle schema
			generatorOptions := generator.DefaultGeneratorOptions()
			outputFile := filepath.Join(tempDir, tt.name+"_output.ts")

			err = generator.GenerateSchemaToFile(parseResult.Tables, parser.PostgreSQL, outputFile, generatorOptions)
			if err != nil {
				if tt.expectError {
					return
				}
				t.Fatalf("Failed to generate schema: %v", err)
			}

			// Step 4: Validate generated output
			generatedContent, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			contentStr := string(generatedContent)

			// Check expected imports
			for _, expectedImport := range tt.expectedImports {
				if !strings.Contains(contentStr, expectedImport) {
					t.Errorf("Generated content missing expected import: %s", expectedImport)
				}
			}

			// Check expected table names (converted to camelCase)
			for _, expectedTable := range tt.expectedTables {
				if !strings.Contains(contentStr, expectedTable) {
					t.Errorf("Generated content missing expected table: %s", expectedTable)
				}
			}

			// Check expected features
			for _, expectedFeature := range tt.expectedFeatures {
				if !strings.Contains(contentStr, expectedFeature) {
					t.Errorf("Generated content missing expected feature: %s", expectedFeature)
				}
			}

			// Validate TypeScript syntax basics
			if len(tt.expectedTables) > 0 {
				if !strings.Contains(contentStr, "import {") {
					t.Error("Generated content should contain import statement")
				}
				if !strings.Contains(contentStr, "} from 'drizzle-orm/pg-core';") {
					t.Error("Generated content should import from drizzle-orm/pg-core")
				}
				if !strings.Contains(contentStr, "export const") {
					t.Error("Generated content should export table constants")
				}
				if !strings.Contains(contentStr, "pgTable(") {
					t.Error("Generated content should use pgTable function")
				}
			}
		})
	}
}

// TestTableDependencyOrdering tests that tables are generated in correct dependency order
func TestTableDependencyOrdering(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dependency_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// SQL with tables in reverse dependency order to test sorting
	sqlContent := `CREATE TABLE comments (
		id BIGSERIAL NOT NULL,
		content TEXT NOT NULL,
		post_id BIGINT NOT NULL,
		user_id BIGINT NOT NULL,
		CONSTRAINT fk_comments_posts FOREIGN KEY (post_id) REFERENCES posts(id),
		CONSTRAINT fk_comments_users FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE posts (
		id BIGSERIAL NOT NULL,
		title VARCHAR(255) NOT NULL,
		user_id BIGINT NOT NULL,
		CONSTRAINT fk_posts_users FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE users (
		id BIGSERIAL NOT NULL,
		name VARCHAR(255) NOT NULL
	);`

	sqlFile := filepath.Join(tempDir, "dependency_test.sql")
	err = os.WriteFile(sqlFile, []byte(sqlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create SQL file: %v", err)
	}

	outputFile := filepath.Join(tempDir, "dependency_output.ts")

	// Process the file
	content, err := reader.ReadSQLFile(sqlFile)
	if err != nil {
		t.Fatalf("Failed to read SQL file: %v", err)
	}

	parseOptions := parser.DefaultParseOptions()
	parseResult, err := parser.ParseSQLContent(content, parser.PostgreSQL, parseOptions)
	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}

	generatorOptions := generator.DefaultGeneratorOptions()
	err = generator.GenerateSchemaToFile(parseResult.Tables, parser.PostgreSQL, outputFile, generatorOptions)
	if err != nil {
		t.Fatalf("Failed to generate schema: %v", err)
	}

	// Read and validate output order
	generatedContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(generatedContent)

	// Find positions of each table definition
	usersPos := strings.Index(contentStr, "export const usersTable")
	postsPos := strings.Index(contentStr, "export const postsTable")
	commentsPos := strings.Index(contentStr, "export const commentsTable")

	if usersPos == -1 || postsPos == -1 || commentsPos == -1 {
		t.Fatal("Not all table definitions found in generated content")
	}

	// Verify correct order: users -> posts -> comments
	if !(usersPos < postsPos && postsPos < commentsPos) {
		t.Errorf("Tables not in correct dependency order. Got: users=%d, posts=%d, comments=%d", usersPos, postsPos, commentsPos)
	}

	// Verify foreign key references are correct
	if !strings.Contains(contentStr, "references(() => usersTable.id)") {
		t.Error("Missing foreign key reference to usersTable.id")
	}
	if !strings.Contains(contentStr, "references(() => postsTable.id)") {
		t.Error("Missing foreign key reference to postsTable.id")
	}
}

// TestNamingConventions tests different naming convention options
func TestNamingConventions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "naming_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sqlContent := `CREATE TABLE user_profiles (
		id BIGSERIAL NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		birth_date DATE
	);`

	tests := []struct {
		name           string
		tableCase      generator.NamingCase
		columnCase     generator.NamingCase
		expectedTable  string
		expectedColumn string
	}{
		{
			name:           "CamelCase for both",
			tableCase:      generator.CamelCase,
			columnCase:     generator.CamelCase,
			expectedTable:  "userProfiles",
			expectedColumn: "firstName",
		},
		{
			name:           "PascalCase for table, camelCase for columns",
			tableCase:      generator.PascalCase,
			columnCase:     generator.CamelCase,
			expectedTable:  "UserProfiles",
			expectedColumn: "firstName",
		},
		{
			name:           "SnakeCase for both",
			tableCase:      generator.SnakeCase,
			columnCase:     generator.SnakeCase,
			expectedTable:  "user_profiles",
			expectedColumn: "first_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlFile := filepath.Join(tempDir, tt.name+"_input.sql")
			err := os.WriteFile(sqlFile, []byte(sqlContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create SQL file: %v", err)
			}

			outputFile := filepath.Join(tempDir, tt.name+"_output.ts")

			// Process with specific naming options
			content, err := reader.ReadSQLFile(sqlFile)
			if err != nil {
				t.Fatalf("Failed to read SQL file: %v", err)
			}

			parseOptions := parser.DefaultParseOptions()
			parseResult, err := parser.ParseSQLContent(content, parser.PostgreSQL, parseOptions)
			if err != nil {
				t.Fatalf("Failed to parse SQL: %v", err)
			}

			generatorOptions := generator.DefaultGeneratorOptions()
			generatorOptions.TableNameCase = tt.tableCase
			generatorOptions.ColumnNameCase = tt.columnCase

			err = generator.GenerateSchemaToFile(parseResult.Tables, parser.PostgreSQL, outputFile, generatorOptions)
			if err != nil {
				t.Fatalf("Failed to generate schema: %v", err)
			}

			// Validate naming conventions
			generatedContent, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			contentStr := string(generatedContent)

			if !strings.Contains(contentStr, "export const "+tt.expectedTable+"Table") {
				t.Errorf("Expected table name %s not found in generated content", tt.expectedTable)
			}

			if !strings.Contains(contentStr, tt.expectedColumn+":") {
				t.Errorf("Expected column name %s not found in generated content", tt.expectedColumn)
			}
		})
	}
}

// TestErrorHandling tests various error conditions
func TestErrorHandling(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		sqlContent  string
		expectError bool
		errorType   string
	}{
		{
			name:        "Invalid SQL syntax",
			sqlContent:  "INVALID SQL SYNTAX",
			expectError: false, // Parser should handle this gracefully with IgnoreUnsupported
		},
		{
			name:        "Malformed CREATE TABLE",
			sqlContent:  "CREATE TABLE users",
			expectError: false, // Should be ignored
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlFile := filepath.Join(tempDir, tt.name+"_input.sql")
			err := os.WriteFile(sqlFile, []byte(tt.sqlContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create SQL file: %v", err)
			}

			outputFile := filepath.Join(tempDir, tt.name+"_output.ts")

			content, err := reader.ReadSQLFile(sqlFile)
			if err != nil {
				t.Fatalf("Failed to read SQL file: %v", err)
			}

			parseOptions := parser.DefaultParseOptions()
			parseResult, err := parser.ParseSQLContent(content, parser.PostgreSQL, parseOptions)

			if tt.expectError && err == nil {
				t.Errorf("Expected parsing error but got none")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected parsing error: %v", err)
				return
			}

			if !tt.expectError {
				generatorOptions := generator.DefaultGeneratorOptions()
				err = generator.GenerateSchemaToFile(parseResult.Tables, parser.PostgreSQL, outputFile, generatorOptions)

				// Generation should succeed even with empty tables
				if err != nil {
					t.Errorf("Unexpected generation error: %v", err)
				}
			}
		})
	}
}
