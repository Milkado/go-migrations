package templates

const CreateTable = ` package migrations

import (
	"database/sql" 
	
	"github.com/Milkado/go-migrations/database"
	"github.com/Milkado/go-migrations/builder/schema"
)

type Migration{{.Timestamp}}{{.ParsedName}} struct {
	database.Migration
}

func (m *Migration{{.Timestamp}}{{.ParsedName}}) Up(db *sql.Tx) error {
	{{if .TableName}}
	_, err := db.Exec(
		schema.Query().Create("{{.TableName}}").Build(),
	)
	{{else}}
	//Write your migration here
	{{end}}
	return err
}

func (m *Migration{{.Timestamp}}{{.ParsedName}}) Down(db *sql.Tx) error {
	{{if .TableName}}
	_, err := db.Exec(
		schema.Query().Drop("{{.TableName}}").Build(),
	)
	{{else}}
	//Write your migration here
	{{end}}
	return err
}

func init() {
	database.Register("{{.Timestamp}}-{{.Name}}", &Migration{{.Timestamp}}{{.ParsedName}}{})
}
`
