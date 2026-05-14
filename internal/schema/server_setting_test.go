package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseSetting = ServerSetting{
	Name:  "work_mem",
	Value: "4MB",
}

func TestServerSettingFullName(t *testing.T) {
	assert.Equal(t, "work_mem", baseSetting.FullName())
}

func TestDiffServerSettings_NoChanges(t *testing.T) {
	settings := []ServerSetting{baseSetting}
	diff := DiffServerSettings(settings, settings)
	assert.True(t, diff.IsEmpty())
}

func TestDiffServerSettings_AddedSetting(t *testing.T) {
	newSetting := ServerSetting{Name: "max_connections", Value: "200"}
	diff := DiffServerSettings([]ServerSetting{}, []ServerSetting{newSetting})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, newSetting, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffServerSettings_RemovedSetting(t *testing.T) {
	diff := DiffServerSettings([]ServerSetting{baseSetting}, []ServerSetting{})
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseSetting, diff.Removed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffServerSettings_ChangedSetting(t *testing.T) {
	updated := ServerSetting{Name: "work_mem", Value: "8MB"}
	diff := DiffServerSettings([]ServerSetting{baseSetting}, []ServerSetting{updated})
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, baseSetting, diff.Changed[0].Old)
	assert.Equal(t, updated, diff.Changed[0].New)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}

func TestDiffServerSettings_IsEmpty_NoChanges(t *testing.T) {
	diff := ServerSettingDiff{}
	assert.True(t, diff.IsEmpty())
}
