package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

func baseAgg() schema.Aggregate {
	return schema.Aggregate{
		Schema:    "public",
		Name:      "sum_int",
		ArgTypes:  []string{"integer"},
		SFuncName: "int4pl",
		SType:     "integer",
	}
}

func strPtrAgg(s string) *string { return &s }

func TestCreateAggregateSQL_Basic(t *testing.T) {
	a := baseAgg()
	sql := CreateAggregateSQL(a)
	assert.Contains(t, sql, "CREATE AGGREGATE public.sum_int(integer)")
	assert.Contains(t, sql, "SFUNC = int4pl")
	assert.Contains(t, sql, "STYPE = integer")
}

func TestCreateAggregateSQL_WithInitCond(t *testing.T) {
	a := baseAgg()
	a.InitCond = strPtrAgg("0")
	sql := CreateAggregateSQL(a)
	assert.Contains(t, sql, "INITCOND = '0'")
}

func TestCreateAggregateSQL_WithFinalFunc(t *testing.T) {
	a := baseAgg()
	a.FinalFunc = strPtrAgg("my_final")
	sql := CreateAggregateSQL(a)
	assert.Contains(t, sql, "FINALFUNC = my_final")
}

func TestDropAggregateSQL(t *testing.T) {
	a := baseAgg()
	sql := DropAggregateSQL(a)
	assert.Equal(t, "DROP AGGREGATE public.sum_int(integer);", sql)
}

func TestAggregateDiffSQL_AddedAndRemoved(t *testing.T) {
	add := baseAgg()
	add.Name = "new_agg"
	remove := baseAgg()
	remove.Name = "old_agg"

	diff := schema.AggregateDiff{
		Added:   []schema.Aggregate{add},
		Removed: []schema.Aggregate{remove},
	}

	stmts := AggregateDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "DROP AGGREGATE")
	assert.Contains(t, stmts[1], "CREATE AGGREGATE")
}

func TestAggregateDiffSQL_Changed(t *testing.T) {
	a := baseAgg()
	a.SFuncName = "int4mi"
	diff := schema.AggregateDiff{Changed: []schema.Aggregate{a}}
	stmts := AggregateDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "DROP AGGREGATE")
	assert.Contains(t, stmts[1], "CREATE AGGREGATE")
}
