// Package parser provides SQL parsing functionality for converting SQL DDL
// statements to structured data that can be used to generate Drizzle ORM schemas.
//
// This package currently supports PostgreSQL syntax and will be extended to support
// MySQL and Spanner in future versions.
package parser

// DatabaseDialect represents the SQL dialect being parsed
type DatabaseDialect string

const (
	// PostgreSQL dialect
	PostgreSQL DatabaseDialect = "postgresql"
	// MySQL dialect (future support)
	MySQL DatabaseDialect = "mysql"
	// Spanner dialect (future support)
	Spanner DatabaseDialect = "spanner"
)

// Table represents a parsed SQL table definition
type Table struct {
	// Name is the table name
	Name string
	// Columns contains all column definitions
	Columns []Column
	// PrimaryKey contains primary key column names
	PrimaryKey []string
	// ForeignKeys contains foreign key constraints
	ForeignKeys []ForeignKey
	// Indexes contains index definitions
	Indexes []Index
	// Constraints contains other constraints (unique, check, etc.)
	Constraints []Constraint
}

// Column represents a parsed column definition
type Column struct {
	// Name is the column name
	Name string
	// Type is the SQL data type (e.g., "VARCHAR", "BIGINT", "TIMESTAMP")
	Type string
	// Length is the column length for types that support it (e.g., VARCHAR(255))
	Length *int
	// Precision is the precision for decimal types
	Precision *int
	// Scale is the scale for decimal types
	Scale *int
	// NotNull indicates if the column has NOT NULL constraint
	NotNull bool
	// Unique indicates if the column has UNIQUE constraint
	Unique bool
	// DefaultValue contains the default value expression if specified
	DefaultValue *string
	// AutoIncrement indicates if the column is auto-incrementing (SERIAL, AUTO_INCREMENT)
	AutoIncrement bool
	// Comment contains column comment if specified
	Comment *string
}

// ForeignKey represents a foreign key constraint
type ForeignKey struct {
	// Name is the constraint name
	Name string
	// Columns are the local columns in the foreign key
	Columns []string
	// ReferencedTable is the referenced table name
	ReferencedTable string
	// ReferencedColumns are the referenced columns
	ReferencedColumns []string
	// OnDelete specifies the action on delete (CASCADE, SET NULL, etc.)
	OnDelete *string
	// OnUpdate specifies the action on update
	OnUpdate *string
}

// Index represents an index definition
type Index struct {
	// Name is the index name
	Name string
	// Columns are the indexed columns
	Columns []string
	// Unique indicates if this is a unique index
	Unique bool
	// Type is the index type (BTREE, HASH, etc.)
	Type *string
}

// Constraint represents a table constraint
type Constraint struct {
	// Name is the constraint name
	Name string
	// Type is the constraint type (CHECK, UNIQUE, etc.)
	Type string
	// Columns are the columns involved in the constraint
	Columns []string
	// Expression is the constraint expression (for CHECK constraints)
	Expression *string
}

// ParseResult contains the results of parsing a SQL file
type ParseResult struct {
	// Tables contains all parsed table definitions
	Tables []Table
	// Dialect is the detected or specified SQL dialect
	Dialect DatabaseDialect
	// Errors contains any parsing errors encountered
	Errors []error
}

// ParseOptions contains options for the SQL parser
type ParseOptions struct {
	// Dialect specifies the SQL dialect to use for parsing
	Dialect DatabaseDialect
	// StrictMode enables strict parsing (fails on unsupported features)
	StrictMode bool
	// IgnoreUnsupported ignores unsupported SQL features instead of failing
	IgnoreUnsupported bool
}

// SQLParser interface defines the contract for SQL parsing implementations
type SQLParser interface {
	// ParseSQL parses SQL content and returns structured table definitions
	ParseSQL(content string, options ParseOptions) (*ParseResult, error)

	// SupportedDialect returns the SQL dialect this parser supports
	SupportedDialect() DatabaseDialect
}
