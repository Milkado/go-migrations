package cmd

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type templateData struct {
	Timestamp string
	Name      string
	ParsedName string
}

func validateMigrationName(name string) error {
	re := `^[a-z0-9]+(-[a-z0-9]+)*$`
	matched, err := regexp.MatchString(re, name)
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("invalid migration name. Use only lowercase letters, numbers and hyphens as separators. Example: create-users-table")
	}

	return nil
}

func GenerateMigrationFile(name string, tmplfile string) {
	name = strings.ReplaceAll(name, " ", "-")
	err := validateMigrationName(name)
	if err != nil {
		fmt.Println(Red, "Error validating migration name, error:")
		fmt.Println(err.Error() + Reset)
		return
	}

	timestamp := time.Now().Format("20060102150405") // Generates a timestamp YYYYMMDDHHMMSS
	parsedname := strings.ReplaceAll(name, "-", "")
	data := templateData{
		Timestamp: timestamp,
		Name:      name,
		ParsedName: parsedname,
	}

	//create template
	tmpl, err := template.New("migration").Parse(tmplfile)
	if err != nil {
		fmt.Println(Red, "Error parsing template, error:")
		fmt.Println(err.Error() + Reset)
		return
	}

	//check if folder migrations exists
	if _, err := os.Stat("database/migrations"); os.IsNotExist(err) {
		os.Mkdir("database/migrations", os.ModePerm)
	}

	// Execute the template
	var bytes bytes.Buffer
	err = tmpl.Execute(&bytes, data)
	if err != nil {
		fmt.Println(Red, "Error executing template, error:")
		fmt.Println(err.Error() + Reset)
		return
	}

	filename := "database/migrations/" + timestamp + "-" + name + ".go"

	//write file from buffer
	err = os.WriteFile(filename, bytes.Bytes(), 0644)
	if err != nil {
		fmt.Println(Red, "Error writing file, error:")
		fmt.Println(err.Error() + Reset)
		return
	}

	fmt.Println(Green + "Migration file created successfully: " + filename + Reset)
}
