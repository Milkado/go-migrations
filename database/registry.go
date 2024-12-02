package database

import (
	"database/sql"
	"fmt"
)

type Migration interface {
	Up(db *sql.Tx) error
	Down(db *sql.Tx) error
}

var registeredMigrations = make(map[string]Migration)

func Register(name string, m Migration) {
	fmt.Printf("DEBUG: Registering migration: %s\n", name)
	registeredMigrations[name] = m
}

func GetMigration(name string) (Migration, bool) {
	m, ok := registeredMigrations[name]
	return m, ok
}
