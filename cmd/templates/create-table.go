package templates

const CreateTable = ` package migrations

import (
	"database/sql" 
)

type Migration{{.Timestamp}}{{.Name}} struct {
	db *sql.DB
}

func (m *Migration{{.Timestamp}}{{.Name}}) Up() error {
	// TODO: Write your migration here
	return nil
}

func (m *Migration{{.Timestamp}}{{.Name}}) Down() error {
	// TODO: Write your rollback migration here
	return nil
}

func (m *Migration{{.Timestamp}}{{.Name}}) GetTimestamp() string {
	return "{{.Timestamp}}"
}

func (m *Migration{{.Timestamp}}{{.Name}}) GetName() string {
	return "{{.Name}}"
}
`