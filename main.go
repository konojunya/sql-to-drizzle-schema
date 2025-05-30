// Package main provides the CLI interface for sql-to-drizzle-schema.
//
// This tool converts SQL DDL files (CREATE TABLE statements, etc.) to
// Drizzle ORM schema definitions in TypeScript format.
//
// Usage:
//
//	sql-to-drizzle-schema [SQL_FILE] -o [OUTPUT_FILE]
//
// Example:
//
//	sql-to-drizzle-schema ./schema.sql -o schema.ts
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/konojunya/sql-to-drizzle-schema/internal/generator"
	"github.com/konojunya/sql-to-drizzle-schema/internal/parser"
	"github.com/konojunya/sql-to-drizzle-schema/internal/reader"
	"github.com/spf13/cobra"
)

var (
	// outputFile stores the path for the generated TypeScript file
	outputFile string
	// dialectFlag stores the SQL dialect to use for parsing
	dialectFlag string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sql-to-drizzle-schema [SQL_FILE]",
	Short: "Convert SQL schemas to Drizzle ORM schema definitions",
	Long: `A CLI tool that converts SQL DDL files to Drizzle ORM schema definitions.

This tool reads SQL files containing CREATE TABLE statements and other DDL
commands, then generates equivalent TypeScript code using Drizzle ORM syntax.

Supported SQL features:
- CREATE TABLE statements
- Column definitions with various data types
- Primary keys and foreign keys
- Constraints and indexes
- Default values

Supported database dialects:
- PostgreSQL (default)
- MySQL (planned)
- Spanner (planned)

Example usage:
  sql-to-drizzle-schema ./database.sql -o schema.ts
  sql-to-drizzle-schema ./database.sql --dialect postgresql -o schema.ts
  sql-to-drizzle-schema ./mysql-schema.sql --dialect mysql -o schema.ts`,
	Args: cobra.ExactArgs(1), // Exactly one SQL file argument is required
	Run: func(cmd *cobra.Command, args []string) {
		// Get the SQL file path from command arguments
		sqlFile := args[0]

		// Set default output file if not specified
		if outputFile == "" {
			outputFile = "schema.ts"
		}

		// Parse and validate dialect
		var dialect parser.DatabaseDialect
		switch strings.ToLower(dialectFlag) {
		case "postgresql", "postgres", "pg":
			dialect = parser.PostgreSQL
		case "mysql":
			dialect = parser.MySQL
		case "spanner":
			dialect = parser.Spanner
		default:
			if dialectFlag != "" {
				fmt.Fprintf(os.Stderr, "Unsupported dialect '%s'. Supported dialects: postgresql, mysql, spanner\n", dialectFlag)
				os.Exit(1)
			}
			// Default to PostgreSQL
			dialect = parser.PostgreSQL
		}

		// Display conversion information to user
		fmt.Printf("Converting SQL file: %s\n", sqlFile)
		fmt.Printf("Output file: %s\n", outputFile)
		fmt.Printf("Database dialect: %s\n", dialect)

		// Read the SQL file content
		content, err := reader.ReadSQLFile(sqlFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading SQL file: %v\n", err)
			os.Exit(1)
		}

		// Parse the SQL content
		fmt.Println("Parsing SQL content...")
		parseOptions := parser.DefaultParseOptions()
		parseOptions.Dialect = dialect
		parseResult, err := parser.ParseSQLContent(content, dialect, parseOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing SQL: %v\n", err)
			os.Exit(1)
		}

		// Display parsing results
		fmt.Printf("Successfully parsed %d table(s):\n", len(parseResult.Tables))
		for _, table := range parseResult.Tables {
			fmt.Printf("  - Table: %s (%d columns)\n", table.Name, len(table.Columns))
			for _, column := range table.Columns {
				fmt.Printf("    - %s: %s", column.Name, column.Type)
				if column.Length != nil {
					fmt.Printf("(%d)", *column.Length)
				}
				if column.NotNull {
					fmt.Print(" NOT NULL")
				}
				if column.AutoIncrement {
					fmt.Print(" AUTO_INCREMENT")
				}
				if column.DefaultValue != nil {
					fmt.Printf(" DEFAULT %s", *column.DefaultValue)
				}
				fmt.Println()
			}
			if len(table.PrimaryKey) > 0 {
				fmt.Printf("    Primary Key: %v\n", table.PrimaryKey)
			}
			if len(table.ForeignKeys) > 0 {
				fmt.Printf("    Foreign Keys: %d\n", len(table.ForeignKeys))
			}
		}

		// Display any parsing errors
		if len(parseResult.Errors) > 0 {
			fmt.Printf("\nWarnings during parsing:\n")
			for _, parseErr := range parseResult.Errors {
				fmt.Printf("  - %v\n", parseErr)
			}
		}

		// Generate Drizzle schema
		fmt.Println("\nGenerating Drizzle ORM schema...")
		generatorOptions := generator.DefaultGeneratorOptions()

		err = generator.GenerateSchemaToFile(parseResult.Tables, dialect, outputFile, generatorOptions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating schema: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Successfully generated Drizzle schema: %s\n", outputFile)
		fmt.Printf("📝 Generated %d table definition(s)\n", len(parseResult.Tables))
	},
}

// init initializes the CLI flags and configuration
func init() {
	// Add the output flag with short (-o) and long (--output) forms
	// If not specified, the default "schema.ts" will be used
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output TypeScript file (default: schema.ts)")

	// Add the dialect flag with short (-d) and long (--dialect) forms
	// If not specified, PostgreSQL will be used as default
	rootCmd.Flags().StringVarP(&dialectFlag, "dialect", "d", "", "Database dialect (postgresql, mysql, spanner) (default: postgresql)")
}

// main is the entry point of the application
func main() {
	// Execute the root command and handle any errors
	if err := rootCmd.Execute(); err != nil {
		// Print error to stderr and exit with non-zero status
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
