package database

import (
	"database/sql"
	"fmt"

	"github.com/Milkado/go-migrations/cmd"
)

var createTalbeSQL = map[string]string{
	"mysql": `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			batch INT NOT NULL,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`,
	"postgres": `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			batch INT NOT NULL,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`,
	"sqlite3": `
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			batch INT NOT NULL,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`,
}

func CreateMigrationsTable(db *sql.DB, driver string) error {
	sql, ok := createTalbeSQL[driver]
	if !ok {
		return fmt.Errorf("driver not supported")
	}

	// Add logging to debug
	fmt.Println(cmd.Yellow + "Creating migrations table..." + cmd.Reset)

	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("error creating migrations table, error: %v", err)
	}

	fmt.Println(cmd.Green + "Migrations table created successfully" + cmd.Reset)
	return nil
}
