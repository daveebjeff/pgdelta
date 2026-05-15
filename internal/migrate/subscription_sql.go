package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateSubscriptionSQL generates a CREATE SUBSCRIPTION statement.
func CreateSubscriptionSQL(s schema.Subscription) string {
	pubs := strings.Join(s.Publications, ", ")
	sql := fmt.Sprintf(
		"CREATE SUBSCRIPTION %s CONNECTION '%s' PUBLICATION %s WITH (slot_name = '%s', enabled = %v);",
		s.Name, s.ConnInfo, pubs, s.SlotName, s.Enabled,
	)
	return sql
}

// DropSubscriptionSQL generates a DROP SUBSCRIPTION statement.
func DropSubscriptionSQL(s schema.Subscription) string {
	return fmt.Sprintf("DROP SUBSCRIPTION %s;", s.Name)
}

// AlterSubscriptionSQL generates ALTER SUBSCRIPTION statements for a changed subscription.
func AlterSubscriptionSQL(s schema.Subscription) []string {
	var stmts []string
	pubs := strings.Join(s.Publications, ", ")
	stmts = append(stmts, fmt.Sprintf(
		"ALTER SUBSCRIPTION %s CONNECTION '%s';",
		s.Name, s.ConnInfo,
	))
	stmts = append(stmts, fmt.Sprintf(
		"ALTER SUBSCRIPTION %s SET PUBLICATION %s;",
		s.Name, pubs,
	))
	if s.Enabled {
		stmts = append(stmts, fmt.Sprintf("ALTER SUBSCRIPTION %s ENABLE;", s.Name))
	} else {
		stmts = append(stmts, fmt.Sprintf("ALTER SUBSCRIPTION %s DISABLE;", s.Name))
	}
	return stmts
}

// SubscriptionDiffSQL generates SQL statements for a SubscriptionDiff.
func SubscriptionDiffSQL(diff schema.SubscriptionDiff) []string {
	var stmts []string
	for _, s := range diff.Removed {
		stmts = append(stmts, DropSubscriptionSQL(s))
	}
	for _, s := range diff.Added {
		stmts = append(stmts, CreateSubscriptionSQL(s))
	}
	for _, s := range diff.Changed {
		stmts = append(stmts, AlterSubscriptionSQL(s)...)
	}
	return stmts
}
