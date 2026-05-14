package migrate_test

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestOperatorMigration_FullCycle(t *testing.T) {
	dsn := os.Getenv("PGDELTA_TEST_DSN")
	if dsn == "" {
		t.Skip("PGDELTA_TEST_DSN not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	op := schema.Operator{
		Schema:     "public",
		Name:       "##",
		LeftType:   "integer",
		RightType:  "integer",
		ResultType: "integer",
		Procedure:  "int4pl",
	}

	createSQL := migrate.CreateOperatorSQL(op)
	if _, err := db.Exec(createSQL); err != nil {
		t.Fatalf("failed to create operator: %v", err)
	}

	var count int
	query := fmt.Sprintf(
		`SELECT COUNT(*) FROM pg_operator o JOIN pg_namespace n ON o.oprnamespace = n.oid
		 WHERE n.nspname = '%s' AND o.oprname = '%s'`,
		op.Schema, strings.TrimSpace(op.Name),
	)
	if err := db.QueryRow(query).Scan(&count); err != nil {
		t.Fatalf("failed to query operator: %v", err)
	}
	if count != 1 {
		t.Errorf("expected operator to exist, count=%d", count)
	}

	dropSQL := migrate.DropOperatorSQL(op)
	if _, err := db.Exec(dropSQL); err != nil {
		t.Fatalf("failed to drop operator: %v", err)
	}
}
