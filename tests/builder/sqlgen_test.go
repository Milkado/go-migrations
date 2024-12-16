package builder_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Milkado/go-migrations/builder/schema"
	"github.com/Milkado/go-migrations/cmd"
	"github.com/spf13/viper"
)

func TestSQLGenerationMysql(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting current directory: %v", err)
	}
	defer os.Chdir(originalDir) // Restore original directory when done

	// Change to project root
	projectRoot := filepath.Dir(filepath.Dir(originalDir))
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("error changing to project root: %v", err)
	}
	//Expected queries with line breaks to make it easier to read
	tests := []struct {
		name     string
		builder  func() string
		expected string
	}{
		{
			name: "Simple table with ID",
			builder: func() string {
				return schema.Query().Create("users").Id().Build()
			},
			expected: `CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY)`,
		},
		{
			name: "Table with multiple columns",
			builder: func() string {
				return schema.Query().Create("users").
					Id().
					String("name", false).
					Integer("age", false).
					Timestamps().
					Build()
			},
			expected: `CREATE TABLE IF NOT EXISTS users (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				age INT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`,
		},
		{
			name: "Table with foreign key",
			builder: func() string {
				return schema.Query().
					Create("posts").
					Id().
					String("title", false).
					ForeignId("user_id", false).
					Refenreces("id").
					OnDelete("CASCADE").
					On("users").
					Build()
			},
			expected: `CREATE TABLE IF NOT EXISTS posts (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				user_id BIGINT NOT NULL,
				FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE)`,
		},
		{
			name: "Alter table add constraint",
			builder: func() string {
				return schema.Query().
					Alter("posts").
					Unique("title").
					Build()
			},
			expected: `ALTER TABLE posts ADD CONSTRAINT title_unique UNIQUE (title)`,
		},
		{
			name: "Drop table",
			builder: func() string {
				return schema.Query().Drop("users").Build()
			},
			expected: `DROP TABLE IF EXISTS users`,
		},
		{
			name: "Alter table rename column",
			builder: func() string {
				return schema.Query().Alter("users").RenameColumn("name", "first_name").Build()
			},
			expected: `ALTER TABLE users RENAME COLUMN name TO first_name`,
		},
		{
			name: "Alter table modify column",
			builder: func() string {
				return schema.Query().Alter("users").ModifyColumn("active", "BOOLEAN", false).Build()
			},
			expected: `ALTER TABLE users MODIFY COLUMN active BOOLEAN NOT NULL`,
		},
		{
			name: "Alter table add two foreign keys",
			builder: func() string {
				return schema.Query().
					Alter("posts").
					ForeignId("user_id", false).
					Refenreces("id").
					OnDelete("CASCADE").
					On("users").
					ForeignId("category_id", false).
					Refenreces("id").
					OnDelete("CASCADE").
					On("categories").
					Build()
			},
			expected: `ALTER TABLE posts ADD COLUMN user_id BIGINT NOT NULL,
				ADD COLUMN category_id BIGINT NOT NULL,
				ADD CONSTRAINT user_id_users_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
				ADD CONSTRAINT category_id_categories_fk FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.builder()
			// Normalize spaces and new lines
			got = normalizeSpaces(got)
			expected := normalizeSpaces(tt.expected)

			if got != expected {
				t.Errorf(cmd.Red+"\nwant: %s\ngot: %s\n"+cmd.Reset, expected, got)
				return
			}
		})
	}
}

func TestSQLGenerationPostgres(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting current directory: %v", err)
	}
	defer os.Chdir(originalDir) // Restore original directory when done

	// Change to project root
	projectRoot := filepath.Dir(filepath.Dir(originalDir))
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("error changing to project root: %v", err)
	}

	//Change env to postgres
	changeEnvDriver("postgres")

	tests := []struct {
		name     string
		builder  func() string
		expected string
	}{
		{
			name: "Simple table with ID",
			builder: func() string {
				return schema.Query().Create("users").Id().Build()
			},
			expected: `CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY)`,
		},
		{
			name: "Table with multiple columns",
			builder: func() string {
				return schema.Query().Create("users").
					Id().
					String("name", false).
					Integer("age", false).
					Timestamps().
					Build()
			},
			expected: `CREATE TABLE IF NOT EXISTS users (
				id BIGSERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				age INT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`,
		},
		{
			name: "Table with foreign key",
			builder: func() string {
				return schema.Query().
					Create("posts").
					Id().
					String("title", false).
					ForeignId("user_id", false).
					Refenreces("id").
					OnDelete("CASCADE").
					On("users").
					Build()
			},
			expected: `CREATE TABLE IF NOT EXISTS posts (
				id BIGSERIAL PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				user_id BIGINT NOT NULL,
				FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE)`,
		},
		{
			name: "Alter table add constraint",
			builder: func() string {
				return schema.Query().
					Alter("posts").
					Unique("title").
					Build()
			},
			expected: `ALTER TABLE posts ADD CONSTRAINT title_unique UNIQUE (title)`,
		},
		{
			name: "Drop table",
			builder: func() string {
				return schema.Query().Drop("users").Build()
			},
			expected: `DROP TABLE IF EXISTS users`,
		},
		{
			name: "Alter table rename column",
			builder: func() string {
				return schema.Query().Alter("users").RenameColumn("name", "first_name").Build()
			},
			expected: `ALTER TABLE users RENAME COLUMN name TO first_name`,
		},
		{
			name: "Alter table modify column",
			builder: func() string {
				return schema.Query().Alter("users").ModifyColumn("active", "BOOLEAN", false).Build()
			},
			expected: `ALTER TABLE users MODIFY COLUMN active BOOLEAN NOT NULL`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.builder()
			// Normalize spaces and new lines
			got = normalizeSpaces(got)
			expected := normalizeSpaces(tt.expected)

			if got != expected {
				t.Errorf(cmd.Red+"\nwant: %s\ngot: %s\n"+cmd.Reset, expected, got)
				return
			}
		})
	}
}

func normalizeSpaces(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ReplaceAll(s, "( ", "(")
	s = strings.ReplaceAll(s, " )", ")")

	return s
}

func changeEnvDriver(driver string) {
	viper.Set("DB_DRIVER", driver)
}
