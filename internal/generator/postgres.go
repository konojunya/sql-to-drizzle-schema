package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/konojunya/sql-to-drizzle-schema/internal/parser"
)

// PostgreSQLTypeMapper implements type mapping for PostgreSQL to Drizzle ORM
type PostgreSQLTypeMapper struct{}

// NewPostgreSQLTypeMapper creates a new PostgreSQL type mapper
func NewPostgreSQLTypeMapper() *PostgreSQLTypeMapper {
	return &PostgreSQLTypeMapper{}
}

// SupportedDialect returns the database dialect this mapper supports
func (m *PostgreSQLTypeMapper) SupportedDialect() parser.DatabaseDialect {
	return parser.PostgreSQL
}

// MapColumnType maps a PostgreSQL column to a Drizzle type definition
func (m *PostgreSQLTypeMapper) MapColumnType(column parser.Column) (*DrizzleType, error) {
	drizzleType := &DrizzleType{
		Function: "",
		Args:     []string{},
		Options:  []string{},
	}

	// Map SQL types to Drizzle types
	switch strings.ToUpper(column.Type) {
	case "BIGSERIAL":
		drizzleType.Function = "bigserial"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name), "{ mode: 'number' }"}
	case "SERIAL":
		drizzleType.Function = "serial"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "SMALLSERIAL":
		drizzleType.Function = "serial"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "BIGINT":
		drizzleType.Function = "bigint"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name), "{ mode: 'number' }"}
	case "INTEGER", "INT", "INT4":
		drizzleType.Function = "integer"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "SMALLINT", "INT2":
		drizzleType.Function = "smallint"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "VARCHAR":
		if column.Length != nil {
			drizzleType.Function = "varchar"
			drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name), fmt.Sprintf("{ length: %d }", *column.Length)}
		} else {
			drizzleType.Function = "varchar"
			drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
		}
	case "TEXT":
		drizzleType.Function = "text"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "BOOLEAN", "BOOL":
		drizzleType.Function = "boolean"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "TIMESTAMP WITH TIME ZONE", "TIMESTAMPTZ":
		drizzleType.Function = "timestamp"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name), "{ withTimezone: true }"}
	case "TIMESTAMP":
		drizzleType.Function = "timestamp"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "DATE":
		drizzleType.Function = "date"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "TIME":
		drizzleType.Function = "time"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "DECIMAL", "NUMERIC":
		if column.Length != nil && column.Scale != nil {
			drizzleType.Function = "decimal"
			drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name), fmt.Sprintf("{ precision: %d, scale: %d }", *column.Length, *column.Scale)}
		} else if column.Length != nil {
			drizzleType.Function = "decimal"
			drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name), fmt.Sprintf("{ precision: %d }", *column.Length)}
		} else {
			drizzleType.Function = "decimal"
			drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
		}
	case "REAL", "FLOAT4":
		drizzleType.Function = "real"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "DOUBLE PRECISION", "DOUBLE", "FLOAT8":
		drizzleType.Function = "doublePrecision"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "UUID":
		drizzleType.Function = "uuid"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "JSON":
		drizzleType.Function = "json"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	case "JSONB":
		drizzleType.Function = "jsonb"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	default:
		// Fallback to text for unknown types
		drizzleType.Function = "text"
		drizzleType.Args = []string{fmt.Sprintf("'%s'", column.Name)}
	}

	// Add constraints as method chains
	if column.NotNull {
		drizzleType.Options = append(drizzleType.Options, "notNull()")
	}

	if column.Unique {
		drizzleType.Options = append(drizzleType.Options, "unique()")
	}

	// Handle default values
	if column.DefaultValue != nil {
		defaultVal := *column.DefaultValue
		switch strings.ToUpper(defaultVal) {
		case "CURRENT_TIMESTAMP", "NOW()":
			if strings.Contains(strings.ToUpper(column.Type), "TIMESTAMP") {
				drizzleType.Options = append(drizzleType.Options, "defaultNow()")
			}
		case "TRUE":
			drizzleType.Options = append(drizzleType.Options, "default(true)")
		case "FALSE":
			drizzleType.Options = append(drizzleType.Options, "default(false)")
		default:
			// For string literals, keep quotes; for numbers, don't quote
			if strings.HasPrefix(defaultVal, "'") && strings.HasSuffix(defaultVal, "'") {
				drizzleType.Options = append(drizzleType.Options, fmt.Sprintf("default(%s)", defaultVal))
			} else if _, err := strconv.Atoi(defaultVal); err == nil {
				// It's a number
				drizzleType.Options = append(drizzleType.Options, fmt.Sprintf("default(%s)", defaultVal))
			} else {
				// Treat as string literal
				drizzleType.Options = append(drizzleType.Options, fmt.Sprintf("default('%s')", defaultVal))
			}
		}
	}

	return drizzleType, nil
}

