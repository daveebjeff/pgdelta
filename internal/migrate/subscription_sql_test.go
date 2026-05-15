package migrate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"pgdelta/internal/migrate"
	"pgdelta/internal/schema"
)

var baseSub = schema.Subscription{
	Name:            "my_sub",
	ConnectionInfo:  "host=localhost dbname=source",
	Publications:    []string{"pub1", "pub2"},
	Enabled:         true,
	SlotName:        "my_sub_slot",
	SynchronousCommit: "off",
}

func TestCreateSubscriptionSQL_Basic(t *testing.T) {
	sql := migrate.CreateSubscriptionSQL(baseSub)
	assert.Contains(t, sql, "CREATE SUBSCRIPTION my_sub")
	assert.Contains(t, sql, "CONNECTION 'host=localhost dbname=source'")
	assert.Contains(t, sql, "PUBLICATION pub1, pub2")
}

func TestCreateSubscriptionSQL_Disabled(t *testing.T) {
	s := baseSub
	s.Enabled = false
	sql := migrate.CreateSubscriptionSQL(s)
	assert.Contains(t, sql, "ENABLED = false")
}

func TestDropSubscriptionSQL(t *testing.T) {
	sql := migrate.DropSubscriptionSQL(baseSub)
	assert.Equal(t, "DROP SUBSCRIPTION my_sub;", sql)
}

func TestAlterSubscriptionSQL_Enable(t *testing.T) {
	old := baseSub
	old.Enabled = false
	new := baseSub
	new.Enabled = true
	sql := migrate.AlterSubscriptionSQL(old, new)
	assert.Contains(t, sql, "ALTER SUBSCRIPTION my_sub ENABLE;")
}

func TestAlterSubscriptionSQL_Publications(t *testing.T) {
	old := baseSub
	new := baseSub
	new.Publications = []string{"pub3"}
	sql := migrate.AlterSubscriptionSQL(old, new)
	assert.Contains(t, sql, "SET PUBLICATION pub3")
}

func TestSubscriptionDiffSQL_AddedAndRemoved(t *testing.T) {
	added := []schema.Subscription{baseSub}
	removed := []schema.Subscription{{Name: "old_sub", ConnectionInfo: "host=old", Publications: []string{"pub_old"}}}
	changed := []schema.SubscriptionChange{}

	result := migrate.SubscriptionDiffSQL(schema.SubscriptionDiff{
		Added:   added,
		Removed: removed,
		Changed: changed,
	})

	assert.Contains(t, result, "CREATE SUBSCRIPTION my_sub")
	assert.Contains(t, result, "DROP SUBSCRIPTION old_sub;")
}

func TestSubscriptionDiffSQL_Changed(t *testing.T) {
	old := baseSub
	new := baseSub
	new.Publications = []string{"pub_new"}

	result := migrate.SubscriptionDiffSQL(schema.SubscriptionDiff{
		Changed: []schema.SubscriptionChange{{Old: old, New: new}},
	})

	assert.Contains(t, result, "SET PUBLICATION pub_new")
}
