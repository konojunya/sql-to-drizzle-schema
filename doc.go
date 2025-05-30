// Package sql-to-drizzle-schema provides a command-line tool for converting SQL DDL files
// to Drizzle ORM schema definitions in TypeScript.
//
// # Overview
//
// This tool bridges the gap between traditional SQL schema definitions and modern
// TypeScript ORM solutions by automatically converting CREATE TABLE statements and
// other DDL commands into equivalent Drizzle ORM schema code.
//
// # Features
//
// The tool currently supports:
//   - Reading SQL files from the filesystem
//   - CLI interface with configurable output options
//   - Error handling with detailed error messages
//
// Planned features include:
//   - Parsing various SQL DDL statements
//   - Converting SQL data types to Drizzle ORM types
//   - Generating TypeScript code with proper imports
//   - Support for multiple SQL dialects (PostgreSQL, MySQL, SQLite)
//
// # Usage
//
// The tool is designed to be used as a command-line utility:
//
//	sql-to-drizzle-schema input.sql -o schema.ts
//
// This will read the SQL file 'input.sql' and generate a TypeScript file 'schema.ts'
// containing the equivalent Drizzle ORM schema definitions.
//
// # Architecture
//
// The project is organized into the following packages:
//
//   - main: CLI interface and command orchestration
//   - internal/reader: File I/O operations for SQL files
//
// This structure follows Go best practices by keeping internal packages private
// and providing a clean public API through the main command.
//
// # Example
//
// Given a SQL file with the following content:
//
//	CREATE TABLE users (
//	    id BIGSERIAL PRIMARY KEY,
//	    name VARCHAR(255) NOT NULL,
//	    email VARCHAR(255) NOT NULL UNIQUE,
//	    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
//	);
//
// The tool will generate equivalent Drizzle ORM schema code:
//
//	import { pgTable, bigserial, varchar, timestamp } from 'drizzle-orm/pg-core';
//
//	export const users = pgTable('users', {
//	  id: bigserial('id', { mode: 'number' }).primaryKey(),
//	  name: varchar('name', { length: 255 }).notNull(),
//	  email: varchar('email', { length: 255 }).notNull().unique(),
//	  createdAt: timestamp('created_at', { withTimezone: true }).defaultNow(),
//	});
//
// # Error Handling
//
// The tool provides comprehensive error handling for common scenarios:
//   - File not found or permission errors
//   - Invalid SQL syntax (planned)
//   - Unsupported SQL features (planned)
//
// All errors include context information to help users identify and resolve issues.
//
// # Development
//
// For development guidelines, architectural decisions, and contribution information,
// see the CLAUDE.md file in the project root.
package main
