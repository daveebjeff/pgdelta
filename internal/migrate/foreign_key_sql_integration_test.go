package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestForeignKeyMigration_FullCycle(t *testing.T) {
	old := schema.Schema{
		ForeignKeys: []schema.ForeignKey{},
	}
	new := schema.Schema{
		ForeignKeys: []schema.ForeignKey{
			{
				Schema:     "public",
				Table:      "orders",
				Name:       "fk_orders_customer",
				Columns:    []string{"customer_id"},
				RefSchema:  "public",
				RefTable:   "customers",
				RefColumns: []string{"id"},
				OnDelete:   "CASCADE",
				OnUpdate:   "NO ACTION",
			},
		},
	}

	diff := schema.DiffSchemas(old, new)
	sqls := migrate.SchemaDiffSQL(diff)

	if len(sqls) == 0 {
		t.Fatal("expected SQL statements for added foreign key, got none")
	}

	found := false
	for _, sql := range sqls {
		if sql == `ALTER TABLE "public"."orders" ADD CONSTRAINT "fk_orders_customer" FOREIGN KEY ("customer_id") REFERENCES "public"."customers" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;` {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected ADD CONSTRAINT statement in output, got: %v", sqls)
	}

	// reverse: drop
	diff2 := schema.DiffSchemas(new, old)
	sqls2 := migrate.SchemaDiffSQL(diff2)

	if len(sqls2) == 0 {
		t.Fatal("expected SQL statements for removed foreign key, got none")
	}

	found2 := false
	for _, sql := range sqls2 {
		if sql == `ALTER TABLE "public"."orders" DROP CONSTRAINT "fk_orders_customer";` {
			found2 = true
			break
		}
	}
	if !found2 {
		t.Errorf("expected DROP CONSTRAINT statement in output, got: %v", sqls2)
	}
}
