package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func strPtrDomain(s string) *string { return &s }

var baseDomain = Domain{
	Schema:   "public",
	Name:     "us_postal_code",
	BaseType: "text",
	NotNull:  true,
}

func TestDomainFullName(t *testing.T) {
	assert.Equal(t, "public.us_postal_code", baseDomain.FullName())
}

func TestDiffDomains_NoChanges(t *testing.T) {
	diff := DiffDomains([]Domain{baseDomain}, []Domain{baseDomain})
	assert.True(t, diff.IsEmpty())
}

func TestDiffDomains_AddedDomain(t *testing.T) {
	diff := DiffDomains(nil, []Domain{baseDomain})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, baseDomain, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffDomains_RemovedDomain(t *testing.T) {
	diff := DiffDomains([]Domain{baseDomain}, nil)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseDomain, diff.Removed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffDomains_ChangedDomain(t *testing.T) {
	modified := baseDomain
	modified.NotNull = false

	diff := DiffDomains([]Domain{baseDomain}, []Domain{modified})
	assert.Len(t, diff.Changed, 1)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}

func TestDiffDomains_ChangedDefault(t *testing.T) {
	withDefault := baseDomain
	withDefault.Default = strPtrDomain("'00000'")

	diff := DiffDomains([]Domain{baseDomain}, []Domain{withDefault})
	assert.Len(t, diff.Changed, 1)
}

func TestDiffDomains_ChangedCheckClause(t *testing.T) {
	withCheck := baseDomain
	withCheck.CheckName = strPtrDomain("valid_postal")
	withCheck.CheckClause = strPtrDomain("VALUE ~ '^[0-9]{5}$'")

	diff := DiffDomains([]Domain{baseDomain}, []Domain{withCheck})
	assert.Len(t, diff.Changed, 1)
}
