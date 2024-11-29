 package migrations

import (
	"database/sql" 
)

type Migration20241129185353createuserstable struct {
	db *sql.DB
}

func (m *Migration20241129185353createuserstable) Up() error {
	// TODO: Write your migration here
	return nil
}

func (m *Migration20241129185353createuserstable) Down() error {
	// TODO: Write your rollback migration here
	return nil
}

func (m *Migration20241129185353createuserstable) GetTimestamp() string {
	return "20241129185353"
}

func (m *Migration20241129185353createuserstable) GetName() string {
	return "createuserstable"
}
