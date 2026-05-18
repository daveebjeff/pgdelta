package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
)

func baseUserMapping() schema.UserMapping {
	return schema.UserMapping{
		User:   "alice",
		Server: "myserver",
		Options: map[string]string{
			"password": "secret",
			"user":     "remote_alice",
		},
	}
}

func TestCreateUserMappingSQL_WithOptions(t *testing.T) {
	m := baseUserMapping()
	got := migrate.CreateUserMappingSQL(m)
	want := "CREATE USER MAPPING FOR alice SERVER myserver OPTIONS (password 'secret', user 'remote_alice');"
	if got != want {
		t.Errorf("expected:\n%s\ngot:\n%s", want, got)
	}
}

func TestCreateUserMappingSQL_NoOptions(t *testing.T) {
	m := schema.UserMapping{User: "bob", Server: "otherserver", Options: nil}
	got := migrate.CreateUserMappingSQL(m)
	want := "CREATE USER MAPPING FOR bob SERVER otherserver;"
	if got != want {
		t.Errorf("expected:\n%s\ngot:\n%s", want, got)
	}
}

func TestDropUserMappingSQL(t *testing.T) {
	m := baseUserMapping()
	got := migrate.DropUserMappingSQL(m)
	want := "DROP USER MAPPING IF EXISTS FOR alice SERVER myserver;"
	if got != want {
		t.Errorf("expected:\n%s\ngot:\n%s", want, got)
	}
}

func TestAlterUserMappingSQL(t *testing.T) {
	m := baseUserMapping()
	got := migrate.AlterUserMappingSQL(m)
	want := "ALTER USER MAPPING FOR alice SERVER myserver OPTIONS (SET password 'secret', SET user 'remote_alice');"
	if got != want {
		t.Errorf("expected:\n%s\ngot:\n%s", want, got)
	}
}

func TestUserMappingDiffSQL_AddedAndRemoved(t *testing.T) {
	m := baseUserMapping()
	diff := schema.SchemaDiff{
		Added:   []interface{}{m},
		Removed: []interface{}{schema.UserMapping{User: "old", Server: "myserver"}},
	}
	stmts := migrate.UserMappingDiffSQL(diff)
	if len(stmts) != 2 {
		t.Errorf("expected 2 statements, got %d", len(stmts))
	}
}
