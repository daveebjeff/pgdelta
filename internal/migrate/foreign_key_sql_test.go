package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

var baseForeignKey = schema.ForeignKey{
	Schema:           "public",
	Table:            "orders",
	Name:             "fk_orders_customer",
	Columns:          []string{"customer_id"},
	RefSchema:        "public",
	RefTable:         "customers",
	RefColumns:       []string{"id"},
	OnDelete:         "CASCADE",
	OnUpdate:         "NO ACTION",
	Deferrable:       false,
	InitiallyDeferred: false,
}

func TestAddForeignKeySQL_Basic(t *testing.T) {
	sql := migrate.AddForeignKeySQL(baseForeignKey)
	expected := `ALTER TABLE "public"."orders" ADD CONSTRAINT "fk_orders_customer" FOREIGN KEY ("customer_id") REFERENCES "public"."customers" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestAddForeignKeySQL_MultiColumn(t *testing.T) {
	fk := baseForeignKey
	fk.Columns = []string{"a", "b"}
	fk.RefColumns = []string{"x", "y"}
	sql := migrate.AddForeignKeySQL(fk)
	expected := `ALTER TABLE "public"."orders" ADD CONSTRAINT "fk_orders_customer" FOREIGN KEY ("a", "b") REFERENCES "public"."customers" ("x", "y") ON DELETE CASCADE ON UPDATE NO ACTION;`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestAddForeignKeySQL_Deferrable(t *testing.T) {
	fk := baseForeignKey
	fk.Deferrable = true
	fk.InitiallyDeferred = true
	sql := migrate.AddForeignKeySQL(fk)
	expected := `ALTER TABLE "public"."orders" ADD CONSTRAINT "fk_orders_customer" FOREIGN KEY ("customer_id") REFERENCES "public"."customers" ("id") ON DELETE CASCADE ON UPDATE NO ACTION DEFERRABLE INITIALLY DEFERRED;`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestDropForeignKeySQL(t *testing.T) {
	sql := migrate.DropForeignKeySQL(baseForeignKey)
	expected := `ALTER TABLE "public"."orders" DROP CONSTRAINT "fk_orders_customer";`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestForeignKeyDiffSQL_AddedAndRemoved(t *testing.T) {
	diff := schema.ForeignKeyDiff{
		Added:   []schema.ForeignKey{baseForeignKey},
		Removed: []schema.ForeignKey{baseForeignKey},
	}
	sqls := migrate.ForeignKeyDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(sqls))
	}
}

func TestForeignKeyDiffSQL_NoChanges(t *testing.T) {
	diff := schema.ForeignKeyDiff{}
	sqls := migrate.ForeignKeyDiffSQL(diff)
	if len(sqls) != 0 {
		t.Fatalf("expected 0 statements, got %d", len(sqls))
	}
}
