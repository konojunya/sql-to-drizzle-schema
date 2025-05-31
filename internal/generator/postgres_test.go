package generator

import (
	"strings"
	"testing"

	"github.com/konojunya/sql-to-drizzle-schema/internal/parser"
)

func TestNewPostgreSQLTypeMapper(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()
	if mapper == nil {
		t.Errorf("NewPostgreSQLTypeMapper() returned nil")
	}
	if mapper.SupportedDialect() != parser.PostgreSQL {
		t.Errorf("NewPostgreSQLTypeMapper() SupportedDialect() = %v, want %v", mapper.SupportedDialect(), parser.PostgreSQL)
	}
}

func TestNewPostgreSQLSchemaGenerator(t *testing.T) {
	generator := NewPostgreSQLSchemaGenerator()
	if generator == nil {
		t.Errorf("NewPostgreSQLSchemaGenerator() returned nil")
	}
	if generator.SupportedDialect() != parser.PostgreSQL {
		t.Errorf("NewPostgreSQLSchemaGenerator() SupportedDialect() = %v, want %v", generator.SupportedDialect(), parser.PostgreSQL)
	}
}

func TestPostgreSQLTypeMapper_MapColumnType(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name         string
		column       parser.Column
		expectedFunc string
		expectedArgs []string
		expectedOpts []string
		wantErr      bool
	}{
		{
			name: "BIGSERIAL column",
			column: parser.Column{
				Name:          "id",
				Type:          "BIGSERIAL",
				NotNull:       true,
				AutoIncrement: true,
			},
			expectedFunc: "bigserial",
			expectedArgs: []string{"'id'", "{ mode: 'number' }"},
			expectedOpts: []string{"notNull()"},
			wantErr:      false,
		},
		{
			name: "VARCHAR with length",
			column: parser.Column{
				Name:    "name",
				Type:    "VARCHAR",
				Length:  intPtr(255),
				NotNull: true,
			},
			expectedFunc: "varchar",
			expectedArgs: []string{"'name'", "{ length: 255 }"},
			expectedOpts: []string{"notNull()"},
			wantErr:      false,
		},
		{
			name: "TEXT column",
			column: parser.Column{
				Name:    "content",
				Type:    "TEXT",
				NotNull: true,
			},
			expectedFunc: "text",
			expectedArgs: []string{"'content'"},
			expectedOpts: []string{"notNull()"},
			wantErr:      false,
		},
		{
			name: "BOOLEAN with default",
			column: parser.Column{
				Name:         "active",
				Type:         "BOOLEAN",
				NotNull:      true,
				DefaultValue: stringPtr("TRUE"),
			},
			expectedFunc: "boolean",
			expectedArgs: []string{"'active'"},
			expectedOpts: []string{"notNull()", "default(true)"},
			wantErr:      false,
		},
		{
			name: "TIMESTAMP WITH TIME ZONE with defaultNow",
			column: parser.Column{
				Name:         "created_at",
				Type:         "TIMESTAMP WITH TIME ZONE",
				NotNull:      true,
				DefaultValue: stringPtr("CURRENT_TIMESTAMP"),
			},
			expectedFunc: "timestamp",
			expectedArgs: []string{"'created_at'", "{ withTimezone: true }"},
			expectedOpts: []string{"notNull()", "defaultNow()"},
			wantErr:      false,
		},
		{
			name: "VARCHAR with UNIQUE constraint",
			column: parser.Column{
				Name:    "email",
				Type:    "VARCHAR",
				Length:  intPtr(255),
				NotNull: true,
				Unique:  true,
			},
			expectedFunc: "varchar",
			expectedArgs: []string{"'email'", "{ length: 255 }"},
			expectedOpts: []string{"notNull()", "unique()"},
			wantErr:      false,
		},
		{
			name: "DECIMAL with precision and scale",
			column: parser.Column{
				Name:    "price",
				Type:    "DECIMAL",
				Length:  intPtr(10),
				Scale:   intPtr(2),
				NotNull: true,
			},
			expectedFunc: "decimal",
			expectedArgs: []string{"'price'", "{ precision: 10, scale: 2 }"},
			expectedOpts: []string{"notNull()"},
			wantErr:      false,
		},
		{
			name: "VARCHAR with string default",
			column: parser.Column{
				Name:         "role",
				Type:         "VARCHAR",
				Length:       intPtr(50),
				NotNull:      true,
				DefaultValue: stringPtr("'user'"),
			},
			expectedFunc: "varchar",
			expectedArgs: []string{"'role'", "{ length: 50 }"},
			expectedOpts: []string{"notNull()", "default('user')"},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapper.MapColumnType(tt.column)

			if tt.wantErr && err == nil {
				t.Errorf("MapColumnType() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("MapColumnType() unexpected error: %v", err)
				return
			}
			if tt.wantErr {
				return
			}

			if result.Function != tt.expectedFunc {
				t.Errorf("MapColumnType() Function = %v, want %v", result.Function, tt.expectedFunc)
			}
			if !slicesEqual(result.Args, tt.expectedArgs) {
				t.Errorf("MapColumnType() Args = %v, want %v", result.Args, tt.expectedArgs)
			}
			if !slicesEqual(result.Options, tt.expectedOpts) {
				t.Errorf("MapColumnType() Options = %v, want %v", result.Options, tt.expectedOpts)
			}
		})
	}
}

