package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseForeignKey = ForeignKey{
	Schema:     "public",
	Table:      "orders",
	Name:       "fk_orders_user_id",
	Columns:    []string{"user_id"},
	RefSchema:  "public",
	RefTable:   "users",
	RefColumns: []string{"id"},
	OnDelete:   "CASCADE",
	OnUpdate:   "NO ACTION",
}

func TestForeignKeyFullName(t *testing.T) {
	assert.Equal(t, "public.orders.fk_orders_user_id", baseForeignKey.FullName())
}

func TestDiffForeignKeys_NoChanges(t *testing.T) {
	diff := DiffForeignKeys([]ForeignKey{baseForeignKey}, []ForeignKey{baseForeignKey})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffForeignKeys_AddedForeignKey(t *testing.T) {
	diff := DiffForeignKeys([]ForeignKey{}, []ForeignKey{baseForeignKey})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, baseForeignKey, diff.Added[0])
	assert.Empty(t, diff.Removed)
}

func TestDiffForeignKeys_RemovedForeignKey(t *testing.T) {
	diff := DiffForeignKeys([]ForeignKey{baseForeignKey}, []ForeignKey{})
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseForeignKey, diff.Removed[0])
	assert.Empty(t, diff.Added)
}

func TestDiffForeignKeys_ChangedForeignKey(t *testing.T) {
	modified := baseForeignKey
	modified.OnDelete = "RESTRICT"

	diff := DiffForeignKeys([]ForeignKey{baseForeignKey}, []ForeignKey{modified})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, "RESTRICT", diff.Changed[0].OnDelete)
}

func TestDiffForeignKeys_ChangedColumns(t *testing.T) {
	modified := baseForeignKey
	modified.Columns = []string{"user_id", "org_id"}
	modified.RefColumns = []string{"id", "org_id"}

	diff := DiffForeignKeys([]ForeignKey{baseForeignKey}, []ForeignKey{modified})
	assert.Len(t, diff.Changed, 1)
}
