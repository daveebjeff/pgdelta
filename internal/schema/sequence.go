package schema

import "fmt"

// Sequence represents a PostgreSQL sequence object.
type Sequence struct {
	Schema    string
	Name      string
	Start     int64
	Increment int64
	MinValue  int64
	MaxValue  int64
	CacheSize int64
	Cycle     bool
}

// FullName returns the fully qualified sequence name.
func (s Sequence) FullName() string {
	return fmt.Sprintf("%s.%s", s.Schema, s.Name)
}

// SequenceDiff represents changes between two sequence snapshots.
type SequenceDiff struct {
	Added   []Sequence
	Removed []Sequence
	Changed []SequenceChange
}

// SequenceChange holds the before/after state of a modified sequence.
type SequenceChange struct {
	Old Sequence
	New Sequence
}

// DiffSequences computes the diff between two slices of sequences.
func DiffSequences(old, new []Sequence) SequenceDiff {
	diff := SequenceDiff{}

	oldMap := make(map[string]Sequence, len(old))
	for _, s := range old {
		oldMap[s.FullName()] = s
	}

	newMap := make(map[string]Sequence, len(new))
	for _, s := range new {
		newMap[s.FullName()] = s
	}

	for _, s := range new {
		if _, exists := oldMap[s.FullName()]; !exists {
			diff.Added = append(diff.Added, s)
		}
	}

	for _, s := range old {
		if _, exists := newMap[s.FullName()]; !exists {
			diff.Removed = append(diff.Removed, s)
		}
	}

	for _, newSeq := range new {
		if oldSeq, exists := oldMap[newSeq.FullName()]; exists {
			if !sequencesEqual(oldSeq, newSeq) {
				diff.Changed = append(diff.Changed, SequenceChange{Old: oldSeq, New: newSeq})
			}
		}
	}

	return diff
}

func sequencesEqual(a, b Sequence) bool {
	return a.Start == b.Start &&
		a.Increment == b.Increment &&
		a.MinValue == b.MinValue &&
		a.MaxValue == b.MaxValue &&
		a.CacheSize == b.CacheSize &&
		a.Cycle == b.Cycle
}
