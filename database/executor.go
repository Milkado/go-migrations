package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/Milkado/go-migrations/cmd"
)

type (
	Executor struct {
		db *sql.DB
	}
)

func NewExecutor(db *sql.DB) *Executor {
	return &Executor{db: db}
}

func (e *Executor) Migrate() error {
	//Get absolute path
	absPath, err := filepath.Abs("./database/migrations")
	if err != nil {
		return err
	}
	//Scan for files
	files, err := ScanMigrationFiles(absPath)
	if err != nil {
		return err
	}

	//Get pending migrations
	pending, err := GetPendingMigrations(e.db, files)
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		log.Println(cmd.Green + "No pending migrations" + cmd.Reset)
		return nil
	}

	//Start transaction
	tx, err := e.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//Get current batch
	var batch int
	err = tx.QueryRow("select coalesce(max(batch), 0) from migrations").Scan(&batch)
	if err != nil {
		return err
	}
	batch++

	//Execute pending migrations
	for _, file := range pending {
		log.Printf(cmd.Yellow+"Applying migration %s"+cmd.Reset, file.Name)

		migration, ok := GetMigration(file.Name)
		if !ok {
			return fmt.Errorf("migration %s not found", file.Name)
		}

		//Execute migration
		if err := migration.Up(tx); err != nil {
			return err
		}

		//Record migration
		_, err = tx.Exec("insert into migrations (name, batch) values (?, ?)", file.Name, batch)
		if err != nil {
			return fmt.Errorf("error recording migration %s, error: %v", file.Name, err)
		}

		log.Printf(cmd.Green+"Migrated: %s"+cmd.Reset, file.Name)
	}
	//Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf(cmd.Green+"Successfully migrated %d files\n"+cmd.Reset, len(pending))
	return nil

}

func (e *Executor) Rollback() error {
	//Get last batch
	var batch int
	err := e.db.QueryRow("select coalesce(max(batch), 0) from migrations").Scan(&batch)
	if err != nil {
		return err
	}

	if batch == 0 {
		log.Println(cmd.Green + "No migrations to rollback" + cmd.Reset)
		return nil
	}

	//Get migrations from last batch
	rows, err := e.db.Query("select name from migrations where batch = ?", batch)
	if err != nil {
		return fmt.Errorf("error getting migrations: %v", err)
	}
	defer rows.Close()

	//Start transaction
	tx, err := e.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}

		log.Printf(cmd.Yellow+"Rolling back migration %s"+cmd.Reset, name)

		migration, exists := GetMigration(name)
		if !exists {
			return fmt.Errorf("migration %s not found", name)
		}

		//If exists, run Down()
		if err := migration.Down(tx); err != nil {
			return err
		}

		//Delete from migrations table
		_, err = tx.Exec("delete from migrations where name = ?", name)
		if err != nil {
			return fmt.Errorf("error deleting migration %s, error: %v", name, err)
		}

		log.Printf(cmd.Green+"Rolled back: %s"+cmd.Reset, name)
	}

	//Commit transaction

	return tx.Commit()
}
