package main

import (
	"flag"
	"fmt"

	"github.com/Milkado/go-migrations/cmd"
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
	default:
		errorC()
	}
}

func errorC() {
	fmt.Println(cmd.Red, "Command not available or not specified" + cmd.Reset)
	fmt.Println(cmd.Yellow + "Usage: go-migrations --c <command>")
	fmt.Println("Available commands:")
	fmt.Println("migration:create --name <name>")
	fmt.Println(cmd.Reset)
}