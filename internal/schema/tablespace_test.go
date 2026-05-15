package schema_test

import (
	"testing"

	"github.com/your-org/pgdelta/internal/schema"
)

var baseTablespace = schema.Tablespace{
	Name:     "fast_storage",
	Owner:    "postgres",
	Location: "/mnt/ssd/pgdata",
}

func TestTablespaceFullName(t *testing.T) {
	ts := baseTablespace
	if ts.FullName() != "fast_storage" {
		t.Errorf("expected 'fast_storage', got %q", ts.FullName())
	}
}

func TestDiffTablespaces_NoChanges(t *testing.T) {
	old := []schema.Tablespace{baseTablespace}
	new := []schema.Tablespace{baseTablespace}
	diff := schema.DiffTablespaces(old, new)
	if !diff.IsEmpty() {
		t.Errorf("expected no diff, got %+v", diff)
	}
}

func TestDiffTablespaces_AddedTablespace(t *testing.T) {
	old := []schema.Tablespace{}
	new := []schema.Tablespace{baseTablespace}
	diff := schema.DiffTablespaces(old, new)
	if len(diff.Added) != 1 || diff.Added[0].Name != "fast_storage" {
		t.Errorf("expected one added tablespace, got %+v", diff)
	}
}

func TestDiffTablespaces_RemovedTablespace(t *testing.T) {
	old := []schema.Tablespace{baseTablespace}
	new := []schema.Tablespace{}
	diff := schema.DiffTablespaces(old, new)
	if len(diff.Removed) != 1 || diff.Removed[0].Name != "fast_storage" {
		t.Errorf("expected one removed tablespace, got %+v", diff)
	}
}

func TestDiffTablespaces_ChangedOwner(t *testing.T) {
	old := []schema.Tablespace{baseTablespace}
	changed := baseTablespace
	changed.Owner = "admin"
	new := []schema.Tablespace{changed}
	diff := schema.DiffTablespaces(old, new)
	if len(diff.Changed) != 1 {
		t.Errorf("expected one changed tablespace, got %+v", diff)
	}
}

func TestDiffTablespaces_ChangedLocation(t *testing.T) {
	old := []schema.Tablespace{baseTablespace}
	changed := baseTablespace
	changed.Location = "/mnt/nvme/pgdata"
	new := []schema.Tablespace{changed}
	diff := schema.DiffTablespaces(old, new)
	if len(diff.Changed) != 1 {
		t.Errorf("expected one changed tablespace, got %+v", diff)
	}
}
