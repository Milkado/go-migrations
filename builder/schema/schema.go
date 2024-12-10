package schema

import "fmt"

type (
	Schema interface {
		Create(table string) Table
		Drop(table string) Table
		Alter(table string) Alteration
	}

	Table interface {
		Id() Table
		String(name string, nullable bool, length ...int) Table
		Text(name string, nullable bool) Table
		Integer(name string, nullable bool) Table
		BigInteger(name string, nullable bool) Table
		Float(name string, nullable bool, precision ...int) Table
		Timestamps() Table
		SoftDeletes() Table
		Unique(columns ...string) Table
		Index(columns ...string) Table
		ForeignId(column string, table string, nullable bool) Reference
		Build() string
	}

	Reference interface {
		Refenreces(column string) Reference
		On(table string) Reference
		OnDelete(action string) Reference
		OnUpdate(action string) Reference
		Build() string
	}

	Alteration interface {
		String(name string, nullable bool, length ...int) Alteration
		Text(name string, nullable bool) Alteration
		Integer(name string, nullable bool) Alteration
		BigInteger(name string, nullable bool) Alteration
		Float(name string, nullable bool, precision ...int) Alteration
		Timestamps() Alteration
		SoftDeletes() Alteration
		Unique(columns ...string) Alteration
		Index(columns ...string) Alteration
		ForeignId(column string, table string, nullable bool) Reference
		DropColumn(column ...string) Alteration
		ModifyColumn(column string, columnType string, nullable bool) Alteration
		RenameColumn(column string, newName string) Alteration
		Build() string
	}

	SchemaBuilder struct {
		sqlBuilder *SQLBuilder
	}

	TableBuilder struct {
		sqlBuilder *SQLBuilder
	}

	ReferenceBuilder struct {
		sqlBuilder *SQLBuilder
		column string
		table string
	}

 	AlterationBuilder struct {
		sqlBuilder *SQLBuilder
	}
)

func Query() *SchemaBuilder {
	return &SchemaBuilder{
		sqlBuilder: &SQLBuilder{},
	}
}

func (s *SchemaBuilder) Create(tableName string) Table {
    s.sqlBuilder.sqlType = CreateTable
	s.sqlBuilder.table = tableName
	return &TableBuilder{sqlBuilder: s.sqlBuilder}
}

func (s *SchemaBuilder) Drop(tableName string) Table {
    s.sqlBuilder.sqlType = DropTable
	s.sqlBuilder.table = tableName
	return &TableBuilder{sqlBuilder: s.sqlBuilder}
}

func (s *SchemaBuilder) Alter(tableName string) Alteration {
    s.sqlBuilder.sqlType = AlterTable
	s.sqlBuilder.table = tableName
	return &AlterationBuilder{sqlBuilder: s.sqlBuilder}
}

func (t *TableBuilder) Id() Table {
	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:       "id",
		Type:       "BIGINT",
		Nullable:   false,
		PrimaryKey: true,
		AutoInc:    true,
	})

	return t
}

func (t *TableBuilder) String(name string, nullable bool, length ...int) Table {
	l := 255
	if len(length) > 0 {
		l = length[0]
	}

	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:     name,
		Type:     "VARCHAR",
		Length:   l,
		Nullable: nullable,
	})
	return t
}

func (t *TableBuilder) Text(name string, nullable bool) Table {
	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:     name,
		Type:     "TEXT",
		Nullable: nullable,
	})

	return t
}

func (t *TableBuilder) Integer(name string, nullable bool) Table {
	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:     name,
		Type:     "INT",
		Nullable: nullable,
	})

	return t
}

func (t *TableBuilder) BigInteger(name string, nullable bool) Table {
	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:     name,
		Type:     "BIGINT",
		Nullable: nullable,
	})

	return t
}

func (t *TableBuilder) Float(name string, nullable bool, precision ...int) Table {
	p := 8
	if len(precision) > 0 {
		p = precision[0]
	}

	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:     name,
		Type:     fmt.Sprintf("FLOAT(%d)", p),
		Nullable: nullable,
	})

	return t
}

func (t *TableBuilder) Timestamps() Table {
	defaultTimestamp := "CURRENT_TIMESTAMP"
	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:     "created_at",
		Type:     "TIMESTAMP",
		Default: &defaultTimestamp,
		Nullable: false,
	})

	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:     "updated_at",
		Type:     "TIMESTAMP",
		Default: &defaultTimestamp,
		Nullable: false,
	})

	return t
}

func (t *TableBuilder) SoftDeletes() Table {
	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:     "deleted_at",
		Type:     "TIMESTAMP",
		Nullable: true,
	})

	return t
}

func (t *TableBuilder) Unique(columns ...string) Table {
	for _, col := range columns {
		name := col + "_unique"
		t.sqlBuilder.indexes = append(t.sqlBuilder.indexes, Index{
			Name:    name,
			Columns: []string{col},
			Type:    "unique",
		})
	}

	return t
}

func (t *TableBuilder) Index(columns ...string) Table {
	for _, col := range columns {
		name := col + "_index"
		t.sqlBuilder.indexes = append(t.sqlBuilder.indexes, Index{
			Name:    name,
			Columns: []string{col},
			Type:    "index",
		})
	}

	return t
}

func (t *TableBuilder) ForeignId(column string, table string, nullable bool) Reference {
	t.sqlBuilder.columns = append(t.sqlBuilder.columns, Column{
		Name:       column,
		Type: "BIGINT",
		Nullable:   nullable,
	})

	return &ReferenceBuilder {
		sqlBuilder: t.sqlBuilder,
		column: column,
		table: table,
	}
}

