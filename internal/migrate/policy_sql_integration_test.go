package migrate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/pgdelta/internal/migrate"
	"github.com/your-org/pgdelta/internal/schema"
)

func TestPolicyMigration_FullCycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	db := setupTestDB(t, ctx)
	defer db.Close()

	_, err := db.ExecContext(ctx, `CREATE TABLE public.orders (id serial PRIMARY KEY, user_id int)`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `ALTER TABLE public.orders ENABLE ROW LEVEL SECURITY`)
	require.NoError(t, err)

	p := schema.Policy{
		Schema:     "public",
		Table:      "orders",
		Name:       "test_policy",
		Command:    "SELECT",
		Permissive: true,
		Roles:      []string{"PUBLIC"},
		Using:      "(true)",
	}

	createSQL := migrate.CreatePolicySQL(p)
	_, err = db.ExecContext(ctx, createSQL)
	require.NoError(t, err, "failed to create policy: %s", createSQL)

	var policyName string
	err = db.QueryRowContext(ctx,
		`SELECT policyname FROM pg_policies WHERE schemaname = 'public' AND tablename = 'orders' AND policyname = 'test_policy'`,
	).Scan(&policyName)
	require.NoError(t, err)
	assert.Equal(t, "test_policy", policyName)

	updated := p
	updated.Using = "(false)"
	alterSQL := migrate.AlterPolicySQL(updated)
	_, err = db.ExecContext(ctx, alterSQL)
	require.NoError(t, err, "failed to alter policy: %s", alterSQL)

	dropSQL := migrate.DropPolicySQL(p)
	_, err = db.ExecContext(ctx, dropSQL)
	require.NoError(t, err, "failed to drop policy: %s", dropSQL)

	err = db.QueryRowContext(ctx,
		`SELECT policyname FROM pg_policies WHERE schemaname = 'public' AND tablename = 'orders' AND policyname = 'test_policy'`,
	).Scan(&policyName)
	assert.Error(t, err, "policy should no longer exist")
}
