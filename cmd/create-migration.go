package cmd

import (
	"strings"

	"github.com/Milkado/go-migrations/cmd/templates"
)

func NewMigration(name string) {
	if strings.Contains(name, "create") {
		migrationNewTable(name)
		return
	}

	migrationAlterTable(name)
}
func migrationNewTable(name string) {
	GenerateMigrationFile(name, templates.CreateTable)
}

func migrationAlterTable(name string) {
	// GenerateMigrationFile(name, templates.AlterTable)
}
