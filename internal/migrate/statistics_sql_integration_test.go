package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestStatisticsMigration_FullCycle(t *testing.T) {
	before := schema.Schema{
		Statistics: []schema.Statistic{
			{
				Schema:  "public",
				Name:    "order_stats",
				Table:   "orders",
				Columns: []string{"user_id", "status"},
				Kinds:   []string{"dependencies", "ndistinct"},
			},
		},
	}

	after := schema.Schema{
		Statistics: []schema.Statistic{
			{
				Schema:  "public",
				Name:    "order_stats",
				Table:   "orders",
				Columns: []string{"user_id", "status"},
				Kinds:   []string{"dependencies"},
			},
			{
				Schema:  "public",
				Name:    "item_stats",
				Table:   "order_items",
				Columns: []string{"product_id", "quantity"},
				Kinds:   []string{"mcv"},
			},
		},
	}

	diff := schema.DiffSchemas(before, after)
	sqls := migrate.SchemaDiffSQL(diff)

	if len(sqls) == 0 {
		t.Fatal("expected SQL statements for statistics changes, got none")
	}

	foundDrop := false
	foundCreate := false
	for _, sql := range sqls {
		if sql == `DROP STATISTICS public.order_stats;` {
			foundDrop = true
		}
		if sql == `CREATE STATISTICS public.order_stats (dependencies) ON user_id, status FROM public.orders;` {
			foundCreate = true
		}
	}

	if !foundDrop {
		t.Error("expected DROP STATISTICS for changed statistic")
	}
	if !foundCreate {
		t.Error("expected CREATE STATISTICS for changed statistic")
	}
}
