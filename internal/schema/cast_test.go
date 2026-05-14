package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseCast = Cast{
	SourceType:   "integer",
	TargetType:   "text",
	FunctionName: "pg_catalog.text",
	Schema:       "pg_catalog",
	CastContext:  "a",
}

func TestCastFullName(t *testing.T) {
	c := baseCast
	assert.Equal(t, "(integer AS text)", c.FullName())
}

func TestDiffCasts_NoChanges(t *testing.T) {
	casts := []Cast{baseCast}
	diff := DiffCasts(casts, casts)
	assert.True(t, diff.IsEmpty())
}

func TestDiffCasts_AddedCast(t *testing.T) {
	newCast := Cast{
		SourceType:   "float",
		TargetType:   "integer",
		FunctionName: "pg_catalog.int4",
		Schema:       "pg_catalog",
		CastContext:  "e",
	}
	diff := DiffCasts([]Cast{baseCast}, []Cast{baseCast, newCast})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, newCast, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffCasts_RemovedCast(t *testing.T) {
	diff := DiffCasts([]Cast{baseCast}, []Cast{})
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseCast, diff.Removed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffCasts_ChangedCast(t *testing.T) {
	modified := baseCast
	modified.CastContext = "i"
	diff := DiffCasts([]Cast{baseCast}, []Cast{modified})
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, modified, diff.Changed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}
