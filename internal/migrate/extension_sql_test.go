package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

var baseExt = schema.Extension{
	Name:    "uuid-ossp",
	Schema:  "public",
	Version: "1.1",
}

func TestCreateExtensionSQL(t *testing.T) {
	sql := CreateExtensionSQL(baseExt)
	assert.Equal(t, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA "public" VERSION '1.1';`, sql)
}

func TestCreateExtensionSQL_NoVersion(t *testing.T) {
	e := schema.Extension{Name: "pgcrypto", Schema: "public"}
	sql := CreateExtensionSQL(e)
	assert.Equal(t, `CREATE EXTENSION IF NOT EXISTS "pgcrypto" WITH SCHEMA "public";`, sql)
}

func TestDropExtensionSQL(t *testing.T) {
	sql := DropExtensionSQL(baseExt)
	assert.Equal(t, `DROP EXTENSION IF EXISTS "uuid-ossp";`, sql)
}

func TestAlterExtensionSQL_WithVersion(t *testing.T) {
	e := schema.Extension{Name: "uuid-ossp", Schema: "public", Version: "1.2"}
	sql := AlterExtensionSQL(e)
	assert.Equal(t, `ALTER EXTENSION "uuid-ossp" UPDATE TO '1.2';`, sql)
}

func TestAlterExtensionSQL_NoVersion(t *testing.T) {
	e := schema.Extension{Name: "uuid-ossp", Schema: "public"}
	sql := AlterExtensionSQL(e)
	assert.Equal(t, `ALTER EXTENSION "uuid-ossp" UPDATE;`, sql)
}

func TestExtensionDiffSQL_AddedAndRemoved(t *testing.T) {
	diff := schema.ExtensionDiff{
		Added:   []schema.Extension{{Name: "pgcrypto", Schema: "public", Version: "1.3"}},
		Removed: []schema.Extension{baseExt},
	}
	stmts := ExtensionDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "DROP EXTENSION")
	assert.Contains(t, stmts[1], "CREATE EXTENSION")
}

func TestExtensionDiffSQL_Changed(t *testing.T) {
	diff := schema.ExtensionDiff{
		Changed: []schema.Extension{{Name: "uuid-ossp", Schema: "public", Version: "1.2"}},
	}
	stmts := ExtensionDiffSQL(diff)
	assert.Len(t, stmts, 1)
	assert.Equal(t, `ALTER EXTENSION "uuid-ossp" UPDATE TO '1.2';`, stmts[0])
}
