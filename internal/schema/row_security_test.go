package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseRowSecurity = RowSecurity{
	Schema:  "public",
	Table:   "orders",
	Enabled: true,
	Forced:  false,
}

func TestRowSecurityFullName(t *testing.T) {
	rs := baseRowSecurity
	assert.Equal(t, "public.orders.row_security", rs.FullName())
}

func TestDiffRowSecurity_NoChanges(t *testing.T) {
	rs := baseRowSecurity
	diff := DiffRowSecurity([]RowSecurity{rs}, []RowSecurity{rs})
	assert.True(t, diff.IsEmpty())
}

func TestDiffRowSecurity_AddedRowSecurity(t *testing.T) {
	rs := baseRowSecurity
	diff := DiffRowSecurity(nil, []RowSecurity{rs})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, rs, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffRowSecurity_RemovedRowSecurity(t *testing.T) {
	rs := baseRowSecurity
	diff := DiffRowSecurity([]RowSecurity{rs}, nil)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, rs, diff.Removed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffRowSecurity_ChangedEnabled(t *testing.T) {
	before := baseRowSecurity
	after := baseRowSecurity
	after.Enabled = false
	diff := DiffRowSecurity([]RowSecurity{before}, []RowSecurity{after})
	assert.Len(t, diff.Changed, 1)
	assert.False(t, diff.Changed[0].Enabled)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}

func TestDiffRowSecurity_ChangedForced(t *testing.T) {
	before := baseRowSecurity
	after := baseRowSecurity
	after.Forced = true
	diff := DiffRowSecurity([]RowSecurity{before}, []RowSecurity{after})
	assert.Len(t, diff.Changed, 1)
	assert.True(t, diff.Changed[0].Forced)
}
