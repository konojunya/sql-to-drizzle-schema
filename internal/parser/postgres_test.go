package parser

import (
	"testing"
)

func TestPostgreSQLParser_SupportedDialect(t *testing.T) {
	parser := NewPostgreSQLParser()
	if parser.SupportedDialect() != PostgreSQL {
		t.Errorf("Expected PostgreSQL dialect, got %v", parser.SupportedDialect())
	}
}

func TestPostgreSQLParser_isCreateTableStatement(t *testing.T) {
	parser := NewPostgreSQLParser()

	tests := []struct {
		name     string
		stmt     string
		expected bool
	}{
		{
			name:     "Valid CREATE TABLE",
			stmt:     "CREATE TABLE users (id INT);",
			expected: true,
		},
		{
			name:     "Case insensitive CREATE TABLE",
			stmt:     "create table users (id int);",
			expected: true,
		},
		{
			name:     "CREATE TABLE with whitespace",
			stmt:     "  CREATE   TABLE   users (id INT);",
			expected: true,
		},
		{
			name:     "Not a CREATE TABLE",
			stmt:     "SELECT * FROM users;",
			expected: false,
		},
		{
			name:     "CREATE INDEX",
			stmt:     "CREATE INDEX idx_users ON users (id);",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.isCreateTableStatement(tt.stmt)
			if result != tt.expected {
				t.Errorf("isCreateTableStatement() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPostgreSQLParser_parseColumnRegex(t *testing.T) {
	parser := NewPostgreSQLParser()
	options := ParseOptions{
		Dialect:           PostgreSQL,
		StrictMode:        false,
		IgnoreUnsupported: false,
	}

	tests := []struct {
		name      string
		columnDef string
		expected  Column
		wantErr   bool
	}{
		{
			name:      "Basic VARCHAR column",
			columnDef: "name VARCHAR(255)",
			expected: Column{
				Name:          "name",
				Type:          "VARCHAR",
				Length:        intPtr(255),
				NotNull:       false,
				Unique:        false,
				AutoIncrement: false,
			},
			wantErr: false,
		},
		{
			name:      "BIGINT with NOT NULL",
			columnDef: "id BIGINT NOT NULL",
			expected: Column{
				Name:          "id",
				Type:          "BIGINT",
				NotNull:       true,
				Unique:        false,
				AutoIncrement: false,
			},
			wantErr: false,
		},
		{
			name:      "BIGSERIAL (auto increment)",
			columnDef: "id BIGSERIAL NOT NULL",
			expected: Column{
				Name:          "id",
				Type:          "BIGSERIAL",
				NotNull:       true,
				Unique:        false,
				AutoIncrement: true,
			},
			wantErr: false,
		},
		{
			name:      "VARCHAR with UNIQUE constraint",
			columnDef: "email VARCHAR(255) NOT NULL UNIQUE",
			expected: Column{
				Name:          "email",
				Type:          "VARCHAR",
				Length:        intPtr(255),
				NotNull:       true,
				Unique:        true,
				AutoIncrement: false,
			},
			wantErr: false,
		},
		{
			name:      "VARCHAR with DEFAULT value",
			columnDef: "role VARCHAR(255) NOT NULL DEFAULT 'user'",
			expected: Column{
				Name:          "role",
				Type:          "VARCHAR",
				Length:        intPtr(255),
				NotNull:       true,
				Unique:        false,
				AutoIncrement: false,
				DefaultValue:  stringPtr("'user'"),
			},
			wantErr: false,
		},
		{
			name:      "TIMESTAMP WITH TIME ZONE",
			columnDef: "created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP",
			expected: Column{
				Name:          "created_at",
				Type:          "TIMESTAMP WITH TIME ZONE",
				NotNull:       true,
				Unique:        false,
				AutoIncrement: false,
				DefaultValue:  stringPtr("CURRENT_TIMESTAMP"),
			},
			wantErr: false,
		},
		{
			name:      "DECIMAL with precision and scale",
			columnDef: "price DECIMAL(10,2) NOT NULL",
			expected: Column{
				Name:          "price",
				Type:          "DECIMAL",
				Length:        intPtr(10),
				Scale:         intPtr(2),
				NotNull:       true,
				Unique:        false,
				AutoIncrement: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseColumnRegex(tt.columnDef, options)

			if tt.wantErr && err == nil {
				t.Errorf("parseColumnRegex() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("parseColumnRegex() unexpected error: %v", err)
				return
			}
			if tt.wantErr {
				return
			}

			if result.Name != tt.expected.Name {
				t.Errorf("parseColumnRegex() Name = %v, want %v", result.Name, tt.expected.Name)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("parseColumnRegex() Type = %v, want %v", result.Type, tt.expected.Type)
			}
			if !compareIntPtr(result.Length, tt.expected.Length) {
				t.Errorf("parseColumnRegex() Length = %v, want %v", result.Length, tt.expected.Length)
			}
			if !compareIntPtr(result.Scale, tt.expected.Scale) {
				t.Errorf("parseColumnRegex() Scale = %v, want %v", result.Scale, tt.expected.Scale)
			}
			if result.NotNull != tt.expected.NotNull {
				t.Errorf("parseColumnRegex() NotNull = %v, want %v", result.NotNull, tt.expected.NotNull)
			}
			if result.Unique != tt.expected.Unique {
				t.Errorf("parseColumnRegex() Unique = %v, want %v", result.Unique, tt.expected.Unique)
			}
			if result.AutoIncrement != tt.expected.AutoIncrement {
				t.Errorf("parseColumnRegex() AutoIncrement = %v, want %v", result.AutoIncrement, tt.expected.AutoIncrement)
			}
			if !compareStringPtr(result.DefaultValue, tt.expected.DefaultValue) {
				t.Errorf("parseColumnRegex() DefaultValue = %v, want %v", result.DefaultValue, tt.expected.DefaultValue)
			}
		})
	}
}

func TestPostgreSQLParser_ParseSQL(t *testing.T) {
	parser := NewPostgreSQLParser()
	options := ParseOptions{
		Dialect:           PostgreSQL,
		StrictMode:        false,
		IgnoreUnsupported: true,
	}

	tests := []struct {
		name           string
		sql            string
		expectedTables int
		expectedErrors int
	}{
		{
			name: "Single table with basic columns",
			sql: `CREATE TABLE users (
				id BIGSERIAL NOT NULL,
				name VARCHAR(255) NOT NULL,
				email VARCHAR(255) NOT NULL UNIQUE,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT pk_users PRIMARY KEY (id)
			);`,
			expectedTables: 1,
			expectedErrors: 0,
		},
		{
			name: "Multiple tables with foreign keys",
			sql: `CREATE TABLE users (
				id BIGSERIAL NOT NULL,
				name VARCHAR(255) NOT NULL,
				CONSTRAINT pk_users PRIMARY KEY (id)
			);
			
			CREATE TABLE posts (
				id BIGSERIAL NOT NULL,
				title VARCHAR(255) NOT NULL,
				user_id BIGINT NOT NULL,
				CONSTRAINT pk_posts PRIMARY KEY (id),
				CONSTRAINT fk_posts_users FOREIGN KEY (user_id) REFERENCES users(id)
			);`,
			expectedTables: 2,
			expectedErrors: 0,
		},
		{
			name: "Table with comments and empty lines",
			sql: `-- This is a comment
			CREATE TABLE users (
				-- User ID
				id BIGSERIAL NOT NULL,
				-- User name
				name VARCHAR(255) NOT NULL
			);`,
			expectedTables: 1,
			expectedErrors: 0,
		},
		{
			name:           "Empty SQL",
			sql:            "",
			expectedTables: 0,
			expectedErrors: 0,
		},
		{
			name:           "Only comments",
			sql:            "-- This is just a comment\n-- Another comment",
			expectedTables: 0,
			expectedErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseSQL(tt.sql, options)

			if err != nil {
				t.Errorf("ParseSQL() unexpected error: %v", err)
				return
			}

			if len(result.Tables) != tt.expectedTables {
				t.Errorf("ParseSQL() tables count = %v, want %v", len(result.Tables), tt.expectedTables)
			}

			if len(result.Errors) != tt.expectedErrors {
				t.Errorf("ParseSQL() errors count = %v, want %v", len(result.Errors), tt.expectedErrors)
			}

			if result.Dialect != PostgreSQL {
				t.Errorf("ParseSQL() dialect = %v, want %v", result.Dialect, PostgreSQL)
			}
		})
	}
}

func TestPostgreSQLParser_parseCreateTableRegex(t *testing.T) {
	parser := NewPostgreSQLParser()
	options := ParseOptions{
		Dialect:           PostgreSQL,
		StrictMode:        false,
		IgnoreUnsupported: true,
	}

	tests := []struct {
		name         string
		sql          string
		expectedName string
		expectedCols int
		expectedPK   []string
		expectedFKs  int
		wantErr      bool
	}{
		{
			name: "Basic table with primary key",
			sql: `CREATE TABLE users (
				id BIGSERIAL NOT NULL,
				name VARCHAR(255) NOT NULL,
				CONSTRAINT pk_users PRIMARY KEY (id)
			);`,
			expectedName: "users",
			expectedCols: 2,
			expectedPK:   []string{"id"},
			expectedFKs:  0,
			wantErr:      false,
		},
		{
			name: "Table with foreign key",
			sql: `CREATE TABLE posts (
				id BIGSERIAL NOT NULL,
				user_id BIGINT NOT NULL,
				CONSTRAINT pk_posts PRIMARY KEY (id),
				CONSTRAINT fk_posts_users FOREIGN KEY (user_id) REFERENCES users(id)
			);`,
			expectedName: "posts",
			expectedCols: 2,
			expectedPK:   []string{"id"},
			expectedFKs:  1,
			wantErr:      false,
		},
		{
			name: "Table with unique constraint",
			sql: `CREATE TABLE role_permissions (
				role_id BIGINT NOT NULL,
				permission_id BIGINT NOT NULL,
				CONSTRAINT unique_role_permission UNIQUE (role_id, permission_id)
			);`,
			expectedName: "role_permissions",
			expectedCols: 2,
			expectedPK:   []string{},
			expectedFKs:  0,
			wantErr:      false,
		},
		{
			name:    "Invalid table statement",
			sql:     "INVALID SQL STATEMENT",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseCreateTableRegex(tt.sql, options)

			if tt.wantErr && err == nil {
				t.Errorf("parseCreateTableRegex() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("parseCreateTableRegex() unexpected error: %v", err)
				return
			}
			if tt.wantErr {
				return
			}

			if result.Name != tt.expectedName {
				t.Errorf("parseCreateTableRegex() Name = %v, want %v", result.Name, tt.expectedName)
			}
			if len(result.Columns) != tt.expectedCols {
				t.Errorf("parseCreateTableRegex() Columns count = %v, want %v", len(result.Columns), tt.expectedCols)
			}
			if len(result.PrimaryKey) != len(tt.expectedPK) {
				t.Errorf("parseCreateTableRegex() PrimaryKey count = %v, want %v", len(result.PrimaryKey), len(tt.expectedPK))
			}
			for i, pk := range tt.expectedPK {
				if i < len(result.PrimaryKey) && result.PrimaryKey[i] != pk {
					t.Errorf("parseCreateTableRegex() PrimaryKey[%d] = %v, want %v", i, result.PrimaryKey[i], pk)
				}
			}
			if len(result.ForeignKeys) != tt.expectedFKs {
				t.Errorf("parseCreateTableRegex() ForeignKeys count = %v, want %v", len(result.ForeignKeys), tt.expectedFKs)
			}
		})
	}
}

// Helper functions for pointer comparisons in tests
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func compareIntPtr(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func compareStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
