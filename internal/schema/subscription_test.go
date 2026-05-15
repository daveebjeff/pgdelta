package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseSubscription = Subscription{
	Name:        "mysub",
	ConnInfo:    "host=localhost dbname=mydb",
	Publications: []string{"mypub"},
	Enabled:     true,
	SlotName:    "mysub_slot",
}

func TestSubscriptionFullName(t *testing.T) {
	s := baseSubscription
	assert.Equal(t, "mysub", s.FullName())
}

func TestDiffSubscriptions_NoChanges(t *testing.T) {
	old := []Subscription{baseSubscription}
	new := []Subscription{baseSubscription}
	diff := DiffSubscriptions(old, new)
	assert.True(t, diff.IsEmpty())
}

func TestDiffSubscriptions_AddedSubscription(t *testing.T) {
	newSub := Subscription{Name: "newsub", ConnInfo: "host=remote dbname=db", Publications: []string{"pub1"}, Enabled: true, SlotName: "newsub_slot"}
	diff := DiffSubscriptions(nil, []Subscription{newSub})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, "newsub", diff.Added[0].Name)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffSubscriptions_RemovedSubscription(t *testing.T) {
	diff := DiffSubscriptions([]Subscription{baseSubscription}, nil)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, "mysub", diff.Removed[0].Name)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffSubscriptions_ChangedSubscription(t *testing.T) {
	modified := baseSubscription
	modified.Enabled = false
	diff := DiffSubscriptions([]Subscription{baseSubscription}, []Subscription{modified})
	assert.Len(t, diff.Changed, 1)
	assert.False(t, diff.Changed[0].Enabled)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}

func TestDiffSubscriptions_ChangedPublications(t *testing.T) {
	modified := baseSubscription
	modified.Publications = []string{"pub1", "pub2"}
	diff := DiffSubscriptions([]Subscription{baseSubscription}, []Subscription{modified})
	assert.Len(t, diff.Changed, 1)
}
