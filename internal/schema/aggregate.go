package schema

import "fmt"

// Aggregate represents a PostgreSQL aggregate function.
type Aggregate struct {
	Schema      string
	Name        string
	ArgTypes    []string
	SFuncName   string
	SType       string
	InitCond    *string
	FinalFunc   *string
}

func (a Aggregate) FullName() string {
	return fmt.Sprintf("%s.%s(%s)", a.Schema, a.Name, joinArgs(a.ArgTypes))
}

func joinArgs(args []string) string {
	result := ""
	for i, a := range args {
		if i > 0 {
			result += ", "
		}
		result += a
	}
	return result
}

// AggregateDiff holds added, removed, and changed aggregates.
type AggregateDiff struct {
	Added   []Aggregate
	Removed []Aggregate
	Changed []Aggregate
}

func (d AggregateDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffAggregates computes the diff between two slices of aggregates.
func DiffAggregates(old, new []Aggregate) AggregateDiff {
	diff := AggregateDiff{}

	oldMap := make(map[string]Aggregate, len(old))
	for _, a := range old {
		oldMap[a.FullName()] = a
	}

	newMap := make(map[string]Aggregate, len(new))
	for _, a := range new {
		newMap[a.FullName()] = a
	}

	for _, a := range new {
		if existing, ok := oldMap[a.FullName()]; !ok {
			diff.Added = append(diff.Added, a)
		} else if !aggregatesEqual(existing, a) {
			diff.Changed = append(diff.Changed, a)
		}
	}

	for _, a := range old {
		if _, ok := newMap[a.FullName()]; !ok {
			diff.Removed = append(diff.Removed, a)
		}
	}

	return diff
}

func aggregatesEqual(a, b Aggregate) bool {
	if a.SFuncName != b.SFuncName || a.SType != b.SType {
		return false
	}
	if ptrStr(a.InitCond) != ptrStr(b.InitCond) {
		return false
	}
	if ptrStr(a.FinalFunc) != ptrStr(b.FinalFunc) {
		return false
	}
	return true
}

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