func (r *ReferenceBuilder) Refenreces(column string) Reference {
	r.sqlBuilder.references = append(r.sqlBuilder.references, References{
		Column:     r.column,
		RefTable:   r.table,
		RefColumn:  column,
	})

	return r
}

func (r *ReferenceBuilder) On(table string) Reference {
	r.table = table
	return r
}

func (r *ReferenceBuilder) OnDelete(action string) Reference {
	if len(r.sqlBuilder.references) > 0 {
		lastIdx := len(r.sqlBuilder.references) - 1
		ref := r.sqlBuilder.references[lastIdx]
		ref.OnDelete = action
		r.sqlBuilder.references[lastIdx] = ref
	}
	return r
}

func (r *ReferenceBuilder) OnUpdate(action string) Reference {
	if len(r.sqlBuilder.references) > 0 {
		lastIdx := len(r.sqlBuilder.references) - 1
		ref := r.sqlBuilder.references[lastIdx]
		ref.OnUpdate = action
		r.sqlBuilder.references[lastIdx] = ref
	}
	return r
}

func (a *AlterationBuilder) String(name string, nullable bool, length ...int) Alteration {
	l := 255
	if len(length) > 0 {
		l = length[0]
	}

	a.sqlBuilder.columns = append(a.sqlBuilder.columns, Column{
		Name:     name,
		Type:     "VARCHAR",
		Length:   l,
		Nullable: nullable,
	})

	return a
}

func (a *AlterationBuilder) Text(name string, nullable bool) Alteration {
	a.sqlBuilder.columns = append(a.sqlBuilder.columns, Column{
		Name:     name,
		Type:     "TEXT",
		Nullable: nullable,
	})

	return a
}

func (a *AlterationBuilder) Integer(name string, nullable bool) Alteration {
	a.sqlBuilder.columns = append(a.sqlBuilder.columns, Column{
		Name:     name,
		Type:     "INT",
		Nullable: nullable,
	})

	return a
}

func (a *AlterationBuilder) BigInteger(name string, nullable bool) Alteration {
	a.sqlBuilder.columns = append(a.sqlBuilder.columns, Column{
		Name:     name,
		Type:     "BIGINT",
		Nullable: nullable,
	})

	return a
}

func (a *AlterationBuilder) Float(name string, nullable bool, precision ...int) Alteration {
	p := 8
	if len(precision) > 0 {
		p = precision[0]
	}

	a.sqlBuilder.columns = append(a.sqlBuilder.columns, Column{
		Name:     name,
		Type:     fmt.Sprintf("FLOAT(%d)", p),
		Nullable: nullable,
	})

	return a
}

func (a *AlterationBuilder) Timestamps() Alteration {
	defaultTimestamp := "CURRENT_TIMESTAMP"
	a.sqlBuilder.columns = append(a.sqlBuilder.columns, Column{
		Name:     "created_at",
		Type:     "TIMESTAMP",
		Default: &defaultTimestamp,
		Nullable: false,
	})

	a.sqlBuilder.columns = append(a.sqlBuilder.columns, Column{
		Name:     "updated_at",
		Type:     "TIMESTAMP",
		Default: &defaultTimestamp,
		Nullable: false,
	})

	return a
}

func (a *AlterationBuilder) SoftDeletes() Alteration {
	a.sqlBuilder.columns = append(a.sqlBuilder.columns, Column{
		Name:     "deleted_at",
		Type:     "TIMESTAMP",
		Nullable: true,
	})

	return a
}

func (a *AlterationBuilder) Unique(columns ...string) Alteration {
	for _, col := range columns {
		name := col + "_unique"
		a.sqlBuilder.indexes = append(a.sqlBuilder.indexes, Index{
			Name:    name,
			Columns: []string{col},
			Type:    "unique",
		})
	}

	return a
}

func (a *AlterationBuilder) Index(columns ...string) Alteration {
	for _, col := range columns {
		name := col + "_index"
		a.sqlBuilder.indexes = append(a.sqlBuilder.indexes, Index{
			Name:    name,
			Columns: []string{col},
			Type:    "index",
		})
	}

	return a
}

func (a *AlterationBuilder) ForeignId(column string, table string, nullable bool) Reference {
	a.sqlBuilder.columns = append(a.sqlBuilder.columns, Column{
		Name:       column,
		Type: "BIGINT",
		Nullable:   nullable,
	})

	return &ReferenceBuilder {
		sqlBuilder: a.sqlBuilder,
		column: column,
		table: table,
	}
}

func (a *AlterationBuilder) DropColumn(columns ...string) Alteration {
	a.sqlBuilder.dropColumns = append(a.sqlBuilder.dropColumns, columns...)

	return a
}

func (a *AlterationBuilder) ModifyColumn(column string, columnType string, nullable bool) Alteration {
	a.sqlBuilder.modifyColumns = append(a.sqlBuilder.modifyColumns, Column{
		Name: column,
		Type: columnType,
		Nullable: nullable,
	})

	return a
}

func (a *AlterationBuilder) RenameColumn(column string, newName string) Alteration {
	a.sqlBuilder.renameColumns = append(a.sqlBuilder.renameColumns, RenameColumn{
		OldName: column,
		NewName: newName,
	})

	return a
}

func (r *ReferenceBuilder) Build() string {
	return r.sqlBuilder.Build()
}

func (t *TableBuilder) Build() string {
	return t.sqlBuilder.Build()
}

func (a *AlterationBuilder) Build() string {
	return a.sqlBuilder.Build()
}