// PostgreSQLSchemaGenerator implements schema generation for PostgreSQL
type PostgreSQLSchemaGenerator struct {
	typeMapper *PostgreSQLTypeMapper
}

// NewPostgreSQLSchemaGenerator creates a new PostgreSQL schema generator
func NewPostgreSQLSchemaGenerator() *PostgreSQLSchemaGenerator {
	return &PostgreSQLSchemaGenerator{
		typeMapper: NewPostgreSQLTypeMapper(),
	}
}

// SupportedDialect returns the database dialect this generator supports
func (g *PostgreSQLSchemaGenerator) SupportedDialect() parser.DatabaseDialect {
	return parser.PostgreSQL
}

// GenerateSchema generates a complete Drizzle schema from parsed tables
func (g *PostgreSQLSchemaGenerator) GenerateSchema(tables []parser.Table, options GeneratorOptions) (*GeneratedSchema, error) {
	schema := &GeneratedSchema{
		Imports: []string{},
		Tables:  []GeneratedTable{},
	}

	// Collect required imports
	importSet := make(map[string]bool)
	importSet["pgTable"] = true // Always need pgTable

	// First pass: collect all required imports
	for _, table := range tables {
		for _, column := range table.Columns {
			drizzleType, err := g.typeMapper.MapColumnType(column)
			if err != nil {
				return nil, fmt.Errorf("failed to map column %s.%s: %w", table.Name, column.Name, err)
			}
			importSet[drizzleType.Function] = true
		}
	}

	// Generate import statement
	var importList []string
	for imp := range importSet {
		importList = append(importList, imp)
	}

	// Sort imports for consistency (basic alphabetical)
	for i := 0; i < len(importList); i++ {
		for j := i + 1; j < len(importList); j++ {
			if importList[i] > importList[j] {
				importList[i], importList[j] = importList[j], importList[i]
			}
		}
	}

	schema.Imports = []string{fmt.Sprintf("import { %s } from 'drizzle-orm/pg-core';", strings.Join(importList, ", "))}

	// Sort tables to handle foreign key dependencies
	// Tables without foreign keys first, then tables with foreign keys
	sortedTables := g.sortTablesByDependencies(tables)

	// Generate table definitions in dependency order
	for _, table := range sortedTables {
		generatedTable, err := g.GenerateTable(table, options)
		if err != nil {
			return nil, fmt.Errorf("failed to generate table %s: %w", table.Name, err)
		}
		schema.Tables = append(schema.Tables, *generatedTable)
	}

	// Build complete content
	var contentBuilder strings.Builder

	// Add imports
	for _, imp := range schema.Imports {
		contentBuilder.WriteString(imp)
		contentBuilder.WriteString("\n")
	}
	contentBuilder.WriteString("\n")

	// Add table definitions
	for i, table := range schema.Tables {
		if i > 0 {
			contentBuilder.WriteString("\n")
		}
		contentBuilder.WriteString(table.Definition)
		contentBuilder.WriteString("\n")
	}

	schema.Content = contentBuilder.String()
	return schema, nil
}

