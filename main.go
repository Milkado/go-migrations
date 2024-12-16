package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/Milkado/go-migrations/cmd"
	"github.com/Milkado/go-migrations/config"
	"github.com/Milkado/go-migrations/database"
	_ "github.com/Milkado/go-migrations/database/migrations"
)

func main() {
	var command string
	var name string

	flag.StringVar(&command, "c", "", "Command to execute")
	flag.StringVar(&name, "name", "", "Name of the migration")
	flag.Parse()

	if command == "" {
		errorC()
		return
	}

	switch command {
	case "migration:create":
		cmd.NewMigration(name)
	case "migrate":
		db := dbConnect()
		defer db.Close()
		executor := database.NewExecutor(db)
		if err := executor.Migrate(); err != nil {
			fmt.Println(cmd.Red, err.Error())
			return
		}
	case "migrate:rollback":
		db := dbConnect()
		defer db.Close()
		executor := database.NewExecutor(db)
		if err := executor.Rollback(); err != nil {
			fmt.Println(cmd.Red, err.Error())
			return
		}
	default:
		errorC()
	}
}

func errorC() {
	fmt.Println(cmd.Red, "Command not available or not specified"+cmd.Reset)
	fmt.Println(cmd.Yellow + "Usage: go-migrations --c <command>")
	fmt.Println("Available commands:")
	fmt.Println("migration:create --name <name>")
	fmt.Println("migrate")
	fmt.Println("migrate:rollback")
	fmt.Println(cmd.Reset)
}

func dbConnect() *sql.DB {
	db, err := database.NewConnectionWithMonitoring(&config.Config)
	if err != nil {
		fmt.Println(cmd.Red, err.Error())
		return nil
	}
	return db
}
