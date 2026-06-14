package digit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinkSetCreate(t *testing.T) {
	set := NewLinkSet(10)
	require.Equal(t, 0, len(set))
}

func TestLinkSetAppend(t *testing.T) {

	set := NewLinkSet(0)

	// Append accepts any number of links (including zero) without inspecting them.
	set.Append()
	require.Equal(t, 0, len(set))

	first := NewLink("friend", "text/html", "http://example.com/friend")
	second := NewLink("parent", "application/json", "http://example.com/parent")
	set.Append(first, second)

	require.Equal(t, 2, len(set))
	require.Equal(t, first, set[0])
	require.Equal(t, second, set[1])

	// Append does not de-duplicate; identical links are stored twice.
	set.Append(first)
	require.Equal(t, 3, len(set))
	require.Equal(t, first, set[2])
}

func TestLinkSetFind(t *testing.T) {

	friend := NewLink("friend", "text/html", "http://example.com/friend")
	parent := NewLink("parent", "application/json", "http://example.com/parent")

	set := NewLinkSet(0)
	set.Append(friend, parent)

	// Find matches on "rel" AND "type" (the href is ignored).
	found := set.Find(NewLink("friend", "text/html", "http://other.example.com"))
	require.Equal(t, friend, found)

	// A link with no match returns an empty link.
	missing := set.Find(NewLink("sibling", "text/plain", ""))
	require.True(t, missing.IsEmpty())

	// A matching "rel" but mismatched "type" does not count as a match.
	mismatch := set.Find(NewLink("friend", "application/json", ""))
	require.True(t, mismatch.IsEmpty())
}

func TestLinkSetApply(t *testing.T) {

	set := NewLinkSet(0)

	// Apply inserts a new link when no match exists.
	friend := NewLink("friend", "text/html", "http://example.com/friend")
	set.Apply(friend)
	require.Equal(t, 1, len(set))
	require.Equal(t, friend, set[0])

	// Apply updates the existing link (same "rel" and "type") in place.
	updated := NewLink("friend", "text/html", "http://example.com/new-friend")
	set.Apply(updated)
	require.Equal(t, 1, len(set))
	require.Equal(t, "http://example.com/new-friend", set[0].Href)

	// A different "rel"/"type" is appended as a new entry.
	parent := NewLink("parent", "application/json", "http://example.com/parent")
	set.Apply(parent)
	require.Equal(t, 2, len(set))
	require.Equal(t, parent, set[1])
}

func TestLinkSet(t *testing.T) {

	set := NewLinkSet(0)

	{
		link := NewLink("friend", "text/html", "http://example.com/friend")
		set.ApplyBy("rel", link)

		require.Equal(t, 1, len(set))
		require.Equal(t, link, set.FindBy("rel", "friend"))
		require.Equal(t, link, set.FindBy("type", "text/html"))
		require.Equal(t, "http://example.com/friend", set.FindBy("rel", "friend").GetString("href"))
	}
	{
		link := NewLink("parent", "application/json", "http://example.com/parent")
		set.ApplyBy("rel", link)

		require.Equal(t, 2, len(set))
		require.Equal(t, link, set.FindBy("rel", "parent"))
		require.Equal(t, link, set.FindBy("type", "application/json"))
	}
	{
		link := NewLink("sibling", "text/markdown", "http://example.com/sibling")
		set.ApplyBy("rel", link)

		require.Equal(t, 3, len(set))
		require.Equal(t, link, set.FindBy("rel", "sibling"))
		require.Equal(t, link, set.FindBy("type", "text/markdown"))
	}
	{
		link := NewLink("friend", "text/html", "http://example.com/friend-but-a-different-one")
		set.ApplyBy("rel", link)

		require.Equal(t, 3, len(set))
		require.Equal(t, link, set.FindBy("rel", "friend"))
		require.Equal(t, link, set.FindBy("type", "text/html"))
		require.Equal(t, "http://example.com/friend-but-a-different-one", set.FindBy("rel", "friend").GetString("href"))
	}

	{
		link := set.FindBy("rel", "nobody")
		require.True(t, link.IsEmpty())
	}

	require.Equal(t, 3, len(set))

	set.Remove(Link{RelationType: "sibling", MediaType: "text/markdown", Href: "http://example.com/sibling"})
	require.Equal(t, 2, len(set))

	set.Remove(Link{RelationType: "parent", MediaType: "application/json", Href: "http://example.com/parent"})
	require.Equal(t, 1, len(set))

	set.RemoveBy("rel", "friend")
	require.Zero(t, len(set))

	set.RemoveBy("rel", "missing")
	require.Zero(t, len(set))
}
