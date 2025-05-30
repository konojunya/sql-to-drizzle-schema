package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PostgreSQLParser implements SQL parsing for PostgreSQL dialect
type PostgreSQLParser struct{}

// NewPostgreSQLParser creates a new PostgreSQL parser
func NewPostgreSQLParser() *PostgreSQLParser {
	return &PostgreSQLParser{}
}

// SupportedDialect returns the SQL dialect this parser supports
func (p *PostgreSQLParser) SupportedDialect() DatabaseDialect {
	return PostgreSQL
}

// ParseSQL parses PostgreSQL SQL content and returns structured table definitions
func (p *PostgreSQLParser) ParseSQL(content string, options ParseOptions) (*ParseResult, error) {
	result := &ParseResult{
		Tables:  []Table{},
		Dialect: PostgreSQL,
		Errors:  []error{},
	}

	// Split content into individual statements
	statements := p.splitStatements(content)

	for _, stmtStr := range statements {
		// Skip empty statements and comments
		stmtStr = strings.TrimSpace(stmtStr)
		if stmtStr == "" {
			continue
		}

		// Remove leading comments but keep the rest
		lines := strings.Split(stmtStr, "\n")
		var cleanLines []string
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmedLine, "--") && trimmedLine != "" {
				cleanLines = append(cleanLines, line)
			}
		}

		if len(cleanLines) == 0 {
			continue
		}

		stmtStr = strings.Join(cleanLines, "\n")

		// Use regex-based parsing for CREATE TABLE statements
		if p.isCreateTableStatement(stmtStr) {
			table, err := p.parseCreateTableRegex(stmtStr, options)
			if err != nil {
				if options.IgnoreUnsupported {
					result.Errors = append(result.Errors, err)
					continue
				}
				return nil, err
			}
			if table != nil {
				result.Tables = append(result.Tables, *table)
			}
		}
	}

	return result, nil
}

// isCreateTableStatement checks if a statement is a CREATE TABLE statement
func (p *PostgreSQLParser) isCreateTableStatement(stmt string) bool {
	// Simple regex to match CREATE TABLE statements
	createTableRegex := regexp.MustCompile(`(?i)^\s*CREATE\s+TABLE\s+`)
	return createTableRegex.MatchString(stmt)
}

