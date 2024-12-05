package cmd_test

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Milkado/go-migrations/cmd"
	"github.com/Milkado/go-migrations/database"
)

func TestCreateMigration(t *testing.T) {
	tests := []struct {
		name         string
		newMigration func()
		expected     string
	}{
		{
			name: "Create migration with valid name",
			newMigration: func() {
				cmd.NewMigration("create-posts-table")
			},
			expected: "create-posts-table",
		},
		// {
		// 	name: "Create alter migration with valid name",
		// 	newMigration: func() {
		// 		cmd.NewMigration("alter-posts-table")
		// 	},
		// 	expected: "alter-posts-table",
		// },
	}

	// Get and store current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting current directory: %v", err)
	}
	defer os.Chdir(originalDir) // Restore original directory when done

	// Change to project root
	projectRoot := filepath.Dir(filepath.Dir(originalDir))
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("error changing to project root: %v", err)
	}

	// Add debug logging
	log.Printf(cmd.Blue+"DEBUG: Current dir: %s"+cmd.Reset, projectRoot)

	// Ensure migrations directory exists
	migrationsDir := filepath.Join(projectRoot, "database", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		t.Fatalf("error creating migrations directory: %v", err)
	}
	log.Printf(cmd.Blue+"DEBUG: migrationsDir: %s"+cmd.Reset, migrationsDir)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.newMigration()
			files, err := database.ScanMigrationFiles("database/migrations")
			if err != nil {
				t.Errorf("error scanning migration files: %v", err)
			}
			//Check if expected exists in file array
			found := false
			for _, file := range files {
				if strings.Contains(file.Name, tt.expected) {
					found = true
				}
			}

			if !found {
				t.Errorf(cmd.Red+"expected file containing %s to exist"+cmd.Reset, tt.expected)
			}
		})
	}

	cleanup(t, projectRoot)
}

func cleanup(t *testing.T, projectRoot string) {
	t.Helper()

	cleanupDir := filepath.Join(projectRoot, "database", "migrations")

	files, err := database.ScanMigrationFiles(cleanupDir)
	if err != nil {
		t.Logf("error scanning files during cleanup: %v", err)
		return
	}

	for _, file := range files {
		if strings.Contains(file.Name, "posts") {
			// Try both file.Name and file.FullPath
			fullPath := filepath.Join(cleanupDir, file.Name)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				// Try alternative path
				fullPath = file.FullPath
			}
			if err := os.Remove(fullPath); err != nil {
				t.Logf(cmd.Red+"error removing file: %s"+cmd.Reset, fullPath)
			} else {
				t.Logf(cmd.Green+"removed file: %s"+cmd.Reset, fullPath)
			}
		}
	}
}
