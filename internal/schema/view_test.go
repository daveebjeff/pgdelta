package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseView = View{
	Schema:     "public",
	Name:       "active_users",
	Definition: "SELECT id, name FROM users WHERE active = true",
}

func TestViewFullName(t *testing.T) {
	assert.Equal(t, "public.active_users", baseView.FullName())
}

func TestDiffViews_NoChanges(t *testing.T) {
	views := []View{baseView}
	diff := DiffViews(views, views)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffViews_AddedView(t *testing.T) {
	newView := View{
		Schema:     "public",
		Name:       "inactive_users",
		Definition: "SELECT id, name FROM users WHERE active = false",
	}
	diff := DiffViews([]View{baseView}, []View{baseView, newView})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, newView, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffViews_RemovedView(t *testing.T) {
	diff := DiffViews([]View{baseView}, []View{})
	assert.Empty(t, diff.Added)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseView, diff.Removed[0])
	assert.Empty(t, diff.Changed)
}

func TestDiffViews_ChangedView(t *testing.T) {
	updatedView := View{
		Schema:     "public",
		Name:       "active_users",
		Definition: "SELECT id, name, email FROM users WHERE active = true",
	}
	diff := DiffViews([]View{baseView}, []View{updatedView})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, baseView, diff.Changed[0].Old)
	assert.Equal(t, updatedView, diff.Changed[0].New)
}
