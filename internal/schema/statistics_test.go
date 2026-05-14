package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseStatistic = Statistic{
	Schema:    "public",
	Name:      "orders_stats",
	TableName: "orders",
	Columns:   []string{"customer_id", "product_id"},
	Kinds:     []string{"dependencies", "ndistinct"},
}

func TestStatisticFullName(t *testing.T) {
	assert.Equal(t, "public.orders_stats", baseStatistic.FullName())
}

func TestDiffStatistics_NoChanges(t *testing.T) {
	old := []Statistic{baseStatistic}
	new := []Statistic{baseStatistic}
	diff := DiffStatistics(old, new)
	assert.True(t, diff.IsEmpty())
}

func TestDiffStatistics_AddedStatistic(t *testing.T) {
	old := []Statistic{}
	new := []Statistic{baseStatistic}
	diff := DiffStatistics(old, new)
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, baseStatistic, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffStatistics_RemovedStatistic(t *testing.T) {
	old := []Statistic{baseStatistic}
	new := []Statistic{}
	diff := DiffStatistics(old, new)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseStatistic, diff.Removed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffStatistics_ChangedStatistic(t *testing.T) {
	old := []Statistic{baseStatistic}
	modified := baseStatistic
	modified.Kinds = []string{"dependencies"}
	new := []Statistic{modified}
	diff := DiffStatistics(old, new)
	assert.Len(t, diff.Changed, 1)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}

func TestDiffStatistics_ChangedColumns(t *testing.T) {
	old := []Statistic{baseStatistic}
	modified := baseStatistic
	modified.Columns = []string{"customer_id"}
	new := []Statistic{modified}
	diff := DiffStatistics(old, new)
	assert.Len(t, diff.Changed, 1)
}
