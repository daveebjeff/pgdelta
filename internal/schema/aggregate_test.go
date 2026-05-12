package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseAggregate() Aggregate {
	return Aggregate{
		Schema:    "public",
		Name:      "my_agg",
		ArgTypes:  []string{"integer"},
		SFuncName: "int4pl",
		SType:     "integer",
	}
}

func TestAggregateFullName(t *testing.T) {
	a := baseAggregate()
	assert.Equal(t, "public.my_agg(integer)", a.FullName())
}

func TestDiffAggregates_NoChanges(t *testing.T) {
	a := baseAggregate()
	diff := DiffAggregates([]Aggregate{a}, []Aggregate{a})
	assert.True(t, diff.IsEmpty())
}

func TestDiffAggregates_AddedAggregate(t *testing.T) {
	a := baseAggregate()
	diff := DiffAggregates([]Aggregate{}, []Aggregate{a})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, a.FullName(), diff.Added[0].FullName())
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffAggregates_RemovedAggregate(t *testing.T) {
	a := baseAggregate()
	diff := DiffAggregates([]Aggregate{a}, []Aggregate{})
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, a.FullName(), diff.Removed[0].FullName())
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffAggregates_ChangedAggregate(t *testing.T) {
	old := baseAggregate()
	new := baseAggregate()
	new.SFuncName = "int4mi"
	diff := DiffAggregates([]Aggregate{old}, []Aggregate{new})
	assert.Len(t, diff.Changed, 1)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}

func TestDiffAggregates_ChangedFinalFunc(t *testing.T) {
	old := baseAggregate()
	new := baseAggregate()
	f := "my_final"
	new.FinalFunc = &f
	diff := DiffAggregates([]Aggregate{old}, []Aggregate{new})
	assert.Len(t, diff.Changed, 1)
}
