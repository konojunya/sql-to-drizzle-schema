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

	"github.com/konojunya/sql-to-drizzle-schema/internal/reader"
	"github.com/spf13/cobra"
)

var (
	// outputFile stores the path for the generated TypeScript file
	outputFile string
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

Example usage:
  sql-to-drizzle-schema ./database.sql -o schema.ts
  sql-to-drizzle-schema ./migrations/*.sql --output drizzle-schema.ts`,
	Args: cobra.ExactArgs(1), // Exactly one SQL file argument is required
	Run: func(cmd *cobra.Command, args []string) {
		// Get the SQL file path from command arguments
		sqlFile := args[0]

		// Set default output file if not specified
		if outputFile == "" {
			outputFile = "schema.ts"
		}

		// Display conversion information to user
		fmt.Printf("Converting SQL file: %s\n", sqlFile)
		fmt.Printf("Output file: %s\n", outputFile)

		// Read the SQL file content
		content, err := reader.ReadSQLFile(sqlFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading SQL file: %v\n", err)
			os.Exit(1)
		}

		// Display the SQL content for debugging purposes
		// TODO: Remove this in production version
		fmt.Printf("SQL content:\n%s\n", content)

		// TODO: Implement SQL parsing and Drizzle schema generation
		// This will include:
		// 1. Parsing SQL DDL statements
		// 2. Converting SQL types to Drizzle types
		// 3. Generating TypeScript code with proper imports
		// 4. Writing the output file
		fmt.Println("Conversion logic will be implemented here...")
	},
}

// init initializes the CLI flags and configuration
func init() {
	// Add the output flag with short (-o) and long (--output) forms
	// If not specified, the default "schema.ts" will be used
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output TypeScript file (default: schema.ts)")
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
