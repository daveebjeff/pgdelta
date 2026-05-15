package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

var baseStatistic = schema.Statistic{
	Schema:  "public",
	Name:    "my_stats",
	Table:   "orders",
	Columns: []string{"customer_id", "product_id"},
	Kinds:   []string{"dependencies", "ndistinct"},
}

func TestCreateStatisticsSQL_Basic(t *testing.T) {
	sql := migrate.CreateStatisticsSQL(baseStatistic)
	expected := `CREATE STATISTICS public.my_stats (dependencies, ndistinct) ON customer_id, product_id FROM public.orders;`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestCreateStatisticsSQL_SingleColumn(t *testing.T) {
	s := schema.Statistic{
		Schema:  "app",
		Name:    "col_stats",
		Table:   "events",
		Columns: []string{"event_type"},
		Kinds:   []string{"mcv"},
	}
	sql := migrate.CreateStatisticsSQL(s)
	expected := `CREATE STATISTICS app.col_stats (mcv) ON event_type FROM app.events;`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestDropStatisticsSQL(t *testing.T) {
	sql := migrate.DropStatisticsSQL(baseStatistic)
	expected := `DROP STATISTICS public.my_stats;`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestStatisticsDiffSQL_AddedAndRemoved(t *testing.T) {
	added := schema.Statistic{
		Schema:  "public",
		Name:    "new_stats",
		Table:   "users",
		Columns: []string{"age", "region"},
		Kinds:   []string{"ndistinct"},
	}
	removed := baseStatistic

	diff := schema.StatisticsDiff{
		Added:   []schema.Statistic{added},
		Removed: []schema.Statistic{removed},
	}

	sqls := migrate.StatisticsDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(sqls))
	}
	if sqls[0] != migrate.DropStatisticsSQL(removed) {
		t.Errorf("expected drop first, got: %s", sqls[0])
	}
	if sqls[1] != migrate.CreateStatisticsSQL(added) {
		t.Errorf("expected create second, got: %s", sqls[1])
	}
}

func TestStatisticsDiffSQL_Changed(t *testing.T) {
	changed := baseStatistic
	changed.Kinds = []string{"dependencies"}

	diff := schema.StatisticsDiff{
		Changed: []schema.StatisticsChange{
			{Old: baseStatistic, New: changed},
		},
	}

	sqls := migrate.StatisticsDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements for changed, got %d", len(sqls))
	}
}
