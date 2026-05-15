package schema

import "fmt"

// ForeignDataWrapper represents a PostgreSQL foreign data wrapper.
type ForeignDataWrapper struct {
	Name       string
	Handler    string
	Validator  string
	Owner      string
	Options    map[string]string
}

func (f ForeignDataWrapper) FullName() string {
	return fmt.Sprintf("fdw:%s", f.Name)
}

// DiffForeignDataWrappers returns added, removed, and changed FDWs.
func DiffForeignDataWrappers(old, new []ForeignDataWrapper) (added, removed, changed []ForeignDataWrapper) {
	oldMap := make(map[string]ForeignDataWrapper, len(old))
	for _, f := range old {
		oldMap[f.Name] = f
	}
	newMap := make(map[string]ForeignDataWrapper, len(new))
	for _, f := range new {
		newMap[f.Name] = f
	}

	for _, f := range new {
		if o, ok := oldMap[f.Name]; !ok {
			added = append(added, f)
		} else if !fdwsEqual(o, f) {
			changed = append(changed, f)
		}
	}
	for _, f := range old {
		if _, ok := newMap[f.Name]; !ok {
			removed = append(removed, f)
		}
	}
	return
}

func fdwsEqual(a, b ForeignDataWrapper) bool {
	if a.Handler != b.Handler || a.Validator != b.Validator || a.Owner != b.Owner {
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