// parseCreateTableRegex parses a CREATE TABLE statement using regex
func (p *PostgreSQLParser) parseCreateTableRegex(stmt string, options ParseOptions) (*Table, error) {
	// Extract table name
	tableNameRegex := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)\s*\(`)
	matches := tableNameRegex.FindStringSubmatch(stmt)
	if len(matches) < 2 {
		return nil, fmt.Errorf("could not extract table name from statement")
	}

	table := &Table{
		Name:        matches[1],
		Columns:     []Column{},
		PrimaryKey:  []string{},
		ForeignKeys: []ForeignKey{},
		Indexes:     []Index{},
		Constraints: []Constraint{},
	}

	// Extract table body (everything between the first ( and last ))
	// Use DOTALL flag to match across newlines
	bodyRegex := regexp.MustCompile(`(?is)CREATE\s+TABLE\s+\w+\s*\((.*)\);?\s*$`)
	bodyMatches := bodyRegex.FindStringSubmatch(stmt)
	if len(bodyMatches) < 2 {
		return nil, fmt.Errorf("could not extract table body from statement")
	}

	tableBody := bodyMatches[1]

	// Parse columns and constraints
	err := p.parseTableBody(table, tableBody, options)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table body: %w", err)
	}

	return table, nil
}

// parseTableBody parses the table body containing columns and constraints
func (p *PostgreSQLParser) parseTableBody(table *Table, body string, options ParseOptions) error {
	// Split by commas, but be careful about parentheses and strings
	items := p.splitTableItems(body)

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		// Check if it's a constraint
		if p.isConstraint(item) {
			err := p.parseConstraint(table, item, options)
			if err != nil && !options.IgnoreUnsupported {
				return err
			}
		} else {
			// It's a column definition
			column, err := p.parseColumnRegex(item, options)
			if err != nil {
				if options.IgnoreUnsupported {
					continue
				}
				return err
			}
			table.Columns = append(table.Columns, *column)
		}
	}

	return nil
}

// parseColumnRegex parses a column definition using regex
func (p *PostgreSQLParser) parseColumnRegex(columnDef string, options ParseOptions) (*Column, error) {
	// Basic column regex: name type [constraints...]
	// Allow more flexible type matching including WITH TIME ZONE
	columnRegex := regexp.MustCompile(`(?i)^\s*(\w+)\s+((?:[A-Z]+(?:\([^)]*\))?(?:\s+WITH\s+TIME\s+ZONE)?)+)\s*(.*)$`)
	matches := columnRegex.FindStringSubmatch(columnDef)

	if len(matches) < 3 {
		return nil, fmt.Errorf("could not parse column definition: %s", columnDef)
	}

	column := &Column{
		Name:          matches[1],
		Type:          strings.ToUpper(strings.TrimSpace(matches[2])),
		NotNull:       false,
		Unique:        false,
		AutoIncrement: false,
	}

	// Parse type with length
	if strings.Contains(column.Type, "(") {
		typeRegex := regexp.MustCompile(`([A-Z]+)\((\d+)(?:,\s*(\d+))?\)`)
		typeMatches := typeRegex.FindStringSubmatch(column.Type)
		if len(typeMatches) >= 3 {
			column.Type = typeMatches[1]
			if length, err := strconv.Atoi(typeMatches[2]); err == nil {
				column.Length = &length
			}
			if len(typeMatches) >= 4 && typeMatches[3] != "" {
				if scale, err := strconv.Atoi(typeMatches[3]); err == nil {
					column.Scale = &scale
				}
			}
		}
	}

	// Handle PostgreSQL specific types
	switch column.Type {
	case "BIGSERIAL":
		column.AutoIncrement = true
	case "SERIAL":
		column.AutoIncrement = true
	case "SMALLSERIAL":
		column.AutoIncrement = true
	}

	// Parse constraints
	if len(matches) > 3 {
		constraints := strings.ToUpper(matches[3])

		if strings.Contains(constraints, "NOT NULL") {
			column.NotNull = true
		}
		if strings.Contains(constraints, "UNIQUE") {
			column.Unique = true
		}

		// Parse DEFAULT value
		defaultRegex := regexp.MustCompile(`(?i)DEFAULT\s+([^,\s]+(?:\s+[^,\s]+)*)`)
		defaultMatches := defaultRegex.FindStringSubmatch(matches[3])
		if len(defaultMatches) >= 2 {
			defaultVal := strings.TrimSpace(defaultMatches[1])
			column.DefaultValue = &defaultVal
		}
	}

	return column, nil
}

// isConstraint checks if an item is a constraint definition
func (p *PostgreSQLParser) isConstraint(item string) bool {
	constraintKeywords := []string{"CONSTRAINT", "PRIMARY KEY", "FOREIGN KEY", "CHECK", "UNIQUE"}
	itemUpper := strings.ToUpper(strings.TrimSpace(item))

	for _, keyword := range constraintKeywords {
		if strings.HasPrefix(itemUpper, keyword) {
			return true
		}
	}
	return false
}

// parseConstraint parses a constraint definition
func (p *PostgreSQLParser) parseConstraint(table *Table, constraintDef string, options ParseOptions) error {
	constraintUpper := strings.ToUpper(strings.TrimSpace(constraintDef))

	// Parse PRIMARY KEY
	if strings.Contains(constraintUpper, "PRIMARY KEY") {
		pkRegex := regexp.MustCompile(`(?i)(?:CONSTRAINT\s+\w+\s+)?PRIMARY\s+KEY\s*\(([^)]+)\)`)
		matches := pkRegex.FindStringSubmatch(constraintDef)
		if len(matches) >= 2 {
			columns := strings.Split(matches[1], ",")
			for _, col := range columns {
				table.PrimaryKey = append(table.PrimaryKey, strings.TrimSpace(col))
			}
		}
		return nil
	}

	// Parse FOREIGN KEY
	if strings.Contains(constraintUpper, "FOREIGN KEY") {
		fkRegex := regexp.MustCompile(`(?i)CONSTRAINT\s+(\w+)\s+FOREIGN\s+KEY\s*\(([^)]+)\)\s+REFERENCES\s+(\w+)\s*\(([^)]+)\)`)
		matches := fkRegex.FindStringSubmatch(constraintDef)
		if len(matches) >= 5 {
			fk := ForeignKey{
				Name:              matches[1],
				Columns:           strings.Split(strings.ReplaceAll(matches[2], " ", ""), ","),
				ReferencedTable:   matches[3],
				ReferencedColumns: strings.Split(strings.ReplaceAll(matches[4], " ", ""), ","),
			}
			table.ForeignKeys = append(table.ForeignKeys, fk)
		}
		return nil
	}

	// For now, ignore other constraints
	if options.IgnoreUnsupported {
		return nil
	}

	return fmt.Errorf("unsupported constraint: %s", constraintDef)
}

// splitTableItems splits table body into individual items (columns and constraints)
func (p *PostgreSQLParser) splitTableItems(body string) []string {
	items := []string{}
	current := ""
	parenDepth := 0
	inString := false
	stringChar := byte(0)

	for i := 0; i < len(body); i++ {
		char := body[i]

		if !inString {
			if char == '\'' || char == '"' {
				inString = true
				stringChar = char
			} else if char == '(' {
				parenDepth++
			} else if char == ')' {
				parenDepth--
			} else if char == ',' && parenDepth == 0 {
				if strings.TrimSpace(current) != "" {
					items = append(items, strings.TrimSpace(current))
				}
				current = ""
				continue
			}
		} else {
			if char == stringChar && (i == 0 || body[i-1] != '\\') {
				inString = false
				stringChar = 0
			}
		}

		current += string(char)
	}

	// Add the last item
	if strings.TrimSpace(current) != "" {
		items = append(items, strings.TrimSpace(current))
	}

	return items
}

// splitStatements splits SQL content into individual statements
// This is a simple implementation that splits on semicolons
func (p *PostgreSQLParser) splitStatements(content string) []string {
	// Remove SQL comments (-- style)
	commentRegex := regexp.MustCompile(`--.*$`)
	content = commentRegex.ReplaceAllString(content, "")

	// Split on semicolons, but be careful about semicolons in strings
	statements := []string{}
	current := ""
	inString := false
	stringChar := byte(0)

	for i := 0; i < len(content); i++ {
		char := content[i]

		if !inString {
			if char == '\'' || char == '"' {
				inString = true
				stringChar = char
			} else if char == ';' {
				if strings.TrimSpace(current) != "" {
					statements = append(statements, current)
				}
				current = ""
				continue
			}
		} else {
			if char == stringChar && (i == 0 || content[i-1] != '\\') {
				inString = false
				stringChar = 0
			}
		}

		current += string(char)
	}

	// Add the last statement if it doesn't end with semicolon
	if strings.TrimSpace(current) != "" {
		statements = append(statements, current)
	}

	return statements
}
