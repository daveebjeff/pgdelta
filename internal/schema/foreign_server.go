package schema

import "fmt"

// ForeignServer represents a PostgreSQL foreign server object.
type ForeignServer struct {
	Name       string
	FDWName    string
	Type       string
	Version    string
	Options    map[string]string
	Owner      string
}

// FullName returns the identifier for the foreign server.
func (fs ForeignServer) FullName() string {
	return fs.Name
}

// DiffForeignServers computes added, removed, and changed foreign servers.
func DiffForeignServers(old, new []ForeignServer) (added, removed, changed []ForeignServer) {
	oldMap := make(map[string]ForeignServer, len(old))
	for _, fs := range old {
		oldMap[fs.Name] = fs
	}
	newMap := make(map[string]ForeignServer, len(new))
	for _, fs := range new {
		newMap[fs.Name] = fs
	}

	for _, fs := range new {
		if o, exists := oldMap[fs.Name]; !exists {
			added = append(added, fs)
		} else if !foreignServersEqual(o, fs) {
			changed = append(changed, fs)
		}
	}
	for _, fs := range old {
		if _, exists := newMap[fs.Name]; !exists {
			removed = append(removed, fs)
		}
	}
	return
}

func foreignServersEqual(a, b ForeignServer) bool {
	if a.FDWName != b.FDWName || a.Type != b.Type || a.Version != b.Version || a.Owner != b.Owner {
		return false
	}
	if len(a.Options) != len(b.Options) {
		return false
	}
	for k, v := range a.Options {
		if b.Options[k] != v {
			return false
		}
	}
	return true
}

// optionsString formats a map of options as a PostgreSQL OPTIONS clause value.
func foreignServerOptionsString(opts map[string]string) string {
	if len(opts) == 0 {
		return ""
	}
	var parts []string
	for k, v := range opts {
		parts = append(parts, fmt.Sprintf("%s '%s'", k, v))
	}
	return joinStrings(parts, ", ")
}

func joinStrings(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}
