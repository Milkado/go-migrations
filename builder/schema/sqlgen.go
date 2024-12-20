package schema

import (
	"fmt"
	"strings"
)

type (
	SQLType    string
	Dialect    string
	SQLBuilder struct {
		sqlType       SQLType
		table         string
		columns       []Column
		dropColumns   []string
		modifyColumns []Column
		renameColumns []RenameColumn
		indexes       []Index
		references    []References
		dialect       Dialect
	}

	Column struct {
		Name       string
		Type       string
		Length     int
		Nullable   bool
		Default    *string
		PrimaryKey bool
		AutoInc    bool
	}

	RenameColumn struct {
		OldName string
		NewName string
	}

	Index struct {
		Name    string
		Columns []string
		Type    string //"index" or "unique"
	}

	References struct {
		Column    string
		RefTable  string
		RefColumn string
		OnDelete  string
		OnUpdate  string
	}
)

const (
	CreateTable         = "CREATE TABLE"
	DropTable           = "DROP TABLE"
	AlterTable          = "ALTER TABLE"
	MySQL       Dialect = "mysql"
	Postgres    Dialect = "postgres"
	SQLite      Dialect = "sqlite"
)

func (b *SQLBuilder) Build() string {
	var sql strings.Builder
	switch b.sqlType {
	case CreateTable:
		return b.buildCreateTable()
	case DropTable:
		return b.buildDropTable()
	case AlterTable:
		return b.buildAlterTable()
	}

	return sql.String()
}

func (b *SQLBuilder) buildCreateTable() string {
	var sql strings.Builder

	sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", b.table))

	//Build columns
	var definitions []string
	for _, col := range b.columns {
		def := col.Name

		//Add length if set
		if col.Length > 0 {
			def += fmt.Sprintf(" %s(%d)", col.Type, col.Length)
		} else {
			def += fmt.Sprintf(" %s", col.Type)
		}

		//Add constraints
		if !col.Nullable && col.Name != "id" {
			def += " NOT NULL"
		}

		if col.Default != nil {
			def += fmt.Sprintf(" DEFAULT %s", *col.Default)
		}

		if col.AutoInc {
			switch b.dialect {
			case MySQL:
				def += " AUTO_INCREMENT"
			case Postgres:
				def += ""
			case SQLite:
				def += " AUTOINCREMENT"
			}
		}

		if col.PrimaryKey {
			def += " PRIMARY KEY"
		}

		definitions = append(definitions, def)
	}

	//Add indexes
	for _, idx := range b.indexes {
		cols := strings.Join(idx.Columns, ", ")

		if idx.Type == "unique" {
			definitions = append(definitions,
				fmt.Sprintf(" UNIQUE KEY %s (%s)", idx.Name, cols),
			)
		} else {
			definitions = append(definitions,
				fmt.Sprintf(" INDEX %s (%s)", idx.Name, cols),
			)
		}
	}

	//Add foreign keys
	for _, ref := range b.references {
		def := fmt.Sprintf(" FOREIGN KEY (%s) REFERENCES %s (%s)",
			ref.Column, ref.RefTable, ref.RefColumn)

		if ref.OnDelete != "" {
			def += fmt.Sprintf(" ON DELETE %s", ref.OnDelete)
		}

		if ref.OnUpdate != "" {
			def += fmt.Sprintf(" ON UPDATE %s", ref.OnUpdate)
		}

		definitions = append(definitions, def)
	}

	sql.WriteString(strings.Join(definitions, ",\n"))
	sql.WriteString("\n)")

	return sql.String()
}

func (b *SQLBuilder) buildDropTable() string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", b.table)
}

func (b *SQLBuilder) buildAlterTable() string {
	//TODO: create builder for alterations

	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("ALTER TABLE %s\n", b.table))

	//Build definitions to add columns/constraints/drops/renames
	var definitions []string
	//Add columns
	for _, col := range b.columns {
		def := fmt.Sprintf("ADD COLUMN %s %s", col.Name, col.Type)

		if col.Length > 0 {
			def += fmt.Sprintf("(%d)", col.Length)
		}

		if !col.Nullable {
			def += " NOT NULL"
		}

		if col.Default != nil {
			def += fmt.Sprintf(" DEFAULT %s", *col.Default)
		}

		definitions = append(definitions, def)
	}

	// Add fk constraints
	for _, ref := range b.references {
		constraintname := ref.Column + "_" + ref.RefTable + "_fk"
		def := fmt.Sprintf("ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
			constraintname, ref.Column, ref.RefTable, ref.RefColumn)

		//Add on delete
		if ref.OnDelete != "" {
			def += fmt.Sprintf(" ON DELETE %s", ref.OnDelete)
		}
		//Add on update
		if ref.OnUpdate != "" {
			def += fmt.Sprintf(" ON UPDATE %s", ref.OnUpdate)
		}

		definitions = append(definitions, def)
	}

	//Add renames
	for _, col := range b.renameColumns {
		definitions = append(definitions,
			fmt.Sprintf("RENAME COLUMN %s TO %s", col.OldName, col.NewName),
		)
	}

	//Add modifications
	for _, col := range b.modifyColumns {
		def := fmt.Sprintf("MODIFY COLUMN %s %s", col.Name, col.Type)

		if col.Length > 0 {
			def += fmt.Sprintf("(%d)", col.Length)

		}

		if !col.Nullable {
			def += " NOT NULL"
		}

		definitions = append(definitions, def)
	}

	//Add uniques and indexes
	for _, idx := range b.indexes {
		cols := strings.Join(idx.Columns, ", ")

		if idx.Type == "unique" {
			definitions = append(definitions,
				fmt.Sprintf("ADD CONSTRAINT %s UNIQUE (%s)", idx.Name, cols),
			)
		} else {
			definitions = append(definitions,
				fmt.Sprintf("ADD INDEX %s (%s)", idx.Name, cols),
			)
		}
	}

	//Add drops
	for _, col := range b.dropColumns {
		definitions = append(definitions,
			fmt.Sprintf("DROP COLUMN %s", col),
		)
	}

	if len(definitions) > 1 {
		sql.WriteString(strings.Join(definitions, ",\n"))
		return sql.String()
	}

	sql.WriteString(strings.Join(definitions, ",\n"))

	return sql.String()
}
