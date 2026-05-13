package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseCollation = Collation{
	Schema:        "public",
	Name:          "my_collation",
	Provider:      "icu",
	Locale:        "en-US",
	Deterministic: true,
}

func TestCollationFullName(t *testing.T) {
	assert.Equal(t, "public.my_collation", baseCollation.FullName())
}

func TestDiffCollations_NoChanges(t *testing.T) {
	diff := DiffCollations([]Collation{baseCollation}, []Collation{baseCollation})
	assert.True(t, diff.IsEmpty())
}

func TestDiffCollations_AddedCollation(t *testing.T) {
	diff := DiffCollations(nil, []Collation{baseCollation})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, baseCollation, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffCollations_RemovedCollation(t *testing.T) {
	diff := DiffCollations([]Collation{baseCollation}, nil)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseCollation, diff.Removed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffCollations_ChangedCollation(t *testing.T) {
	modified := baseCollation
	modified.Locale = "fr-FR"

	diff := DiffCollations([]Collation{baseCollation}, []Collation{modified})
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, modified, diff.Changed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}

func TestDiffCollations_ChangedDeterministic(t *testing.T) {
	modified := baseCollation
	modified.Deterministic = false

	diff := DiffCollations([]Collation{baseCollation}, []Collation{modified})
	assert.Len(t, diff.Changed, 1)
	assert.False(t, diff.Changed[0].Deterministic)
}
