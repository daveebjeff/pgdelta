package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseForeignServer = ForeignServer{
	Name:    "my_server",
	FDWName: "postgres_fdw",
	Type:    "",
	Version: "1.0",
	Options: map[string]string{"host": "localhost", "port": "5432"},
	Owner:   "admin",
}

func TestForeignServerFullName(t *testing.T) {
	fs := baseForeignServer
	assert.Equal(t, "my_server", fs.FullName())
}

func TestDiffForeignServers_NoChanges(t *testing.T) {
	fs := baseForeignServer
	added, removed, changed := DiffForeignServers([]ForeignServer{fs}, []ForeignServer{fs})
	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Empty(t, changed)
}

func TestDiffForeignServers_AddedServer(t *testing.T) {
	fs := baseForeignServer
	added, removed, changed := DiffForeignServers(nil, []ForeignServer{fs})
	assert.Equal(t, []ForeignServer{fs}, added)
	assert.Empty(t, removed)
	assert.Empty(t, changed)
}

func TestDiffForeignServers_RemovedServer(t *testing.T) {
	fs := baseForeignServer
	added, removed, changed := DiffForeignServers([]ForeignServer{fs}, nil)
	assert.Empty(t, added)
	assert.Equal(t, []ForeignServer{fs}, removed)
	assert.Empty(t, changed)
}

func TestDiffForeignServers_ChangedServer(t *testing.T) {
	old := baseForeignServer
	new := baseForeignServer
	new.Version = "2.0"
	added, removed, changed := DiffForeignServers([]ForeignServer{old}, []ForeignServer{new})
	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Equal(t, []ForeignServer{new}, changed)
}

func TestDiffForeignServers_ChangedOptions(t *testing.T) {
	old := baseForeignServer
	new := baseForeignServer
	new.Options = map[string]string{"host": "remotehost", "port": "5432"}
	added, removed, changed := DiffForeignServers([]ForeignServer{old}, []ForeignServer{new})
	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Equal(t, []ForeignServer{new}, changed)
}
