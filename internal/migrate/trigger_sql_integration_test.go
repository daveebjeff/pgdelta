package migrate_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestTriggerMigration_FullCycle(t *testing.T) {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.trigger_test_orders (id serial PRIMARY KEY, status text)`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	defer db.Exec(`DROP TABLE IF EXISTS public.trigger_test_orders`)

	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION public.noop_trigger_fn()
		RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END;$$
	`)
	if err != nil {
		t.Fatalf("failed to create trigger function: %v", err)
	}
	defer db.Exec(`DROP FUNCTION IF EXISTS public.noop_trigger_fn()`)

	trg := schema.Trigger{
		Schema:     "public",
		Name:       "trg_integration_test",
		Table:      "trigger_test_orders",
		Definition: "CREATE TRIGGER trg_integration_test BEFORE INSERT ON public.trigger_test_orders FOR EACH ROW EXECUTE FUNCTION public.noop_trigger_fn()",
	}

	createSQL := migrate.CreateTriggerSQL(trg)
	_, err = db.Exec(createSQL)
	if err != nil {
		t.Fatalf("CreateTriggerSQL failed: %v\nSQL: %s", err, createSQL)
	}

	var count int
	err = db.QueryRow(
		`SELECT COUNT(*) FROM information_schema.triggers WHERE trigger_name = $1`,
		"trg_integration_test",
	).Scan(&count)
	if err != nil || count != 1 {
		t.Fatalf("trigger not found after create: count=%d err=%v", count, err)
	}

	dropSQL := migrate.DropTriggerSQL(trg)
	_, err = db.Exec(dropSQL)
	if err != nil {
		t.Fatalf("DropTriggerSQL failed: %v\nSQL: %s", err, dropSQL)
	}

	err = db.QueryRow(
		`SELECT COUNT(*) FROM information_schema.triggers WHERE trigger_name = $1`,
		"trg_integration_test",
	).Scan(&count)
	if err != nil {
		t.Fatalf("query after drop failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected trigger to be dropped, but count=%d", count)
	}

	fmt.Println("TestTriggerMigration_FullCycle passed")
}
