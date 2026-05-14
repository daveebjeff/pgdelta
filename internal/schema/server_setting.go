package schema

import "fmt"

// ServerSetting represents a PostgreSQL server-level configuration setting
// (e.g., set via ALTER SYSTEM or pg_settings).
type ServerSetting struct {
	Name  string
	Value string
}

// FullName returns the unique identifier for the setting.
func (s ServerSetting) FullName() string {
	return s.Name
}

// ServerSettingDiff captures added, removed, and changed server settings.
type ServerSettingDiff struct {
	Added   []ServerSetting
	Removed []ServerSetting
	Changed []ServerSettingChange
}

// ServerSettingChange records a setting whose value has changed.
type ServerSettingChange struct {
	Old ServerSetting
	New ServerSetting
}

// IsEmpty returns true when there are no differences.
func (d ServerSettingDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffServerSettings computes the diff between two slices of ServerSetting.
func DiffServerSettings(old, new []ServerSetting) ServerSettingDiff {
	oldMap := make(map[string]ServerSetting, len(old))
	for _, s := range old {
		oldMap[s.Name] = s
	}

	newMap := make(map[string]ServerSetting, len(new))
	for _, s := range new {
		newMap[s.Name] = s
	}

	var diff ServerSettingDiff

	for _, s := range new {
		if o, exists := oldMap[s.Name]; !exists {
			diff.Added = append(diff.Added, s)
		} else if !serverSettingsEqual(o, s) {
			diff.Changed = append(diff.Changed, ServerSettingChange{Old: o, New: s})
		}
	}

	for _, s := range old {
		if _, exists := newMap[s.Name]; !exists {
			diff.Removed = append(diff.Removed, s)
		}
	}

	return diff
}

func serverSettingsEqual(a, b ServerSetting) bool {
	return a.Value == b.Value
}

var _ = fmt.Sprintf // ensure fmt is used if needed in future
