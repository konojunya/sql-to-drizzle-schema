# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project that converts SQL schemas to Drizzle ORM schema definitions. The project uses Go 1.24.1 and implements a CLI tool using the Cobra framework.

**Purpose**: Convert SQL DDL files (CREATE TABLE statements, etc.) to TypeScript Drizzle ORM schema definitions.

**Target Users**: Developers migrating from traditional SQL schemas to Drizzle ORM, or those who prefer defining schemas in SQL first.

## Architecture

The project follows Go best practices with a clean package structure:

```
sql-to-drizzle-schema/
├── main.go                    # CLI entry point using Cobra
├── internal/                  # Internal packages (not importable by external projects)
│   ├── reader/               # File reading utilities
│   │   └── file.go           # SQL file reading functionality
│   ├── parser/               # SQL parsing functionality
│   │   ├── types.go          # Type definitions for parsed SQL structures
│   │   ├── postgres.go       # PostgreSQL-specific parser implementation
│   │   └── parser.go         # Parser factory and common functionality
│   └── generator/            # Drizzle schema generation functionality
│       ├── types.go          # Type definitions for schema generation
│       ├── postgres.go       # PostgreSQL to Drizzle type mapping and generation
│       └── generator.go      # Generator factory and file operations
├── example/                  # Example SQL files for testing
│   └── postgres/
│       └── create-table.sql  # PostgreSQL example schema
├── doc.go                    # Package-level documentation
├── go.mod                    # Go module definition
├── go.sum                    # Go dependencies lock file
└── CLAUDE.md                 # This file
```

### Package Structure

- **main**: CLI interface using Cobra, handles command-line arguments and orchestrates the conversion process
- **internal/reader**: File I/O operations for reading SQL files with proper error handling
- **internal/parser**: SQL parsing functionality with support for PostgreSQL (extensible for MySQL/Spanner)
  - **types.go**: Type definitions for parsed SQL structures (Table, Column, Constraint, etc.)
  - **postgres.go**: PostgreSQL-specific parser using regex-based parsing
  - **parser.go**: Parser factory and common functionality
- **internal/generator**: Drizzle ORM schema generation functionality
  - **types.go**: Type definitions for schema generation (GeneratorOptions, DrizzleType, etc.)
  - **postgres.go**: PostgreSQL to Drizzle type mapping and TypeScript code generation
  - **generator.go**: Generator factory and file operations
- **example**: Sample SQL files for testing and documentation purposes

### Dependencies

- `github.com/spf13/cobra`: CLI framework for building command-line applications
- Standard library packages: `fmt`, `os`, `io`, `regexp`, `strings` for basic operations

## Common Commands

### Using Makefile (Recommended)

```bash
# Show all available commands
make help

# Development workflow
make dev                    # Quick development cycle (format, test, build)
make check                  # Run all checks (format, vet, lint, test)
make ci                     # Run CI pipeline

# Building
make build                  # Build the binary
make build-all             # Build for multiple platforms
make install               # Install to GOPATH/bin

# Testing
make test                   # Run all tests
make test-coverage         # Run tests with coverage
make test-coverage-view    # Run tests with coverage and open in browser
make test-verbose          # Run tests with verbose output
make bench                 # Run benchmarks

# Code quality
make fmt                   # Format code
make lint                  # Run linter (requires golangci-lint)
make vet                   # Run go vet
make security              # Run security scan (requires gosec)

# Dependencies
make deps                  # Download and verify dependencies
make tidy                  # Tidy dependencies

# Examples and debugging
make example               # Build and run with example file
make debug-example         # Run example with verbose output

# Utilities
make clean                 # Clean build artifacts
make setup                 # Setup development environment
make docs                  # Generate Go documentation
make status                # Show project status
```

### Direct Go Commands

```bash
# Build the project
go build -o sql-to-drizzle-schema

# Build and run with example file
go build -o sql-to-drizzle-schema && ./sql-to-drizzle-schema ./example/postgres/create-table.sql -o output.ts

# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -run TestFunctionName ./path/to/package

# Format code (always run before committing)
go fmt ./...

# Run linter (if golangci-lint is available)
golangci-lint run

# Install dependencies
go mod tidy

# Download dependencies
go mod download

# Generate documentation
godoc -http=:6060
```

## CLI Usage

```bash
# Convert SQL file to Drizzle schema
./sql-to-drizzle-schema input.sql -o schema.ts

# Use default output filename (schema.ts)
./sql-to-drizzle-schema input.sql

# Get help
./sql-to-drizzle-schema --help
```

## Development Guidelines

### Code Style
- Follow standard Go conventions (gofmt, golint)
- Use meaningful variable and function names
- Add comprehensive comments for all exported functions
- Include package-level documentation
- Use error wrapping with `fmt.Errorf` and `%w` verb

