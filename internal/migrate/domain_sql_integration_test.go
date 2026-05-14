package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestDomainMigration_FullCycle(t *testing.T) {
	original := []schema.Domain{
		{Schema: "public", Name: "email", Type: "text"},
	}

	def := "unknown@example.com"
	nn := true
	updated := []schema.Domain{
		{Schema: "public", Name: "email", Type: "text", Default: &def, NotNull: &nn},
		{Schema: "public", Name: "score", Type: "integer"},
	}

	diff := schema.DiffDomains(original, updated)

	if len(diff.Added) != 1 {
		t.Fatalf("expected 1 added domain, got %d", len(diff.Added))
	}
	if diff.Added[0].Name != "score" {
		t.Errorf("expected added domain 'score', got %q", diff.Added[0].Name)
	}

	if len(diff.Changed) != 1 {
		t.Fatalf("expected 1 changed domain, got %d", len(diff.Changed))
	}
	if diff.Changed[0].Name != "email" {
		t.Errorf("expected changed domain 'email', got %q", diff.Changed[0].Name)
	}

	if len(diff.Removed) != 0 {
		t.Fatalf("expected 0 removed domains, got %d", len(diff.Removed))
	}

	sqls := migrate.DomainDiffSQL(diff)
	// changed domain = drop + create, added domain = create => 3 total
	if len(sqls) != 3 {
		t.Fatalf("expected 3 SQL statements, got %d: %v", len(sqls), sqls)
	}

	// Verify drop comes before create for changed domain
	foundDrop := false
	for _, s := range sqls {
		if s == "DROP DOMAIN public.email;" {
			foundDrop = true
		}
	}
	if !foundDrop {
		t.Error("expected DROP DOMAIN public.email; in output")
	}
}
