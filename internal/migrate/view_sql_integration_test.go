package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
)

func TestViewMigration_FullCycle(t *testing.T) {
	base := []schema.View{
		{Schema: "public", Name: "user_summary", Definition: "SELECT id, name FROM users"},
		{Schema: "public", Name: "order_summary", Definition: "SELECT id, total FROM orders"},
	}
	target := []schema.View{
		{Schema: "public", Name: "user_summary", Definition: "SELECT id, name, email FROM users"},
		{Schema: "public", Name: "product_summary", Definition: "SELECT id, title FROM products"},
	}

	diff := schema.DiffViews(base, target)
	sqls := migrate.ViewDiffSQL(diff)

	if len(sqls) == 0 {
		t.Fatal("expected SQL statements, got none")
	}

	hasDropOrderSummary := false
	hasCreateProductSummary := false
	hasReplaceUserSummary := false

	for _, sql := range sqls {
		switch sql {
		case "DROP VIEW IF EXISTS public.order_summary;":
			hasDropOrderSummary = true
		case "CREATE OR REPLACE VIEW public.product_summary AS SELECT id, title FROM products;":
			hasCreateProductSummary = true
		case "CREATE OR REPLACE VIEW public.user_summary AS SELECT id, name, email FROM users;":
			hasReplaceUserSummary = true
		}
	}

	if !hasDropOrderSummary {
		t.Error("expected DROP VIEW for order_summary")
	}
	if !hasCreateProductSummary {
		t.Error("expected CREATE OR REPLACE VIEW for product_summary")
	}
	if !hasReplaceUserSummary {
		t.Error("expected CREATE OR REPLACE VIEW for updated user_summary")
	}
}
