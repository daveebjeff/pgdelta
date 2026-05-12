package schema

import "fmt"

// Publication represents a PostgreSQL logical replication publication.
type Publication struct {
	Name       string
	AllTables  bool
	Insert     bool
	Update     bool
	Delete     bool
	Truncate   bool
	Tables     []string
}

func (p Publication) FullName() string {
	return fmt.Sprintf("publication:%s", p.Name)
}

// PublicationDiff holds added, removed, and changed publications.
type PublicationDiff struct {
	Added   []Publication
	Removed []Publication
	Changed []Publication
}

func (d PublicationDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffPublications computes the diff between two slices of publications.
func DiffPublications(old, new []Publication) PublicationDiff {
	oldMap := make(map[string]Publication, len(old))
	for _, p := range old {
		oldMap[p.Name] = p
	}
	newMap := make(map[string]Publication, len(new))
	for _, p := range new {
		newMap[p.Name] = p
	}

	var diff PublicationDiff
	for _, p := range new {
		if o, exists := oldMap[p.Name]; !exists {
			diff.Added = append(diff.Added, p)
		} else if !publicationsEqual(o, p) {
			diff.Changed = append(diff.Changed, p)
		}
	}
	for _, p := range old {
		if _, exists := newMap[p.Name]; !exists {
			diff.Removed = append(diff.Removed, p)
		}
	}
	return diff
}

func publicationsEqual(a, b Publication) bool {
	if a.AllTables != b.AllTables || a.Insert != b.Insert ||
		a.Update != b.Update || a.Delete != b.Delete ||
		a.Truncate != b.Truncate || len(a.Tables) != len(b.Tables) {
		return false
	}
	aSet := make(map[string]struct{}, len(a.Tables))
	for _, t := range a.Tables {
		aSet[t] = struct{}{}
	}
	for _, t := range b.Tables {
		if _, ok := aSet[t]; !ok {
			return false
		}
	}
	return true
}
