package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var basePublication = Publication{
	Name:      "my_pub",
	AllTables: false,
	Insert:    true,
	Update:    true,
	Delete:    false,
	Truncate:  false,
	Tables:    []string{"public.orders"},
}

func TestPublicationFullName(t *testing.T) {
	p := basePublication
	assert.Equal(t, "publication:my_pub", p.FullName())
}

func TestDiffPublications_NoChanges(t *testing.T) {
	old := []Publication{basePublication}
	new := []Publication{basePublication}
	diff := DiffPublications(old, new)
	assert.True(t, diff.IsEmpty())
}

func TestDiffPublications_AddedPublication(t *testing.T) {
	newPub := Publication{Name: "new_pub", AllTables: true, Insert: true, Update: true}
	diff := DiffPublications([]Publication{basePublication}, []Publication{basePublication, newPub})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, "new_pub", diff.Added[0].Name)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffPublications_RemovedPublication(t *testing.T) {
	diff := DiffPublications([]Publication{basePublication}, []Publication{})
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, "my_pub", diff.Removed[0].Name)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffPublications_ChangedPublication(t *testing.T) {
	changed := basePublication
	changed.Delete = true
	diff := DiffPublications([]Publication{basePublication}, []Publication{changed})
	assert.Len(t, diff.Changed, 1)
	assert.True(t, diff.Changed[0].Delete)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}

func TestDiffPublications_ChangedTables(t *testing.T) {
	changed := basePublication
	changed.Tables = []string{"public.orders", "public.items"}
	diff := DiffPublications([]Publication{basePublication}, []Publication{changed})
	assert.Len(t, diff.Changed, 1)
}
