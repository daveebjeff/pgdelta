package schema

import "fmt"

// Function represents a PostgreSQL function/procedure.
type Function struct {
	Schema     string
	Name       string
	Arguments  string // e.g. "a integer, b integer"
	ReturnType string
	Language   string
	Body       string
	Volatility string // VOLATILE, STABLE, IMMUTABLE
}

// FullName returns the schema-qualified function name with arguments.
func (f Function) FullName() string {
	return fmt.Sprintf("%s.%s(%s)", f.Schema, f.Name, f.Arguments)
}

// FunctionDiff holds added, removed, and changed functions.
type FunctionDiff struct {
	Added   []Function
	Removed []Function
	Changed []Function
}

// DiffFunctions compares two slices of functions and returns a FunctionDiff.
func DiffFunctions(old, new []Function) FunctionDiff {
	diff := FunctionDiff{}

	oldMap := make(map[string]Function, len(old))
	for _, f := range old {
		oldMap[f.FullName()] = f
	}

	newMap := make(map[string]Function, len(new))
	for _, f := range new {
		newMap[f.FullName()] = f
	}

	for _, f := range new {
		if _, exists := oldMap[f.FullName()]; !exists {
			diff.Added = append(diff.Added, f)
		}
	}

	for _, f := range old {
		if _, exists := newMap[f.FullName()]; !exists {
			diff.Removed = append(diff.Removed, f)
		}
	}

	for _, f := range new {
		if oldF, exists := oldMap[f.FullName()]; exists {
			if !functionsEqual(oldF, f) {
				diff.Changed = append(diff.Changed, f)
			}
		}
	}

	return diff
}

func functionsEqual(a, b Function) bool {
	return a.Body == b.Body &&
		a.ReturnType == b.ReturnType &&
		a.Language == b.Language &&
		a.Volatility == b.Volatility
}
