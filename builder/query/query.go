package query

import "github.com/Milkado/go-migrations/config"

type (
	Query interface {
		Table(table string) Table
	}

	Table interface {
		Insert([]Row) Table
		Update([]Values) Table
        Delete () Table
		Where(column string, value string, operator ...string) Table
		Build() string
	}

	Row []Values

	Values struct {
		Column string
		Value string
		Type string
	}

	BuilderSQL struct {
		SQLBuilder *SQLBuilder
	}

	TableBuilder struct {
		SQLBuilder *SQLBuilder
	}
)

func DB() *BuilderSQL {
	return &BuilderSQL{
		SQLBuilder: &SQLBuilder{
			dialect: Dialect(config.Env("DB_DRIVER")),
		},
	}
}

func (s *BuilderSQL) Table(table string) Table {
	s.SQLBuilder.table = table
	return &TableBuilder{SQLBuilder: s.SQLBuilder}
}

func (t *TableBuilder) Insert(values []Row) Table {
	t.SQLBuilder.generator = Insert
	for _, v := range values {
		for _, val := range v {
			typer := "string"
			if val.Type != "" {
				typer = val.Type
			}
			if len(t.SQLBuilder.insert.columns) != len(v) {
				t.SQLBuilder.insert.columns = append(t.SQLBuilder.insert.columns, val.Column)
			} 
			t.SQLBuilder.insert.values = append(t.SQLBuilder.insert.values, ValuesGen{
				value: val.Value,
				typer: typer,
			})
		}
	}

	return t
}

func (t *TableBuilder) Update(values []Values) Table {
	t.SQLBuilder.generator = Update
	for _, v := range values {
		typer := "string"
		if v.Type != "" {
			typer = v.Type
		}
		t.SQLBuilder.update.columns = append(t.SQLBuilder.update.columns, v.Column)
		t.SQLBuilder.update.values = append(t.SQLBuilder.update.values, ValuesGen{
			value: v.Value,
			typer: typer,
		})
	}
	return t
}

func (t *TableBuilder) Delete() Table {
    t.SQLBuilder.generator = Delete
    return t
}

func (t *TableBuilder) Where(column string, value string, operator ...string) Table {
	oper := "="
	if len(operator) > 0 {
		oper = operator[0]
	}
	t.SQLBuilder.where.column = column
	t.SQLBuilder.where.value = value
	t.SQLBuilder.where.operator = oper
	return t
}

func (t *TableBuilder) Build() string {
	return t.SQLBuilder.Build()
}
