package migrate

import (
	"strings"
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
)

var sampleIndex = schema.Index{
	SchemaName: "public",
	TableName:  "users",
	Name:       "idx_users_email",
	Columns:    []string{"email"},
	Unique:     true,
	Method:     schema.IndexMethodBTree,
}

func TestCreateIndexSQL_Unique(t *testing.T) {
	sql := CreateIndexSQL(sampleIndex)
	if !strings.Contains(sql, "CREATE UNIQUE INDEX") {
		t.Errorf("expected UNIQUE in SQL, got: %s", sql)
	}
	if !strings.Contains(sql, "idx_users_email") {
		t.Errorf("expected index name in SQL, got: %s", sql)
	}
	if !strings.Contains(sql, "public.users") {
		t.Errorf("expected table reference in SQL, got: %s", sql)
	}
	if !strings.Contains(sql, "(email)") {
		t.Errorf("expected column list in SQL, got: %s", sql)
	}
}

func TestCreateIndexSQL_NonUnique(t *testing.T) {
	idx := sampleIndex
	idx.Unique = false
	sql := CreateIndexSQL(idx)
	if strings.Contains(sql, "UNIQUE") {
		t.Errorf("did not expect UNIQUE in SQL, got: %s", sql)
	}
}

func TestCreateIndexSQL_NonBTreeMethod(t *testing.T) {
	idx := sampleIndex
	idx.Method = schema.IndexMethodGIN
	sql := CreateIndexSQL(idx)
	if !strings.Contains(sql, "USING gin") {
		t.Errorf("expected USING gin in SQL, got: %s", sql)
	}
}

func TestCreateIndexSQL_MultipleColumns(t *testing.T) {
	idx := sampleIndex
	idx.Columns = []string{"first_name", "last_name"}
	sql := CreateIndexSQL(idx)
	if !strings.Contains(sql, "(first_name, last_name)") {
		t.Errorf("expected multi-column list in SQL, got: %s", sql)
	}
}

func TestDropIndexSQL(t *testing.T) {
	sql := DropIndexSQL(sampleIndex)
	expected := "DROP INDEX public.idx_users_email;"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
}

func TestIndexDiffSQL(t *testing.T) {
	newIdx := schema.Index{
		SchemaName: "public",
		TableName:  "users",
		Name:       "idx_users_name",
		Columns:    []string{"name"},
		Unique:     false,
		Method:     schema.IndexMethodBTree,
	}
	diff := schema.IndexDiff{
		Added:   []schema.Index{newIdx},
		Removed: []schema.Index{sampleIndex},
	}
	stmts := IndexDiffSQL(diff)
	if len(stmts) != 2 {
		t.Errorf("expected 2 statements, got %d", len(stmts))
	}
	if !strings.HasPrefix(stmts[0], "DROP") {
		t.Errorf("expected first statement to be DROP, got: %s", stmts[0])
	}
	if !strings.HasPrefix(stmts[1], "CREATE") {
		t.Errorf("expected second statement to be CREATE, got: %s", stmts[1])
	}
}
