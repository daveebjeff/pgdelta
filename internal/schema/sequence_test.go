package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseSequence() Sequence {
	return Sequence{
		Schema:    "public",
		Name:      "users_id_seq",
		Start:     1,
		Increment: 1,
		MinValue:  1,
		MaxValue:  9223372036854775807,
		CacheSize: 1,
		Cycle:     false,
	}
}

func TestSequenceFullName(t *testing.T) {
	s := baseSequence()
	assert.Equal(t, "public.users_id_seq", s.FullName())
}

func TestDiffSequences_NoChanges(t *testing.T) {
	s := baseSequence()
	diff := DiffSequences([]Sequence{s}, []Sequence{s})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffSequences_AddedSequence(t *testing.T) {
	s := baseSequence()
	diff := DiffSequences([]Sequence{}, []Sequence{s})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, s, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffSequences_RemovedSequence(t *testing.T) {
	s := baseSequence()
	diff := DiffSequences([]Sequence{s}, []Sequence{})
	assert.Empty(t, diff.Added)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, s, diff.Removed[0])
	assert.Empty(t, diff.Changed)
}

func TestDiffSequences_ChangedSequence(t *testing.T) {
	old := baseSequence()
	newSeq := baseSequence()
	newSeq.Increment = 5
	newSeq.CacheSize = 10

	diff := DiffSequences([]Sequence{old}, []Sequence{newSeq})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, old, diff.Changed[0].Old)
	assert.Equal(t, newSeq, diff.Changed[0].New)
}
