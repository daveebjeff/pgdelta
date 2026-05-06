package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

func baseSeq() schema.Sequence {
	return schema.Sequence{
		Schema:    "public",
		Name:      "orders_id_seq",
		Start:     1,
		Increment: 1,
		MinValue:  1,
		MaxValue:  9223372036854775807,
		CacheSize: 1,
		Cycle:     false,
	}
}

func TestCreateSequenceSQL(t *testing.T) {
	s := baseSeq()
	sql := CreateSequenceSQL(s)
	assert.Contains(t, sql, "CREATE SEQUENCE public.orders_id_seq")
	assert.Contains(t, sql, "START WITH 1")
	assert.Contains(t, sql, "INCREMENT BY 1")
	assert.Contains(t, sql, "NO CYCLE")
}

func TestCreateSequenceSQL_Cycle(t *testing.T) {
	s := baseSeq()
	s.Cycle = true
	sql := CreateSequenceSQL(s)
	assert.Contains(t, sql, "CYCLE")
	assert.NotContains(t, sql, "NO CYCLE")
}

func TestDropSequenceSQL(t *testing.T) {
	s := baseSeq()
	sql := DropSequenceSQL(s)
	assert.Equal(t, "DROP SEQUENCE public.orders_id_seq;", sql)
}

func TestAlterSequenceSQL_IncrementAndCache(t *testing.T) {
	old := baseSeq()
	newSeq := baseSeq()
	newSeq.Increment = 10
	newSeq.CacheSize = 20

	sql := AlterSequenceSQL(old, newSeq)
	assert.Contains(t, sql, "ALTER SEQUENCE public.orders_id_seq")
	assert.Contains(t, sql, "INCREMENT BY 10")
	assert.Contains(t, sql, "CACHE 20")
}

func TestAlterSequenceSQL_CycleChange(t *testing.T) {
	old := baseSeq()
	newSeq := baseSeq()
	newSeq.Cycle = true

	sql := AlterSequenceSQL(old, newSeq)
	assert.Contains(t, sql, "ALTER SEQUENCE public.orders_id_seq")
	assert.Contains(t, sql, "CYCLE")
	assert.NotContains(t, sql, "NO CYCLE")
}

func TestAlterSequenceSQL_NoChanges(t *testing.T) {
	s := baseSeq()
	sql := AlterSequenceSQL(s, s)
	assert.Empty(t, sql)
}

func TestSequenceDiffSQL(t *testing.T) {
	added := baseSeq()
	removed := baseSeq()
	removed.Name = "old_seq"

	diff := schema.SequenceDiff{
		Added:   []schema.Sequence{added},
		Removed: []schema.Sequence{removed},
	}

	stmts := SequenceDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "CREATE SEQUENCE")
	assert.Contains(t, stmts[1], "DROP SEQUENCE")
}
