package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestCreateTablespaceSQL_WithOwner(t *testing.T) {
	ts := schema.Tablespace{
		Name:     "fast_storage",
		Owner:    "admin",
		Location: "/mnt/ssd/pg_data",
	}
	got := migrate.CreateTablespaceSQL(ts)
	want := "CREATE TABLESPACE fast_storage OWNER admin LOCATION '/mnt/ssd/pg_data';"
	if got != want {
		t.Errorf("CreateTablespaceSQL() = %q, want %q", got, want)
	}
}

func TestCreateTablespaceSQL_NoOwner(t *testing.T) {
	ts := schema.Tablespace{
		Name:     "archive",
		Location: "/mnt/hdd/pg_archive",
	}
	got := migrate.CreateTablespaceSQL(ts)
	want := "CREATE TABLESPACE archive LOCATION '/mnt/hdd/pg_archive';"
	if got != want {
		t.Errorf("CreateTablespaceSQL() = %q, want %q", got, want)
	}
}

func TestDropTablespaceSQL(t *testing.T) {
	ts := schema.Tablespace{Name: "old_storage"}
	got := migrate.DropTablespaceSQL(ts)
	want := "DROP TABLESPACE old_storage;"
	if got != want {
		t.Errorf("DropTablespaceSQL() = %q, want %q", got, want)
	}
}

func TestAlterTablespaceOwnerSQL(t *testing.T) {
	ts := schema.Tablespace{Name: "fast_storage", Owner: "newowner"}
	got := migrate.AlterTablespaceOwnerSQL(ts)
	want := "ALTER TABLESPACE fast_storage OWNER TO newowner;"
	if got != want {
		t.Errorf("AlterTablespaceOwnerSQL() = %q, want %q", got, want)
	}
}

func TestTablespaceDiffSQL_AddedAndRemoved(t *testing.T) {
	diff := schema.SchemaDiff{
		AddedTablespaces: []schema.Tablespace{
			{Name: "new_ts", Owner: "admin", Location: "/mnt/new"},
		},
		RemovedTablespaces: []schema.Tablespace{
			{Name: "old_ts"},
		},
	}
	stmts := migrate.TablespaceDiffSQL(diff)
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
	if stmts[0] != "CREATE TABLESPACE new_ts OWNER admin LOCATION '/mnt/new';" {
		t.Errorf("unexpected create statement: %q", stmts[0])
	}
	if stmts[1] != "DROP TABLESPACE old_ts;" {
		t.Errorf("unexpected drop statement: %q", stmts[1])
	}
}

func TestTablespaceDiffSQL_Changed(t *testing.T) {
	diff := schema.SchemaDiff{
		ChangedTablespaces: []schema.Tablespace{
			{Name: "my_ts", Owner: "newowner"},
		},
	}
	stmts := migrate.TablespaceDiffSQL(diff)
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	if stmts[0] != "ALTER TABLESPACE my_ts OWNER TO newowner;" {
		t.Errorf("unexpected alter statement: %q", stmts[0])
	}
}
