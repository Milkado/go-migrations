package migrations

import (
	"database/sql"

	"github.com/Milkado/go-migrations/builder/schema"
	"github.com/Milkado/go-migrations/database"
)

type Migration20241202174839createuserstable struct {
	database.Migration
}

func (m *Migration20241202174839createuserstable) Up(db *sql.Tx) error {
	// TODO: Write your migration here
	// _, err := db.Exec(`
	// CREATE TABLE IF NOT EXISTS users (
	// 	id INT AUTO_INCREMENT PRIMARY KEY,
	// 	name VARCHAR(255) NOT NULL,
	// 	email VARCHAR(255) NOT NULL,
	// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	// )
	// `)
	_, err := db.Exec(schema.Query().Create("users").
		Id().
		String("name", false).
		String("email", false).
		Timestamps().
		Build(),
	)
	return err
}

func (m *Migration20241202174839createuserstable) Down(db *sql.Tx) error {
	// TODO: Write your rollback migration here
	_, err := db.Exec(schema.Query().Drop("users").Build())
	return err
}

func init() {
	database.Register("20241202174839-create-users-table", &Migration20241202174839createuserstable{})
}
