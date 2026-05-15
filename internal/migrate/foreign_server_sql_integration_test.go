package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForeignServerMigration_FullCycle(t *testing.T) {
	db := getTestDB(t)

	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS postgres_fdw;`)
	require.NoError(t, err)

	fs := schema.ForeignServer{
		Name:    "integration_server",
		FDWName: "postgres_fdw",
		Version: "1.0",
		Options: map[string]string{"host": "localhost", "dbname": "testdb"},
	}

	// Create
	creatSQL := migrate.CreateForeignServerSQL(fs)
	_, err = db.Exec(createSQL)
	require.NoError(t, err)

	// Verify exists
	var name string
	err = db.QueryRow(`SELECT srvname FROM pg_foreign_server WHERE srvname = $1`, fs.Name).Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, fs.Name, name)

	// Drop
	dropSQL := migrate.DropForeignServerSQL(fs)
	_, err = db.Exec(dropSQL)
	require.NoError(t, err)

	// Verify removed
	err = db.QueryRow(`SELECT srvname FROM pg_foreign_server WHERE srvname = $1`, fs.Name).Scan(&name)
	assert.Error(t, err)
}
