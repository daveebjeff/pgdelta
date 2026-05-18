package schema_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
)

func baseUserMapping() schema.UserMapping {
	return schema.UserMapping{
		User:   "alice",
		Server: "myserver",
		Options: map[string]string{
			"user":     "remote_alice",
			"password": "secret",
		},
	}
}

func TestUserMappingFullName(t *testing.T) {
	m := baseUserMapping()
	if got := m.FullName(); got != "alice@myserver" {
		t.Errorf("expected alice@myserver, got %s", got)
	}
}

func TestDiffUserMappings_NoChanges(t *testing.T) {
	m := baseUserMapping()
	diff := schema.DiffUserMappings([]schema.UserMapping{m}, []schema.UserMapping{m})
	if len(diff.Added)+len(diff.Removed)+len(diff.Changed) != 0 {
		t.Errorf("expected no diff, got %+v", diff)
	}
}

func TestDiffUserMappings_AddedUserMapping(t *testing.T) {
	m := baseUserMapping()
	diff := schema.DiffUserMappings(nil, []schema.UserMapping{m})
	if len(diff.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(diff.Added))
	}
}

func TestDiffUserMappings_RemovedUserMapping(t *testing.T) {
	m := baseUserMapping()
	diff := schema.DiffUserMappings([]schema.UserMapping{m}, nil)
	if len(diff.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(diff.Removed))
	}
}

func TestDiffUserMappings_ChangedUserMapping(t *testing.T) {
	old := baseUserMapping()
	new := baseUserMapping()
	new.Options["password"] = "newsecret"
	diff := schema.DiffUserMappings([]schema.UserMapping{old}, []schema.UserMapping{new})
	if len(diff.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(diff.Changed))
	}
}
