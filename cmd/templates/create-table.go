package templates

const CreateTable = ` package migrations

import (
	"database/sql" 
	
	"github.com/Milkado/go-migrations/database"
)

type Migration{{.Timestamp}}{{.ParsedName}} struct {
	database.Migration
}

func (m *Migration{{.Timestamp}}{{.ParsedName}}) Up(db *sql.Tx) error {
	// TODO: Write your migration here
	return nil
}

func (m *Migration{{.Timestamp}}{{.ParsedName}}) Down(db *sql.Tx) error {
	// TODO: Write your rollback migration here
	return nil
}

func init() {
	database.Register("{{.Timestamp}}-{{.Name}}", &Migration{{.Timestamp}}{{.ParsedName}}{})
}
`