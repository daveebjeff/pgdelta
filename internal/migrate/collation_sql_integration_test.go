package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestCollationMigration_FullCycle(t *testing.T) {
	initial := schema.Schema{}

	target := schema.Schema{
		Collations: []schema.Collation{
			{
				Schema:        "public",
				Name:          "english_ci",
				Provider:      "icu",
				Locale:        "en-US-u-ks-level1",
				Deterministic: false,
			},
		},
	}

	diff := schema.DiffSchemas(initial, target)
	if len(diff.AddedCollations) != 1 {
		t.Fatalf("expected 1 added collation, got %d", len(diff.AddedCollations))
	}

	sqls := migrate.CollationDiffSQL(diff)
	if len(sqls) != 1 {
		t.Fatalf("expected 1 SQL statement, got %d", len(sqls))
	}

	expected := `CREATE COLLATION public.english_ci (PROVIDER = icu, LOCALE = 'en-US-u-ks-level1', DETERMINISTIC = false);`
	if sqls[0] != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sqls[0])
	}

	// Now diff target -> initial (removal)
	removeDiff := schema.DiffSchemas(target, initial)
	if len(removeDiff.RemovedCollations) != 1 {
		t.Fatalf("expected 1 removed collation, got %d", len(removeDiff.RemovedCollations))
	}

	dropSQLs := migrate.CollationDiffSQL(removeDiff)
	if len(dropSQLs) != 1 {
		t.Fatalf("expected 1 drop SQL statement, got %d", len(dropSQLs))
	}

	expectedDrop := `DROP COLLATION public.english_ci;`
	if dropSQLs[0] != expectedDrop {
		t.Errorf("expected:\n%s\ngot:\n%s", expectedDrop, dropSQLs[0])
	}
}
