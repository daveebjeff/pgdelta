package migrate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"pgdelta/internal/migrate"
	"pgdelta/internal/schema"
)

func TestAlterSystemSetSQL_Basic(t *testing.T) {
	s := schema.ServerSetting{
		Name:  "work_mem",
		Value: "64MB",
	}
	sql := migrate.AlterSystemSetSQL(s)
	assert.Equal(t, "ALTER SYSTEM SET work_mem = '64MB';", sql)
}

func TestAlterSystemSetSQL_NumericValue(t *testing.T) {
	s := schema.ServerSetting{
		Name:  "max_connections",
		Value: "200",
	}
	sql := migrate.AlterSystemSetSQL(s)
	assert.Equal(t, "ALTER SYSTEM SET max_connections = '200';", sql)
}

func TestAlterSystemResetSQL(t *testing.T) {
	s := schema.ServerSetting{
		Name:  "work_mem",
		Value: "64MB",
	}
	sql := migrate.AlterSystemResetSQL(s)
	assert.Equal(t, "ALTER SYSTEM RESET work_mem;", sql)
}

func TestServerSettingDiffSQL_AddedAndRemoved(t *testing.T) {
	added := schema.ServerSetting{Name: "work_mem", Value: "64MB"}
	removed := schema.ServerSetting{Name: "old_param", Value: "off"}

	diff := schema.ServerSettingDiff{
		Added:   []schema.ServerSetting{added},
		Removed: []schema.ServerSetting{removed},
	}

	sql := migrate.ServerSettingDiffSQL(diff)
	assert.Contains(t, sql, "ALTER SYSTEM SET work_mem = '64MB';")
	assert.Contains(t, sql, "ALTER SYSTEM RESET old_param;")
}

func TestServerSettingDiffSQL_Changed(t *testing.T) {
	old := schema.ServerSetting{Name: "work_mem", Value: "32MB"}
	new := schema.ServerSetting{Name: "work_mem", Value: "128MB"}

	diff := schema.ServerSettingDiff{
		Changed: []schema.ServerSettingChange{
			{Old: old, New: new},
		},
	}

	sql := migrate.ServerSettingDiffSQL(diff)
	assert.Contains(t, sql, "ALTER SYSTEM SET work_mem = '128MB';")
}

func TestServerSettingDiffSQL_Empty(t *testing.T) {
	diff := schema.ServerSettingDiff{}
	sql := migrate.ServerSettingDiffSQL(diff)
	assert.Empty(t, sql)
}
