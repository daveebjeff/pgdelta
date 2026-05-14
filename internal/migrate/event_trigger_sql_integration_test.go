package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventTriggerMigration_FullCycle(t *testing.T) {
	db := openTestDB(t)

	_, err := db.Exec(`
		CREATE OR REPLACE FUNCTION pgdelta_test_event_func()
		RETURNS event_trigger LANGUAGE plpgsql AS $$
		BEGIN
			RAISE NOTICE 'event trigger fired';
		END;
		$$;
	`)
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec(`DROP EVENT TRIGGER IF EXISTS pgdelta_test_trigger;`)
		db.Exec(`DROP FUNCTION IF EXISTS pgdelta_test_event_func();`)
	})

	et := schema.EventTrigger{
		Name:     "pgdelta_test_trigger",
		Event:    "ddl_command_end",
		FuncName: "pgdelta_test_event_func",
		Enabled:  "ENABLE",
		Tags:     []string{"CREATE TABLE"},
	}

	createSql := migrate.CreateEventTriggerSQL(et)
	_, err = db.Exec(createSql)
	require.NoError(t, err)

	var count int
	err = db.QueryRow(`SELECT count(*) FROM pg_event_trigger WHERE evtname = 'pgdelta_test_trigger'`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	dropSql := migrate.DropEventTriggerSQL(et)
	_, err = db.Exec(dropSql)
	require.NoError(t, err)

	err = db.QueryRow(`SELECT count(*) FROM pg_event_trigger WHERE evtname = 'pgdelta_test_trigger'`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
