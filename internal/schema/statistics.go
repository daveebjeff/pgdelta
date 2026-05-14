package schema

import "fmt"

// Statistic represents a PostgreSQL extended statistics object.
type Statistic struct {
	Schema     string
	Name       string
	TableName  string
	Columns    []string
	Kinds      []string // e.g. "dependencies", "ndistinct", "mcv"
}

func (s Statistic) FullName() string {
	return fmt.Sprintf("%s.%s", s.Schema, s.Name)
}

// StatisticDiff holds added and removed extended statistics.
type StatisticDiff struct {
	Added   []Statistic
	Removed []Statistic
	Changed []Statistic
}

func (d StatisticDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffStatistics compares two slices of Statistic and returns a StatisticDiff.
func DiffStatistics(old, new []Statistic) StatisticDiff {
	oldMap := make(map[string]Statistic, len(old))
	for _, s := range old {
		oldMap[s.FullName()] = s
	}

	newMap := make(map[string]Statistic, len(new))
	for _, s := range new {
		newMap[s.FullName()] = s
	}

	var diff StatisticDiff

	for _, s := range new {
		if o, exists := oldMap[s.FullName()]; !exists {
			diff.Added = append(diff.Added, s)
		} else if !statisticsEqual(o, s) {
			diff.Changed = append(diff.Changed, s)
		}
	}

	for _, s := range old {
		if _, exists := newMap[s.FullName()]; !exists {
			diff.Removed = append(diff.Removed, s)
		}
	}

	return diff
}

func statisticsEqual(a, b Statistic) bool {
	if a.TableName != b.TableName {
		return false
	}
	if len(a.Columns) != len(b.Columns) {
		return false
	}
	for i := range a.Columns {
		if a.Columns[i] != b.Columns[i] {
			return false
		}
	}
	if len(a.Kinds) != len(b.Kinds) {
		return false
	}
	for i := range a.Kinds {
		if a.Kinds[i] != b.Kinds[i] {
			return false
		}
	}
	return true
}
