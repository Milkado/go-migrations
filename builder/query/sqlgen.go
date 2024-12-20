package query

import (
	"fmt"
	"strings"
)

type (
	Dialect   string
	Generator string

	SQLBuilder struct {
		table     string
		dialect   Dialect
		generator Generator
		insert    struct {
			columns []string
			values  []ValuesGen
		}
		update struct {
			columns []string
			values  []ValuesGen
		}
		where struct {
			column   string
			value    string
			operator string
		}
	}

	ValuesGen struct {
		value string
		typer string
	}
)

const (
	MySQL    Dialect   = "mysql"
	Postgres Dialect   = "postgres"
	SQLite   Dialect   = "sqlite"
	Insert   Generator = "insert"
	Update   Generator = "update"
	Delete   Generator = "delete"
	Select   Generator = "select"
)

func (b *SQLBuilder) Build() string {
	var sql strings.Builder
	//Build query
	switch b.generator {
	case Insert:
		sql.WriteString(b.buildInsert())
	case Update:
		sql.WriteString(b.buildUpdate())
	case Delete:
		sql.WriteString(b.buildDelete())
	case Select:
		break
	}
	return sql.String()
}

func (b *SQLBuilder) buildInsert() string {
	var sql strings.Builder

	sql.WriteString(fmt.Sprintf("INSERT INTO %s ", b.table))

	//Build columns
	var columns []string
	sql.WriteString("(")
	columns = append(columns, b.insert.columns...)
	sql.WriteString(strings.Join(columns, ", "))

	//Builder values
	var values []string
	var row []string
	sql.WriteString(") VALUES ")

	//Build rows
	rowSize := len(b.insert.columns)
	for i, v := range b.insert.values {
		var val string
		switch v.typer {
		case "int":
			val = v.value
		case "string":
			val = fmt.Sprintf("'%s'", v.value)
		case "bool":
			if v.value == "true" {
				val = "1"
			} else {
				val = "0"
			}
		}

		row = append(row, val)

		if (i+1)%rowSize == 0 {
			values = append(values, fmt.Sprintf("(%s)", strings.Join(row, ", ")))
			row = []string{}
		}
	}
	sql.WriteString(strings.Join(values, ", "))
	return sql.String()
}

func (b *SQLBuilder) buildUpdate() string {
	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("UPDATE %s SET ", b.table))

	//Build statement
	var statements []string
	for i, v := range b.update.columns {
		def := fmt.Sprintf("%s = ", v)
		switch b.update.values[i].typer {
		case "int":
			def += b.update.values[i].value
		case "string":
			def += fmt.Sprintf("'%s'", b.update.values[i].value)
		case "bool":
			if b.update.values[i].value == "true" {
				def += "1"
			} else {
				def += "0"
			}
		}

		statements = append(statements, def)
	}

	sql.WriteString(strings.Join(statements, ", "))

	//Build where
	if b.where.column != "" {
		sql.WriteString(fmt.Sprintf(" WHERE %s %s ", b.where.column, b.where.operator))

		if b.where.column == "id" {
			sql.WriteString(b.where.value)
		} else {
			sql.WriteString(fmt.Sprintf("'%s'", b.where.value))
		}
	}

	return sql.String()
}

func (b *SQLBuilder) buildDelete() string {
	var sql strings.Builder
	sql.WriteString(fmt.Sprintf("DELETE FROM %s", b.table))

	//Build where
	if b.where.column != "" {
		sql.WriteString(fmt.Sprintf(" WHERE %s %s ", b.where.column, b.where.operator))

		if b.where.column == "id" {
			sql.WriteString(b.where.value)
		} else {
			sql.WriteString(fmt.Sprintf("'%s'", b.where.value))
		}
	}

	return sql.String()
}