### Testing
- Write tests for all public functions
- Use table-driven tests where appropriate
- Include both positive and negative test cases
- Test error conditions thoroughly

### Documentation
- Document all exported functions with examples
- Keep CLAUDE.md updated with architectural changes
- Use godoc-compatible comments
- Include usage examples in function documentation

### Error Handling
- Always handle errors explicitly
- Use wrapped errors for better debugging
- Provide context in error messages
- Fail fast with meaningful error messages

## Current Status

The project has reached a functional state with complete PostgreSQL support:
- ✅ CLI framework using Cobra with dialect selection (--dialect flag)
- ✅ File reading functionality with error handling
- ✅ Package structure and comprehensive documentation
- ✅ PostgreSQL SQL parsing (CREATE TABLE statements)
  - ✅ Column definitions with types, constraints, defaults
  - ✅ Primary key constraints
  - ✅ Foreign key constraints (basic support)
  - ✅ PostgreSQL-specific types (BIGSERIAL, TIMESTAMP WITH TIME ZONE, etc.)
- ✅ Drizzle ORM schema generation for PostgreSQL
  - ✅ Complete type mapping (BIGSERIAL → bigserial, VARCHAR → varchar, etc.)
  - ✅ Constraint mapping (NOT NULL → .notNull(), DEFAULT → .default(), etc.)
  - ✅ Naming convention support (camelCase, PascalCase, snake_case)
  - ✅ TypeScript code generation with proper imports
- ✅ TypeScript output generation with formatted code
- ✅ Complete end-to-end conversion pipeline
- ✅ Foreign key relationships in generated schema
  - ✅ Automatic .references() generation for foreign key columns
  - ✅ Table dependency sorting for proper declaration order
  - ✅ Support for single-column foreign keys
- ✅ Comprehensive test suite
  - ✅ Unit tests for all internal packages (parser, generator, reader)
  - ✅ Integration tests for end-to-end conversion
  - ✅ Test coverage: reader (100%), parser (83.7%), generator (78.2%)
  - ✅ Error handling and edge case testing
  - ✅ Naming convention testing
  - ✅ Foreign key dependency ordering tests
- 🚧 MySQL parser (planned)
- 🚧 Spanner parser (planned)
- 🚧 Multi-column foreign keys (planned)

## CI/CD Pipeline

The project uses GitHub Actions for comprehensive CI/CD automation:

### Automated Workflows

1. **CI Pipeline** (`.github/workflows/ci.yaml`):
   - **Linting & Formatting**: golangci-lint, go fmt, go vet
   - **Security Scanning**: gosec security analysis
   - **Cross-platform Testing**: Linux, macOS, Windows
   - **Multi-version Testing**: Go 1.23, 1.24.1
   - **Integration Testing**: End-to-end conversion validation
   - **Build Verification**: Cross-platform binary compilation
   - **Vulnerability Scanning**: govulncheck for dependencies

2. **Release Automation** (`.github/workflows/release.yaml`):
   - **Automated Releases**: Triggered by version tags (v*.*.*)
   - **Cross-platform Binaries**: Linux, macOS, Windows (AMD64, ARM64)
   - **Release Notes**: Auto-generated changelog
   - **Go Module Publishing**: Automatic proxy cache updates
   - **Documentation Updates**: Installation instructions

3. **Dependency Management**:
   - **Dependabot**: Weekly dependency updates
   - **Security Alerts**: Automated vulnerability detection
   - **Auto-merging**: Compatible with automated PR merging

### CI Commands Integration

The CI pipeline leverages the Makefile targets:

```bash
# CI pipeline commands used
make ci                     # Full CI pipeline
make deps                   # Dependency management
make check                  # All quality checks
make test-coverage         # Coverage reporting
make build-all             # Cross-platform builds
make security              # Security scanning
```

### Development Integration

- **Pre-commit Checks**: Use `make check` before pushing
- **Local CI Simulation**: Run `make ci` locally
- **Release Preparation**: Use `make release-prep`

### CI/CD Security

For security reasons, the CI pipeline has different behaviors based on the contributor:

1. **Repository Owner/Maintainers**: CI runs automatically on all pushes and PRs
2. **External Contributors**: CI is skipped by default and must be triggered manually
3. **Manual CI Trigger**: Maintainers can comment `/run-ci` on a PR to run the full CI pipeline

This prevents external contributors from consuming excessive CI resources while still allowing thorough testing of contributions after review.

## Future Enhancements

- Support for multiple SQL dialects (PostgreSQL, MySQL, SQLite)
- Advanced SQL features (indexes, triggers, views)
- Configuration file support
- Batch processing of multiple files
- Integration with existing Drizzle projects
- Plugin system for custom transformations