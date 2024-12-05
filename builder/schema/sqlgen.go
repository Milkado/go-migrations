package schema

import (
	"fmt"
	"strings"
)

type (
	SQLBuilder struct {
		table      string
		columns    []Column
		indexes    []Index
		references []References
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

func (b *SQLBuilder) Build() string {
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
		if !col.Nullable {
			def += " NOT NULL"
		}

		if col.Default != nil {
			def += fmt.Sprintf(" DEFAULT %s", *col.Default)
		}

		if col.AutoInc {
			def += " AUTO_INCREMENT"
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
