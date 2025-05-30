package parser

import "fmt"

// NewParser creates a new SQL parser for the specified dialect
func NewParser(dialect DatabaseDialect) (SQLParser, error) {
	switch dialect {
	case PostgreSQL:
		return NewPostgreSQLParser(), nil
	case MySQL:
		return nil, fmt.Errorf("MySQL dialect support is not yet implemented")
	case Spanner:
		return nil, fmt.Errorf("Spanner dialect support is not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", dialect)
	}
}

// ParseSQLContent is a convenience function that creates a parser and parses SQL content
func ParseSQLContent(content string, dialect DatabaseDialect, options ParseOptions) (*ParseResult, error) {
	parser, err := NewParser(dialect)
	if err != nil {
		return nil, err
	}

	// Set the dialect in options if not already set
	if options.Dialect == "" {
		options.Dialect = dialect
	}

	return parser.ParseSQL(content, options)
}

// DefaultParseOptions returns sensible default options for parsing
func DefaultParseOptions() ParseOptions {
	return ParseOptions{
		Dialect:           PostgreSQL,
		StrictMode:        false,
		IgnoreUnsupported: true,
	}
}
