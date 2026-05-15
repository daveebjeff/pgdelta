package schema

import "fmt"

// Subscription represents a PostgreSQL logical replication subscription.
type Subscription struct {
	Name        string
	ConnInfo    string
	Publications []string
	Enabled     bool
	SlotName    string
}

// FullName returns the unique identifier for the subscription.
func (s Subscription) FullName() string {
	return s.Name
}

// SubscriptionDiff holds added, removed, and changed subscriptions.
type SubscriptionDiff struct {
	Added   []Subscription
	Removed []Subscription
	Changed []Subscription
}

// IsEmpty returns true if there are no subscription changes.
func (d SubscriptionDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffSubscriptions computes the diff between two sets of subscriptions.
func DiffSubscriptions(old, new []Subscription) SubscriptionDiff {
	oldMap := make(map[string]Subscription, len(old))
	for _, s := range old {
		oldMap[s.FullName()] = s
	}
	newMap := make(map[string]Subscription, len(new))
	for _, s := range new {
		newMap[s.FullName()] = s
	}

	var diff SubscriptionDiff
	for _, s := range new {
		if o, ok := oldMap[s.FullName()]; !ok {
			diff.Added = append(diff.Added, s)
		} else if !subscriptionsEqual(o, s) {
			diff.Changed = append(diff.Changed, s)
		}
	}
	for _, s := range old {
		if _, ok := newMap[s.FullName()]; !ok {
			diff.Removed = append(diff.Removed, s)
		}
	}
	return diff
}

func subscriptionsEqual(a, b Subscription) bool {
	if a.ConnInfo != b.ConnInfo || a.Enabled != b.Enabled || a.SlotName != b.SlotName {
		return false
	}
	if len(a.Publications) != len(b.Publications) {
		return false
	}
	for i := range a.Publications {
		if a.Publications[i] != b.Publications[i] {
			return false
		}
	}
	return true
}

// ensure fmt is used
var _ = fmt.Sprintf
