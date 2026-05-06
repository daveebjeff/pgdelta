package migrate

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctionMigration_FullCycle(t *testing.T) {
	dsn := os.Getenv("PGDELTA_TEST_DSN")
	if dsn == "" {
		t.Skip("PGDELTA_TEST_DSN not set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	defer conn.Close(ctx)

	f := schema.Function{
		Schema:     "public",
		Name:       "pgdelta_test_add",
		Arguments:  "a integer, b integer",
		ReturnType: "integer",
		Language:   "plpgsql",
		Body:       "BEGIN RETURN a + b; END;",
		Volatility: "IMMUTABLE",
	}

	// Create the function.
	createSQL := CreateFunctionSQL(f)
	_, err = conn.Exec(ctx, createSQL)
	require.NoError(t, err)

	// Verify it exists and works.
	var result int
	err = conn.QueryRow(ctx, "SELECT public.pgdelta_test_add(3, 4)").Scan(&result)
	require.NoError(t, err)
	assert.Equal(t, 7, result)

	// Update the function body.
	f.Body = "BEGIN RETURN a + b + 10; END;"
	updateSQL := CreateFunctionSQL(f)
	_, err = conn.Exec(ctx, updateSQL)
	require.NoError(t, err)

	err = conn.QueryRow(ctx, "SELECT public.pgdelta_test_add(3, 4)").Scan(&result)
	require.NoError(t, err)
	assert.Equal(t, 17, result)

	// Drop the function.
	dropSQL := DropFunctionSQL(f)
	_, err = conn.Exec(ctx, dropSQL)
	require.NoError(t, err)

	// Verify it no longer exists.
	err = conn.QueryRow(ctx, "SELECT public.pgdelta_test_add(1, 2)").Scan(&result)
	assert.Error(t, err)
}
