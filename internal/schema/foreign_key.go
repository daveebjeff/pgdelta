package schema

import "fmt"

type ForeignKey struct {
	Schema          string
	Table           string
	Name            string
	Columns         []string
	RefSchema       string
	RefTable        string
	RefColumns      []string
	OnDelete        string
	OnUpdate        string
	Deferrable      bool
	InitiallyDeferred bool
}

func (fk ForeignKey) FullName() string {
	return fmt.Sprintf("%s.%s.%s", fk.Schema, fk.Table, fk.Name)
}

type ForeignKeyDiff struct {
	Added   []ForeignKey
	Removed []ForeignKey
	Changed []ForeignKey
}

func DiffForeignKeys(old, new []ForeignKey) ForeignKeyDiff {
	diff := ForeignKeyDiff{}

	oldMap := make(map[string]ForeignKey, len(old))
	for _, fk := range old {
		oldMap[fk.FullName()] = fk
	}

	newMap := make(map[string]ForeignKey, len(new))
	for _, fk := range new {
		newMap[fk.FullName()] = fk
	}

	for _, fk := range new {
		if oldFK, exists := oldMap[fk.FullName()]; !exists {
			diff.Added = append(diff.Added, fk)
		} else if !foreignKeysEqual(oldFK, fk) {
			diff.Changed = append(diff.Changed, fk)
		}
	}

	for _, fk := range old {
		if _, exists := newMap[fk.FullName()]; !exists {
			diff.Removed = append(diff.Removed, fk)
		}
	}

	return diff
}

func foreignKeysEqual(a, b ForeignKey) bool {
	if a.RefSchema != b.RefSchema || a.RefTable != b.RefTable {
		return false
	}
	if a.OnDelete != b.OnDelete || a.OnUpdate != b.OnUpdate {
		return false
	}
	if a.Deferrable != b.Deferrable || a.InitiallyDeferred != b.InitiallyDeferred {
		return false
	}
	if len(a.Columns) != len(b.Columns) || len(a.RefColumns) != len(b.RefColumns) {
		return false
	}
	for i := range a.Columns {
		if a.Columns[i] != b.Columns[i] {
			return false
		}
	}
	for i := range a.RefColumns {
		if a.RefColumns[i] != b.RefColumns[i] {
			return false
		}
	}
	return true
}
