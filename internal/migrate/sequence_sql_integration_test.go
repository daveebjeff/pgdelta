package migrate_test

import (
	"context"
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSequenceMigration_FullCycle(t *testing.T) {
	db := connectTestDB(t)
	defer db.Close()

	ctx := context.Background()

	_, err := db.ExecContext(ctx, "DROP SEQUENCE IF EXISTS public.test_seq;")
	require.NoError(t, err)

	seq := schema.Sequence{
		Schema:    "public",
		Name:      "test_seq",
		Start:     1,
		Increment: 1,
		MinValue:  1,
		MaxValue:  9223372036854775807,
		CacheSize: 1,
		Cycle:     false,
	}

	createSQL := migrate.CreateSequenceSQL(seq)
	_, err = db.ExecContext(ctx, createSQL)
	require.NoError(t, err, "CREATE SEQUENCE should succeed")

	var exists bool
	err = db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_sequences WHERE schemaname=$1 AND sequencename=$2)",
		"public", "test_seq",
	).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "sequence should exist after creation")

	alteredSeq := seq
	alteredSeq.Increment = 5
	alteredSeq.CacheSize = 10
	alterSQL := migrate.AlterSequenceSQL(seq, alteredSeq)
	_, err = db.ExecContext(ctx, alterSQL)
	require.NoError(t, err, "ALTER SEQUENCE should succeed")

	dropSQL := migrate.DropSequenceSQL(seq)
	_, err = db.ExecContext(ctx, dropSQL)
	require.NoError(t, err, "DROP SEQUENCE should succeed")

	err = db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_sequences WHERE schemaname=$1 AND sequencename=$2)",
		"public", "test_seq",
	).Scan(&exists)
	require.NoError(t, err)
	assert.False(t, exists, "sequence should not exist after drop")
}
