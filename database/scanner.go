package database

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type MigrationFile struct {
	Timestamp string
	Name      string
	FullPath  string
}

func ScanMigrationFiles(directory string) ([]MigrationFile, error) {
	var migrations []MigrationFile

	//Read all .go files in migrations directory
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		//Skip non .go files
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") || strings.HasSuffix(file.Name(), "_test.go") {
			continue
		}

		//Parse file name format: YYYYMMDDHHMMSS-name.go
		fileName := file.Name()
		trimmedName := strings.TrimSuffix(fileName, ".go")
		parts := strings.Split(strings.TrimSuffix(fileName, ".go"), "-")
		if len(parts) < 2 {
			continue
		}

		// Load and compile the migration file
		filePath := filepath.Join(directory, fileName)
		cmd := exec.Command("go", "run", filePath)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: Failed to load migration %s: %v\n", fileName, err)
		}

		migrations = append(migrations, MigrationFile{
			Timestamp: parts[0],
			Name:      trimmedName,
			FullPath:  filepath.Join(directory, fileName),
		})
	}

	//Sort by timestamp
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Timestamp < migrations[j].Timestamp
	})

	return migrations, nil
}

func GetPendingMigrations(db *sql.DB, files []MigrationFile) ([]MigrationFile, error) {
	var pending []MigrationFile

	for _, file := range files {
		//Check if migration is pending
		var exists bool
		err := db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM migrations WHERE name = ?)",
			file.Name,
		).Scan(&exists)

		if err != nil {
			return nil, err
		}

		if !exists {
			pending = append(pending, file)
		}
	}

	return pending, nil
}
