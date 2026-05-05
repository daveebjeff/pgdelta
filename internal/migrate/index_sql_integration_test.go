package migrate_test

import (
	"strings"
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
)

// TestIndexMigration_FullCycle tests a realistic scenario where indexes are
// diffed and migration SQL is generated end-to-end.
func TestIndexMigration_FullCycle(t *testing.T) {
	oldIndexes := []schema.Index{
		{
			SchemaName: "public",
			TableName:  "orders",
			Name:       "idx_orders_user_id",
			Columns:    []string{"user_id"},
			Unique:     false,
			Method:     schema.IndexMethodBTree,
		},
	}

	newIndexes := []schema.Index{
		{
			SchemaName: "public",
			TableName:  "orders",
			Name:       "idx_orders_user_id",
			Columns:    []string{"user_id", "created_at"},
			Unique:     false,
			Method:     schema.IndexMethodBTree,
		},
		{
			SchemaName: "public",
			TableName:  "orders",
			Name:       "idx_orders_status",
			Columns:    []string{"status"},
			Unique:     false,
			Method:     schema.IndexMethodHash,
		},
	}

	diff := schema.DiffIndexes(oldIndexes, newIndexes)

	if len(diff.Changed) != 1 {
		t.Fatalf("expected 1 changed index, got %d", len(diff.Changed))
	}
	if len(diff.Added) != 1 {
		t.Fatalf("expected 1 added index, got %d", len(diff.Added))
	}

	stmts := migrate.IndexDiffSQL(diff)

	// Changed index: DROP + CREATE = 2, Added: CREATE = 1 => total 3
	if len(stmts) != 3 {
		t.Fatalf("expected 3 SQL statements, got %d: %v", len(stmts), stmts)
	}

	hasHash := false
	for _, s := range stmts {
		if strings.Contains(s, "USING hash") {
			hasHash = true
		}
	}
	if !hasHash {
		t.Errorf("expected a USING hash statement among: %v", stmts)
	}
}
