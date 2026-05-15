package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

var baseFDW = schema.ForeignDataWrapper{
	Name:      "myfdw",
	Handler:   "myfdw_handler",
	Validator: "myfdw_validator",
	Owner:     "postgres",
	Options:   map[string]string{"encoding": "utf8"},
}

func TestCreateFDWSQL_WithHandlerAndValidator(t *testing.T) {
	sql := CreateForeignDataWrapperSQL(baseFDW)
	assert.Contains(t, sql, "CREATE FOREIGN DATA WRAPPER myfdw")
	assert.Contains(t, sql, "HANDLER myfdw_handler")
	assert.Contains(t, sql, "VALIDATOR myfdw_validator")
	assert.Contains(t, sql, "OPTIONS (encoding 'utf8')")
}

func TestCreateFDWSQL_NoHandlerNoValidator(t *testing.T) {
	f := schema.ForeignDataWrapper{Name: "simplefdw"}
	sql := CreateForeignDataWrapperSQL(f)
	assert.Contains(t, sql, "NO HANDLER")
	assert.Contains(t, sql, "NO VALIDATOR")
	assert.NotContains(t, sql, "OPTIONS")
}

func TestDropFDWSQL(t *testing.T) {
	sql := DropForeignDataWrapperSQL(baseFDW)
	assert.Equal(t, "DROP FOREIGN DATA WRAPPER myfdw;", sql)
}

func TestAlterFDWSQL_WithOptions(t *testing.T) {
	sql := AlterForeignDataWrapperSQL(baseFDW)
	assert.Contains(t, sql, "ALTER FOREIGN DATA WRAPPER myfdw")
	assert.Contains(t, sql, "OPTIONS (SET encoding 'utf8')")
}

func TestFDWDiffSQL_AddedAndRemoved(t *testing.T) {
	old := []schema.ForeignDataWrapper{baseFDW}
	newFDW := schema.ForeignDataWrapper{Name: "newfdw", Handler: "newfdw_handler"}
	new := []schema.ForeignDataWrapper{newFDW}
	stmts := FDWDiffSQL(old, new)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "DROP FOREIGN DATA WRAPPER myfdw")
	assert.Contains(t, stmts[1], "CREATE FOREIGN DATA WRAPPER newfdw")
}

func TestFDWDiffSQL_Changed(t *testing.T) {
	updated := baseFDW
	updated.Handler = "updated_handler"
	stmts := FDWDiffSQL([]schema.ForeignDataWrapper{baseFDW}, []schema.ForeignDataWrapper{updated})
	assert.Len(t, stmts, 1)
	assert.Contains(t, stmts[0], "ALTER FOREIGN DATA WRAPPER myfdw")
	assert.Contains(t, stmts[0], "HANDLER updated_handler")
}