func TestPostgreSQLSchemaGenerator_GenerateTable(t *testing.T) {
	generator := NewPostgreSQLSchemaGenerator()
	options := DefaultGeneratorOptions()

	tests := []struct {
		name            string
		table           parser.Table
		options         GeneratorOptions
		expectedExport  string
		expectedContent []string
		wantErr         bool
	}{
		{
			name: "Simple table",
			table: parser.Table{
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
			options:        options,
			expectedExport: "users",
			expectedContent: []string{
				"export const users = pgTable('users', {",
				"id: bigserial('id', { mode: 'number' }).notNull().primaryKey()",
				"name: varchar('name', { length: 255 }).notNull()",
				"});",
			},
			wantErr: false,
		},
		{
			name: "Table with foreign key",
			table: parser.Table{
				Name: "posts",
				Columns: []parser.Column{
					{
						Name:    "id",
						Type:    "BIGSERIAL",
						NotNull: true,
					},
					{
						Name:    "user_id",
						Type:    "BIGINT",
						NotNull: true,
					},
				},
				PrimaryKey: []string{"id"},
				ForeignKeys: []parser.ForeignKey{
					{
						Name:              "fk_posts_users",
						Columns:           []string{"user_id"},
						ReferencedTable:   "users",
						ReferencedColumns: []string{"id"},
					},
				},
			},
			options:        options,
			expectedExport: "posts",
			expectedContent: []string{
				"export const posts = pgTable('posts', {",
				"id: bigserial('id', { mode: 'number' }).notNull().primaryKey()",
				"userId: bigint('user_id', { mode: 'number' }).notNull().references(() => users.id)",
				"});",
			},
			wantErr: false,
		},
		{
			name: "Table with unique constraint",
			table: parser.Table{
				Name: "role_permissions",
				Columns: []parser.Column{
					{
						Name:    "role_id",
						Type:    "BIGINT",
						NotNull: true,
					},
					{
						Name:    "permission_id",
						Type:    "BIGINT",
						NotNull: true,
					},
				},
				Constraints: []parser.Constraint{
					{
						Name:    "unique_role_permission",
						Type:    "UNIQUE",
						Columns: []string{"role_id", "permission_id"},
					},
				},
			},
			options:        options,
			expectedExport: "rolePermissions",
			expectedContent: []string{
				"export const rolePermissions = pgTable('role_permissions', {",
				"roleId: bigint('role_id', { mode: 'number' }).notNull()",
				"permissionId: bigint('permission_id', { mode: 'number' }).notNull()",
				"});",
				"export const uniqueRolePermission = unique('unique_role_permission').on(rolePermissions.roleId, rolePermissions.permissionId);",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generator.GenerateTable(tt.table, tt.options)

			if tt.wantErr && err == nil {
				t.Errorf("GenerateTable() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("GenerateTable() unexpected error: %v", err)
				return
			}
			if tt.wantErr {
				return
			}

			if result.OriginalName != tt.table.Name {
				t.Errorf("GenerateTable() OriginalName = %v, want %v", result.OriginalName, tt.table.Name)
			}
			if result.ExportName != tt.expectedExport {
				t.Errorf("GenerateTable() ExportName = %v, want %v", result.ExportName, tt.expectedExport)
			}

			// Check that expected content strings are present
			for _, expected := range tt.expectedContent {
				if !strings.Contains(result.Definition, expected) {
					t.Errorf("GenerateTable() Definition missing expected content: %s\nActual:\n%s", expected, result.Definition)
				}
			}
		})
	}
}

func TestPostgreSQLSchemaGenerator_GenerateSchema(t *testing.T) {
	generator := NewPostgreSQLSchemaGenerator()
	options := DefaultGeneratorOptions()

	tests := []struct {
		name            string
		tables          []parser.Table
		options         GeneratorOptions
		expectedTables  int
		expectedImports []string
		wantErr         bool
	}{
		{
			name: "Single table schema",
			tables: []parser.Table{
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
				},
			},
			options:        options,
			expectedTables: 1,
			expectedImports: []string{
				"bigserial",
				"pgTable",
				"varchar",
			},
			wantErr: false,
		},
		{
			name: "Multiple tables with dependencies",
			tables: []parser.Table{
				{
					Name: "posts",
					Columns: []parser.Column{
						{
							Name:    "id",
							Type:    "BIGSERIAL",
							NotNull: true,
						},
						{
							Name:    "user_id",
							Type:    "BIGINT",
							NotNull: true,
						},
					},
					ForeignKeys: []parser.ForeignKey{
						{
							Columns:           []string{"user_id"},
							ReferencedTable:   "users",
							ReferencedColumns: []string{"id"},
						},
					},
				},
				{
					Name: "users",
					Columns: []parser.Column{
						{
							Name:    "id",
							Type:    "BIGSERIAL",
							NotNull: true,
						},
					},
				},
			},
			options:        options,
			expectedTables: 2,
			expectedImports: []string{
				"bigint",
				"bigserial",
				"pgTable",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generator.GenerateSchema(tt.tables, tt.options)

			if tt.wantErr && err == nil {
				t.Errorf("GenerateSchema() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("GenerateSchema() unexpected error: %v", err)
				return
			}
			if tt.wantErr {
				return
			}

			if len(result.Tables) != tt.expectedTables {
				t.Errorf("GenerateSchema() Tables count = %v, want %v", len(result.Tables), tt.expectedTables)
			}

			// Check imports are present
			importStr := strings.Join(result.Imports, " ")
			for _, expectedImport := range tt.expectedImports {
				if !strings.Contains(importStr, expectedImport) {
					t.Errorf("GenerateSchema() missing expected import: %s in %s", expectedImport, importStr)
				}
			}

			// Check content is generated
			if result.Content == "" {
				t.Errorf("GenerateSchema() Content is empty")
			}
		})
	}
}

func TestPostgreSQLSchemaGenerator_convertCase(t *testing.T) {
	generator := NewPostgreSQLSchemaGenerator()

	tests := []struct {
		name     string
		input    string
		caseType NamingCase
		expected string
	}{
		{
			name:     "snake_case to camelCase",
			input:    "user_profiles",
			caseType: CamelCase,
			expected: "userProfiles",
		},
		{
			name:     "snake_case to PascalCase",
			input:    "user_profiles",
			caseType: PascalCase,
			expected: "UserProfiles",
		},
		{
			name:     "snake_case to snake_case",
			input:    "user_profiles",
			caseType: SnakeCase,
			expected: "user_profiles",
		},
		{
			name:     "snake_case to kebab-case",
			input:    "user_profiles",
			caseType: KebabCase,
			expected: "user-profiles",
		},
		{
			name:     "single word",
			input:    "users",
			caseType: CamelCase,
			expected: "users",
		},
		{
			name:     "single word to PascalCase",
			input:    "users",
			caseType: PascalCase,
			expected: "Users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.convertCase(tt.input, tt.caseType)
			if result != tt.expected {
				t.Errorf("convertCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPostgreSQLSchemaGenerator_sortTablesByDependencies(t *testing.T) {
	generator := NewPostgreSQLSchemaGenerator()

	tables := []parser.Table{
		{
			Name: "comments",
			ForeignKeys: []parser.ForeignKey{
				{Columns: []string{"user_id"}, ReferencedTable: "users"},
				{Columns: []string{"post_id"}, ReferencedTable: "posts"},
			},
		},
		{
			Name: "posts",
			ForeignKeys: []parser.ForeignKey{
				{Columns: []string{"user_id"}, ReferencedTable: "users"},
			},
		},
		{
			Name: "users",
		},
	}

	result := generator.sortTablesByDependencies(tables)

	// users should come first (no dependencies)
	// posts should come second (depends on users)
	// comments should come last (depends on both users and posts)
	expectedOrder := []string{"users", "posts", "comments"}

	if len(result) != len(expectedOrder) {
		t.Errorf("sortTablesByDependencies() returned %d tables, want %d", len(result), len(expectedOrder))
		return
	}

	for i, expectedName := range expectedOrder {
		if result[i].Name != expectedName {
			t.Errorf("sortTablesByDependencies() table[%d] = %s, want %s", i, result[i].Name, expectedName)
		}
	}
}

// Helper functions for tests
func stringPtr(s string) *string {
	return &s
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