// sortTablesByDependencies sorts tables so that referenced tables come before referencing tables
func (g *PostgreSQLSchemaGenerator) sortTablesByDependencies(tables []parser.Table) []parser.Table {
	// Create a map for quick lookup
	tableMap := make(map[string]parser.Table)
	for _, table := range tables {
		tableMap[table.Name] = table
	}

	// Simple topological sort
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	sorted := []parser.Table{}

	var visit func(tableName string)
	visit = func(tableName string) {
		if visited[tableName] || visiting[tableName] {
			return
		}

		visiting[tableName] = true
		table := tableMap[tableName]

		// Visit all dependencies (referenced tables) first
		for _, fk := range table.ForeignKeys {
			if _, exists := tableMap[fk.ReferencedTable]; exists {
				visit(fk.ReferencedTable)
			}
		}

		visiting[tableName] = false
		visited[tableName] = true
		sorted = append(sorted, table)
	}

	// Visit all tables
	for _, table := range tables {
		visit(table.Name)
	}

	return sorted
}

// GenerateTable generates a single table definition
func (g *PostgreSQLSchemaGenerator) GenerateTable(table parser.Table, options GeneratorOptions) (*GeneratedTable, error) {
	exportName := g.convertCase(table.Name, options.TableNameCase)

	var builder strings.Builder
	indent := strings.Repeat(" ", options.IndentSize)

	// Add comment if enabled
	if options.IncludeComments {
		builder.WriteString(fmt.Sprintf("// %s table\n", table.Name))
	}

	// Start table definition
	builder.WriteString(fmt.Sprintf("export const %s%s = pgTable('%s', {\n", options.ExportPrefix, exportName, table.Name))

	// Generate columns
	for i, column := range table.Columns {
		drizzleType, err := g.typeMapper.MapColumnType(column)
		if err != nil {
			return nil, fmt.Errorf("failed to map column %s: %w", column.Name, err)
		}

		columnName := g.convertCase(column.Name, options.ColumnNameCase)

		// Build column definition
		builder.WriteString(fmt.Sprintf("%s%s: %s(%s)", indent, columnName, drizzleType.Function, strings.Join(drizzleType.Args, ", ")))

		// Add method chains
		for _, option := range drizzleType.Options {
			builder.WriteString(fmt.Sprintf(".%s", option))
		}

		// Add primary key if this column is in the primary key
		for _, pkCol := range table.PrimaryKey {
			if pkCol == column.Name {
				builder.WriteString(".primaryKey()")
				break
			}
		}

		// Add foreign key reference if this column has one
		for _, fk := range table.ForeignKeys {
			// Check if this column is part of a foreign key (support single-column FKs for now)
			if len(fk.Columns) == 1 && fk.Columns[0] == column.Name {
				referencedTableName := g.convertCase(fk.ReferencedTable, options.TableNameCase)
				if len(fk.ReferencedColumns) == 1 {
					referencedColumnName := g.convertCase(fk.ReferencedColumns[0], options.ColumnNameCase)
					builder.WriteString(fmt.Sprintf(".references(() => %s.%s)", referencedTableName, referencedColumnName))
				}
				break
			}
		}

		// Add comma except for last column
		if i < len(table.Columns)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("});")

	return &GeneratedTable{
		OriginalName: table.Name,
		ExportName:   exportName,
		Definition:   builder.String(),
	}, nil
}

// convertCase converts a string to the specified naming case
func (g *PostgreSQLSchemaGenerator) convertCase(input string, caseType NamingCase) string {
	switch caseType {
	case CamelCase:
		return g.toCamelCase(input)
	case PascalCase:
		return g.toPascalCase(input)
	case SnakeCase:
		return input // Keep as-is
	case KebabCase:
		return strings.ReplaceAll(input, "_", "-")
	default:
		return input
	}
}

// toCamelCase converts snake_case to camelCase
func (g *PostgreSQLSchemaGenerator) toCamelCase(input string) string {
	words := strings.Split(input, "_")
	if len(words) == 0 {
		return input
	}

	result := words[0]
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result += strings.ToUpper(words[i][:1]) + words[i][1:]
		}
	}
	return result
}

// toPascalCase converts snake_case to PascalCase
func (g *PostgreSQLSchemaGenerator) toPascalCase(input string) string {
	words := strings.Split(input, "_")
	var result string

	for _, word := range words {
		if len(word) > 0 {
			result += strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return result
}
