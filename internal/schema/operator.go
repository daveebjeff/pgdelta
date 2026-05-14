package schema

import "fmt"

// Operator represents a PostgreSQL operator definition.
type Operator struct {
	Schema      string
	Name        string
	LeftType    string
	RightType   string
	ResultType  string
	Procedure   string
	Commutator  string
	Negator     string
}

func (o Operator) FullName() string {
	return fmt.Sprintf("%s.%s(%s,%s)", o.Schema, o.Name, o.LeftType, o.RightType)
}

// OperatorDiff holds added, removed, and changed operators.
type OperatorDiff struct {
	Added   []Operator
	Removed []Operator
	Changed []Operator
}

func (d OperatorDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffOperators computes the diff between two slices of operators.
func DiffOperators(old, new []Operator) OperatorDiff {
	diff := OperatorDiff{}

	oldMap := make(map[string]Operator, len(old))
	for _, o := range old {
		oldMap[o.FullName()] = o
	}

	newMap := make(map[string]Operator, len(new))
	for _, o := range new {
		newMap[o.FullName()] = o
	}

	for key, o := range newMap {
		if _, exists := oldMap[key]; !exists {
			diff.Added = append(diff.Added, o)
		} else if !operatorsEqual(oldMap[key], o) {
			diff.Changed = append(diff.Changed, o)
		}
	}

	for key, o := range oldMap {
		if _, exists := newMap[key]; !exists {
			diff.Removed = append(diff.Removed, o)
		}
	}

	return diff
}

func operatorsEqual(a, b Operator) bool {
	return a.ResultType == b.ResultType &&
		a.Procedure == b.Procedure &&
		a.Commutator == b.Commutator &&
		a.Negator == b.Negator
}
