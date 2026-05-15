package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseFDW = ForeignDataWrapper{
	Name:      "myfdw",
	Handler:   "myfdw_handler",
	Validator: "myfdw_validator",
	Owner:     "postgres",
	Options:   map[string]string{"host": "localhost"},
}

func TestFDWFullName(t *testing.T) {
	assert.Equal(t, "fdw:myfdw", baseFDW.FullName())
}

func TestDiffFDWs_NoChanges(t *testing.T) {
	added, removed, changed := DiffForeignDataWrappers([]ForeignDataWrapper{baseFDW}, []ForeignDataWrapper{baseFDW})
	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Empty(t, changed)
}

func TestDiffFDWs_AddedFDW(t *testing.T) {
	added, removed, changed := DiffForeignDataWrappers(nil, []ForeignDataWrapper{baseFDW})
	assert.Len(t, added, 1)
	assert.Empty(t, removed)
	assert.Empty(t, changed)
}

func TestDiffFDWs_RemovedFDW(t *testing.T) {
	added, removed, changed := DiffForeignDataWrappers([]ForeignDataWrapper{baseFDW}, nil)
	assert.Empty(t, added)
	assert.Len(t, removed, 1)
	assert.Empty(t, changed)
}

func TestDiffFDWs_ChangedFDW(t *testing.T) {
	updated := baseFDW
	updated.Handler = "new_handler"
	added, removed, changed := DiffForeignDataWrappers([]ForeignDataWrapper{baseFDW}, []ForeignDataWrapper{updated})
	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Len(t, changed, 1)
}

func TestDiffFDWs_ChangedOptions(t *testing.T) {
	updated := baseFDW
	updated.Options = map[string]string{"host": "remotehost"}
	added, removed, changed := DiffForeignDataWrappers([]ForeignDataWrapper{baseFDW}, []ForeignDataWrapper{updated})
	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Len(t, changed, 1)
}
