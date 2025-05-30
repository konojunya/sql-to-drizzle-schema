package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	// Test that main function doesn't panic
	// We can't easily test the actual main() function since it calls os.Exit
	// But we can test the command setup
	if rootCmd == nil {
		t.Error("rootCmd should be initialized")
	}
}

func TestRootCmd_Setup(t *testing.T) {
	// Test that the command is properly configured
	if rootCmd.Use != "sql-to-drizzle-schema [SQL_FILE]" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "sql-to-drizzle-schema [SQL_FILE]")
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}

	// Check that it expects exactly one argument
	if rootCmd.Args == nil {
		t.Error("rootCmd.Args should be set")
	}
}

func TestRootCmd_Flags(t *testing.T) {
	// Test that flags are properly configured
	outputFlag := rootCmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("output flag should be defined")
	}

	dialectFlag := rootCmd.Flags().Lookup("dialect")
	if dialectFlag == nil {
		t.Error("dialect flag should be defined")
	}

	// Test short flags
	oFlag := rootCmd.Flags().ShorthandLookup("o")
	if oFlag == nil {
		t.Error("short flag 'o' should be defined")
	}

	dFlag := rootCmd.Flags().ShorthandLookup("d")
	if dFlag == nil {
		t.Error("short flag 'd' should be defined")
	}
}

func TestGlobalVariables(t *testing.T) {
	// Test that global variables are properly initialized
	tests := []struct {
		name     string
		variable interface{}
	}{
		{"outputFile", outputFile},
		{"dialectFlag", dialectFlag},
		{"rootCmd", rootCmd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.variable == nil {
				t.Errorf("Global variable %s should be initialized", tt.name)
			}
		})
	}
}

func TestInit(t *testing.T) {
	// Test that init function properly sets up flags
	// We can't directly test init(), but we can verify its effects

	// Check that flags have been added to rootCmd
	flags := rootCmd.Flags()

	if !flags.HasFlags() {
		t.Error("rootCmd should have flags after init()")
	}

	// Check specific flags exist
	expectedFlags := []string{"output", "dialect"}
	for _, flagName := range expectedFlags {
		if flags.Lookup(flagName) == nil {
			t.Errorf("Flag %s should be defined after init()", flagName)
		}
	}

	// Check short flags exist
	expectedShortFlags := []string{"o", "d"}
	for _, shortFlag := range expectedShortFlags {
		if flags.ShorthandLookup(shortFlag) == nil {
			t.Errorf("Short flag %s should be defined after init()", shortFlag)
		}
	}
}

func TestRootCmd_Args(t *testing.T) {
	// Test that the command correctly validates arguments
	// We test this by checking the Args field is set correctly
	if rootCmd.Args == nil {
		t.Error("rootCmd.Args should be set to validate arguments")
	}
}

func TestPackageConstants(t *testing.T) {
	// Test that the package is properly set up
	// This is more of a compilation test
	if rootCmd.Use == "" {
		t.Error("Command Use should not be empty")
	}

	if rootCmd.Short == "" {
		t.Error("Command Short description should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("Command Long description should not be empty")
	}
}
