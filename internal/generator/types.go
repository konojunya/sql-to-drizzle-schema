// Package generator provides functionality to generate Drizzle ORM schema definitions
// from parsed SQL table structures.
//
// This package converts the parsed SQL structures into TypeScript code that uses
// Drizzle ORM syntax for different database dialects.
package generator

import "github.com/konojunya/sql-to-drizzle-schema/internal/parser"

// GeneratorOptions contains options for schema generation
type GeneratorOptions struct {
	// TableNameCase specifies the naming convention for table exports
	TableNameCase NamingCase
	// ColumnNameCase specifies the naming convention for column names
	ColumnNameCase NamingCase
	// IncludeComments includes comments in the generated schema
	IncludeComments bool
	// ExportPrefix adds a prefix to exported table names
	ExportPrefix string
	// IndentSize specifies the number of spaces for indentation
	IndentSize int
}

// NamingCase represents different naming conventions
type NamingCase string

const (
	// CamelCase converts to camelCase (userProfiles)
	CamelCase NamingCase = "camel"
	// PascalCase converts to PascalCase (UserProfiles)
	PascalCase NamingCase = "pascal"
	// SnakeCase keeps snake_case (user_profiles)
	SnakeCase NamingCase = "snake"
	// KebabCase converts to kebab-case (user-profiles)
	KebabCase NamingCase = "kebab"
)

// GeneratedSchema represents the complete generated schema
type GeneratedSchema struct {
	// Imports contains the import statements needed for the schema
	Imports []string
	// Tables contains the generated table definitions
	Tables []GeneratedTable
	// Content contains the complete generated TypeScript content
	Content string
}

// GeneratedTable represents a single generated table definition
type GeneratedTable struct {
	// OriginalName is the original SQL table name
	OriginalName string
	// ExportName is the exported TypeScript variable name
	ExportName string
	// Definition contains the table definition code
	Definition string
}

// DrizzleType represents a Drizzle ORM column type
type DrizzleType struct {
	// Function is the Drizzle function name (e.g., "varchar", "bigserial")
	Function string
	// Args contains arguments for the function
	Args []string
	// Options contains method chain options (e.g., ".notNull()", ".default()")
	Options []string
}

// SchemaGenerator interface defines the contract for schema generation
type SchemaGenerator interface {
	// GenerateSchema generates a complete Drizzle schema from parsed tables
	GenerateSchema(tables []parser.Table, options GeneratorOptions) (*GeneratedSchema, error)

	// GenerateTable generates a single table definition
	GenerateTable(table parser.Table, options GeneratorOptions) (*GeneratedTable, error)

	// SupportedDialect returns the database dialect this generator supports
	SupportedDialect() parser.DatabaseDialect
}

// ColumnTypeMapper interface defines the contract for mapping SQL types to Drizzle types
type ColumnTypeMapper interface {
	// MapColumnType maps a SQL column to a Drizzle type definition
	MapColumnType(column parser.Column) (*DrizzleType, error)

	// SupportedDialect returns the database dialect this mapper supports
	SupportedDialect() parser.DatabaseDialect
}

// DefaultGeneratorOptions returns sensible default options for schema generation
func DefaultGeneratorOptions() GeneratorOptions {
	return GeneratorOptions{
		TableNameCase:   CamelCase,
		ColumnNameCase:  CamelCase,
		IncludeComments: true,
		ExportPrefix:    "",
		IndentSize:      2,
	}
}
