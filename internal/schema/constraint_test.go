package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseConstraint() Constraint {
	return Constraint{
		Schema:     "public",
		Table:      "orders",
		Name:       "orders_pkey",
		Type:       "PRIMARY KEY",
		Definition: "PRIMARY KEY (id)",
	}
}

func TestConstraintFullName(t *testing.T) {
	c := baseConstraint()
	assert.Equal(t, "public.orders.orders_pkey", c.FullName())
}

func TestDiffConstraints_NoChanges(t *testing.T) {
	c := baseConstraint()
	diff := DiffConstraints([]Constraint{c}, []Constraint{c})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffConstraints_AddedConstraint(t *testing.T) {
	c := baseConstraint()
	newC := Constraint{
		Schema:     "public",
		Table:      "orders",
		Name:       "orders_status_check",
		Type:       "CHECK",
		Definition: "CHECK (status IN ('pending', 'done'))",
	}
	diff := DiffConstraints([]Constraint{c}, []Constraint{c, newC})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, newC, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffConstraints_RemovedConstraint(t *testing.T) {
	c := baseConstraint()
	diff := DiffConstraints([]Constraint{c}, []Constraint{})
	assert.Empty(t, diff.Added)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, c, diff.Removed[0])
	assert.Empty(t, diff.Changed)
}

func TestDiffConstraints_ChangedConstraint(t *testing.T) {
	old := baseConstraint()
	updated := old
	updated.Definition = "PRIMARY KEY (id, tenant_id)"
	diff := DiffConstraints([]Constraint{old}, []Constraint{updated})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, updated, diff.Changed[0])
}
