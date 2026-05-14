package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestCreateDomainSQL_Basic(t *testing.T) {
	d := schema.Domain{
		Schema: "public",
		Name:   "email",
		Type:   "text",
	}
	got := migrate.CreateDomainSQL(d)
	want := "CREATE DOMAIN public.email AS text;"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCreateDomainSQL_WithNotNull(t *testing.T) {
	nn := true
	d := schema.Domain{
		Schema:  "public",
		Name:    "positive_int",
		Type:    "integer",
		NotNull: &nn,
	}
	got := migrate.CreateDomainSQL(d)
	want := "CREATE DOMAIN public.positive_int AS integer NOT NULL;"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCreateDomainSQL_WithDefault(t *testing.T) {
	def := "0"
	d := schema.Domain{
		Schema:  "public",
		Name:    "score",
		Type:    "integer",
		Default: &def,
	}
	got := migrate.CreateDomainSQL(d)
	want := "CREATE DOMAIN public.score AS integer DEFAULT 0;"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDropDomainSQL(t *testing.T) {
	d := schema.Domain{
		Schema: "public",
		Name:   "email",
		Type:   "text",
	}
	got := migrate.DropDomainSQL(d)
	want := "DROP DOMAIN public.email;"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDomainDiffSQL_AddedAndRemoved(t *testing.T) {
	old := []schema.Domain{
		{Schema: "public", Name: "old_domain", Type: "text"},
	}
	new := []schema.Domain{
		{Schema: "public", Name: "new_domain", Type: "integer"},
	}
	diff := schema.DiffDomains(old, new)
	sqls := migrate.DomainDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(sqls))
	}
}

func TestDomainDiffSQL_Changed(t *testing.T) {
	old := []schema.Domain{
		{Schema: "public", Name: "score", Type: "integer"},
	}
	def := "0"
	new := []schema.Domain{
		{Schema: "public", Name: "score", Type: "integer", Default: &def},
	}
	diff := schema.DiffDomains(old, new)
	sqls := migrate.DomainDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements (drop+create), got %d", len(sqls))
	}
}
