 package migrations

import (
	"database/sql" 
	
	"github.com/Milkado/go-migrations/database"
	"github.com/Milkado/go-migrations/builder/schema"
)

type Migration20241212152938createpasswordresettable struct {
	database.Migration
}

func (m *Migration20241212152938createpasswordresettable) Up(db *sql.Tx) error {
	
	_, err := db.Exec(
		schema.Query().Create("password_reset").
			Id().
			String("token", false, 36).
			ForeignId("user_id", false).Refenreces("id").On("users").
			Timestamps().
			Build(),
	)
	
	return err
}

func (m *Migration20241212152938createpasswordresettable) Down(db *sql.Tx) error {
	
	_, err := db.Exec(
		schema.Query().Drop("password_reset").Build(),
	)
	
	return err
}

func init() {
	database.Register("20241212152938-create-password-reset-table", &Migration20241212152938createpasswordresettable{})
}
