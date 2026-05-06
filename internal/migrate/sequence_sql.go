package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateSequenceSQL generates a CREATE SEQUENCE statement.
func CreateSequenceSQL(s schema.Sequence) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "CREATE SEQUENCE %s", s.FullName())
	fmt.Fprintf(&sb, "\n    START WITH %d", s.Start)
	fmt.Fprintf(&sb, "\n    INCREMENT BY %d", s.Increment)
	fmt.Fprintf(&sb, "\n    MINVALUE %d", s.MinValue)
	fmt.Fprintf(&sb, "\n    MAXVALUE %d", s.MaxValue)
	fmt.Fprintf(&sb, "\n    CACHE %d", s.CacheSize)
	if s.Cycle {
		sb.WriteString("\n    CYCLE")
	} else {
		sb.WriteString("\n    NO CYCLE")
	}
	sb.WriteString(";")
	return sb.String()
}

// DropSequenceSQL generates a DROP SEQUENCE statement.
func DropSequenceSQL(s schema.Sequence) string {
	return fmt.Sprintf("DROP SEQUENCE %s;", s.FullName())
}

// SequenceDiffSQL generates ALTER SEQUENCE statements for changed sequences.
func SequenceDiffSQL(diff schema.SequenceDiff) []string {
	var stmts []string

	for _, s := range diff.Added {
		stmts = append(stmts, CreateSequenceSQL(s))
	}

	for _, s := range diff.Removed {
		stmts = append(stmts, DropSequenceSQL(s))
	}

	for _, change := range diff.Changed {
		stmts = append(stmts, AlterSequenceSQL(change.Old, change.New))
	}

	return stmts
}

// AlterSequenceSQL generates an ALTER SEQUENCE statement reflecting changes.
func AlterSequenceSQL(old, new schema.Sequence) string {
	var parts []string

	if old.Increment != new.Increment {
		parts = append(parts, fmt.Sprintf("INCREMENT BY %d", new.Increment))
	}
	if old.MinValue != new.MinValue {
		parts = append(parts, fmt.Sprintf("MINVALUE %d", new.MinValue))
	}
	if old.MaxValue != new.MaxValue {
		parts = append(parts, fmt.Sprintf("MAXVALUE %d", new.MaxValue))
	}
	if old.CacheSize != new.CacheSize {
		parts = append(parts, fmt.Sprintf("CACHE %d", new.CacheSize))
	}
	if old.Cycle != new.Cycle {
		if new.Cycle {
			parts = append(parts, "CYCLE")
		} else {
			parts = append(parts, "NO CYCLE")
		}
	}

	if len(parts) == 0 {
		return ""
	}
	return fmt.Sprintf("ALTER SEQUENCE %s %s;", new.FullName(), strings.Join(parts, " "))
}
