package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseMaterializedView() MaterializedView {
	return MaterializedView{
		Schema:     "public",
		Name:       "mv_orders",
		Definition: "SELECT id, total FROM orders WHERE total > 0",
		WithData:    true,
	}
}

func TestMaterializedViewFullName(t *testing.T) {
	mv := baseMaterializedView()
	assert.Equal(t, "public.mv_orders", mv.FullName())
}

func TestDiffMaterializedViews_NoChanges(t *testing.T) {
	mv := baseMaterializedView()
	diff := DiffMaterializedViews([]MaterializedView{mv}, []MaterializedView{mv})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffMaterializedViews_AddedView(t *testing.T) {
	mv := baseMaterializedView()
	diff := DiffMaterializedViews(nil, []MaterializedView{mv})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, mv, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffMaterializedViews_RemovedView(t *testing.T) {
	mv := baseMaterializedView()
	diff := DiffMaterializedViews([]MaterializedView{mv}, nil)
	assert.Empty(t, diff.Added)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, mv, diff.Removed[0])
	assert.Empty(t, diff.Changed)
}

func TestDiffMaterializedViews_ChangedDefinition(t *testing.T) {
	old := baseMaterializedView()
	new := baseMaterializedView()
	new.Definition = "SELECT id, total FROM orders WHERE total > 100"
	diff := DiffMaterializedViews([]MaterializedView{old}, []MaterializedView{new})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, new, diff.Changed[0])
}

func TestDiffMaterializedViews_ChangedWithData(t *testing.T) {
	old := baseMaterializedView()
	new := baseMaterializedView()
	new.WithData = false
	diff := DiffMaterializedViews([]MaterializedView{old}, []MaterializedView{new})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
}
