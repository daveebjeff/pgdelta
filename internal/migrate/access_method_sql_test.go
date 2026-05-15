package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

var baseAM = schema.AccessMethod{
	Name:    "my_am",
	Type:    "index",
	Handler: "my_am_handler",
}

func TestCreateAccessMethodSQL(t *testing.T) {
	sql := CreateAccessMethodSQL(baseAM)
	assert.Equal(t, "CREATE ACCESS METHOD my_am TYPE INDEX HANDLER my_am_handler;", sql)
}

func TestDropAccessMethodSQL(t *testing.T) {
	sql := DropAccessMethodSQL(baseAM)
	assert.Equal(t, "DROP ACCESS METHOD my_am;", sql)
}

func TestAccessMethodDiffSQL_AddedAndRemoved(t *testing.T) {
	added := schema.AccessMethod{Name: "new_am", Type: "index", Handler: "new_handler"}
	removed := schema.AccessMethod{Name: "old_am", Type: "index", Handler: "old_handler"}

	diff := schema.AccessMethodDiff{
		Added:   []schema.AccessMethod{added},
		Removed: []schema.AccessMethod{removed},
	}

	stmts := AccessMethodDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Equal(t, "DROP ACCESS METHOD old_am;", stmts[0])
	assert.Equal(t, "CREATE ACCESS METHOD new_am TYPE INDEX HANDLER new_handler;", stmts[1])
}

func TestAccessMethodDiffSQL_Changed(t *testing.T) {
	changed := schema.AccessMethod{Name: "my_am", Type: "index", Handler: "updated_handler"}
	diff := schema.AccessMethodDiff{
		Changed: []schema.AccessMethod{changed},
	}

	stmts := AccessMethodDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Equal(t, "DROP ACCESS METHOD my_am;", stmts[0])
	assert.Equal(t, "CREATE ACCESS METHOD my_am TYPE INDEX HANDLER updated_handler;", stmts[1])
}

func TestAccessMethodDiffSQL_Empty(t *testing.T) {
	stmts := AccessMethodDiffSQL(schema.AccessMethodDiff{})
	assert.Empty(t, stmts)
}
