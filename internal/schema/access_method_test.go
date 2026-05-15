package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseAccessMethod = AccessMethod{
	Name:    "my_am",
	Type:    "index",
	Handler: "my_am_handler",
}

func TestAccessMethodFullName(t *testing.T) {
	am := baseAccessMethod
	assert.Equal(t, "my_am", am.FullName())
}

func TestDiffAccessMethods_NoChanges(t *testing.T) {
	am := baseAccessMethod
	diff := DiffAccessMethods([]AccessMethod{am}, []AccessMethod{am})
	assert.True(t, diff.IsEmpty())
}

func TestDiffAccessMethods_AddedAccessMethod(t *testing.T) {
	am := baseAccessMethod
	diff := DiffAccessMethods(nil, []AccessMethod{am})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, am, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffAccessMethods_RemovedAccessMethod(t *testing.T) {
	am := baseAccessMethod
	diff := DiffAccessMethods([]AccessMethod{am}, nil)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, am, diff.Removed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffAccessMethods_ChangedAccessMethod(t *testing.T) {
	old := baseAccessMethod
	new := baseAccessMethod
	new.Handler = "other_handler"
	diff := DiffAccessMethods([]AccessMethod{old}, []AccessMethod{new})
	assert.Len(t, diff.Changed, 1)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}
