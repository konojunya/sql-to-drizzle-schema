package generator

import (
	"fmt"
	"os"

	"github.com/konojunya/sql-to-drizzle-schema/internal/parser"
)

// NewSchemaGenerator creates a new schema generator for the specified dialect
func NewSchemaGenerator(dialect parser.DatabaseDialect) (SchemaGenerator, error) {
	switch dialect {
	case parser.PostgreSQL:
		return NewPostgreSQLSchemaGenerator(), nil
	case parser.MySQL:
		return nil, fmt.Errorf("MySQL schema generation is not yet implemented")
	case parser.Spanner:
		return nil, fmt.Errorf("Spanner schema generation is not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", dialect)
	}
}

// GenerateSchemaToFile is a convenience function that generates schema and writes to file
func GenerateSchemaToFile(tables []parser.Table, dialect parser.DatabaseDialect, outputFile string, options GeneratorOptions) error {
	generator, err := NewSchemaGenerator(dialect)
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	schema, err := generator.GenerateSchema(tables, options)
	if err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}

	err = WriteSchemaToFile(schema.Content, outputFile)
	if err != nil {
		return fmt.Errorf("failed to write schema to file: %w", err)
	}

	return nil
}

// WriteSchemaToFile writes the generated schema content to a file
func WriteSchemaToFile(content, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write content to file %s: %w", filename, err)
	}

	return nil
}