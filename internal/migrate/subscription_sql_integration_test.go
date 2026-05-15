package migrate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pgdelta/internal/migrate"
	"pgdelta/internal/schema"
)

func TestSubscriptionMigration_FullCycle(t *testing.T) {
	sub := schema.Subscription{
		Name:              "test_sub",
		ConnectionInfo:    "host=localhost port=5432 dbname=source user=replicator",
		Publications:      []string{"pub_orders", "pub_customers"},
		Enabled:           true,
		SlotName:          "test_sub_slot",
		SynchronousCommit: "off",
	}

	// Create
	creatSQL := migrate.CreateSubscriptionSQL(sub)
	require.NotEmpty(t, createSQL)
	assert.Contains(t, createSQL, "CREATE SUBSCRIPTION test_sub")
	assert.Contains(t, createSQL, "CONNECTION 'host=localhost port=5432 dbname=source user=replicator'")
	assert.Contains(t, createSQL, "PUBLICATION pub_orders, pub_customers")

	// Alter: disable
	modified := sub
	modified.Enabled = false
	alterSQL := migrate.AlterSubscriptionSQL(sub, modified)
	require.NotEmpty(t, alterSQL)
	assert.Contains(t, alterSQL, "ALTER SUBSCRIPTION test_sub DISABLE;")

	// Alter: change publications
	modified2 := sub
	modified2.Publications = []string{"pub_all"}
	alterSQL2 := migrate.AlterSubscriptionSQL(sub, modified2)
	require.NotEmpty(t, alterSQL2)
	assert.Contains(t, alterSQL2, "SET PUBLICATION pub_all")

	// Drop
	dropSQL := migrate.DropSubscriptionSQL(sub)
	require.NotEmpty(t, dropSQL)
	assert.Equal(t, "DROP SUBSCRIPTION test_sub;", dropSQL)

	// Full diff cycle
	diff := schema.SubscriptionDiff{
		Added:   []schema.Subscription{sub},
		Removed: []schema.Subscription{},
		Changed: []schema.SubscriptionChange{},
	}
	diffSQL := migrate.SubscriptionDiffSQL(diff)
	assert.Contains(t, diffSQL, "CREATE SUBSCRIPTION test_sub")
}
