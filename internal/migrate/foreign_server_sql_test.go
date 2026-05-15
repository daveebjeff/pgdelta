package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

var baseServer = schema.ForeignServer{
	Name:    "my_server",
	FDWName: "postgres_fdw",
	Version: "1.0",
	Options: map[string]string{"host": "localhost", "dbname": "mydb"},
	Owner:   "admin",
}

func TestCreateForeignServerSQL_Basic(t *testing.T) {
	sql := CreateForeignServerSQL(baseServer)
	assert.Contains(t, sql, "CREATE SERVER my_server")
	assert.Contains(t, sql, "VERSION '1.0'")
	assert.Contains(t, sql, "FOREIGN DATA WRAPPER postgres_fdw")
	assert.Contains(t, sql, "OPTIONS (")
	assert.Contains(t, sql, "OWNER TO admin")
}

func TestCreateForeignServerSQL_NoOptions(t *testing.T) {
	fs := baseServer
	fs.Options = nil
	fs.Version = ""
	sql := CreateForeignServerSQL(fs)
	assert.Contains(t, sql, "CREATE SERVER my_server")
	assert.NotContains(t, sql, "OPTIONS")
	assert.NotContains(t, sql, "VERSION")
}

func TestDropForeignServerSQL(t *testing.T) {
	sql := DropForeignServerSQL(baseServer)
	assert.Equal(t, "DROP SERVER IF EXISTS my_server CASCADE;", sql)
}

func TestAlterForeignServerSQL_Version(t *testing.T) {
	fs := baseServer
	fs.Version = "2.0"
	sql := AlterForeignServerSQL(fs)
	assert.Contains(t, sql, "ALTER SERVER my_server VERSION '2.0'")
}

func TestForeignServerDiffSQL_AddedAndRemoved(t *testing.T) {
	added := []schema.ForeignServer{baseServer}
	removed := []schema.ForeignServer{{Name: "old_server", FDWName: "postgres_fdw"}}
	stmts := ForeignServerDiffSQL(added, removed, nil)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "DROP SERVER")
	assert.Contains(t, stmts[1], "CREATE SERVER")
}

func TestForeignServerDiffSQL_Changed(t *testing.T) {
	changed := []schema.ForeignServer{baseServer}
	stmts := ForeignServerDiffSQL(nil, nil, changed)
	assert.Len(t, stmts, 1)
	assert.Contains(t, stmts[0], "ALTER SERVER")
}
