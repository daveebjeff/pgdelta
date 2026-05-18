package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
)

func TestUserMappingMigration_FullCycle(t *testing.T) {
	original := schema.UserMapping{
		User:   "alice",
		Server: "myserver",
		Options: map[string]string{
			"user":     "remote_alice",
			"password": "secret",
		},
	}

	updated := schema.UserMapping{
		User:   "alice",
		Server: "myserver",
		Options: map[string]string{
			"user":     "remote_alice",
			"password": "newsecret",
		},
	}

	diff := schema.DiffUserMappings(
		[]schema.UserMapping{original},
		[]schema.UserMapping{updated},
	)

	if len(diff.Changed) != 1 {
		t.Fatalf("expected 1 changed mapping, got %d", len(diff.Changed))
	}

	stmts := migrate.UserMappingDiffSQL(diff)
	if len(stmts) != 1 {
		t.Fatalf("expected 1 SQL statement, got %d", len(stmts))
	}

	expected := "ALTER USER MAPPING FOR alice SERVER myserver OPTIONS (SET password 'newsecret', SET user 'remote_alice');"
	if stmts[0] != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, stmts[0])
	}

	// Test drop cycle
	dropDiff := schema.DiffUserMappings([]schema.UserMapping{original}, nil)
	dropStmts := migrate.UserMappingDiffSQL(dropDiff)
	if len(dropStmts) != 1 {
		t.Fatalf("expected 1 drop statement, got %d", len(dropStmts))
	}
	if dropStmts[0] != "DROP USER MAPPING IF EXISTS FOR alice SERVER myserver;" {
		t.Errorf("unexpected drop statement: %s", dropStmts[0])
	}
}
