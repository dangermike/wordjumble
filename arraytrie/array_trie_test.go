package arraytrie_test

import (
	"sort"
	"testing"

	"github.com/dangermike/wordjumble/arraytrie"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	at := arraytrie.New()
	require.False(t, at.ContainsString(""))
	require.False(t, at.ContainsString("x"))
	at = nil
	require.False(t, at.ContainsString(""))
	require.False(t, at.ContainsString("x"))
}

func TestFull(t *testing.T) {
	at := arraytrie.New()
	at.LoadString("full")
	at.LoadString("fu")
	require.False(t, at.ContainsString(""))
	require.False(t, at.ContainsString("x"))
	require.False(t, at.ContainsString("fullx"))
	require.True(t, at.ContainsString("full"))
	require.True(t, at.ContainsString("fu"))
}

func TestPermuteAll(t *testing.T) {
	at := arraytrie.New()
	at.LoadString("full")
	at.LoadString("fu")
	rets := []string{}
	for _, ret := range at.PermuteAll([]byte("fulxd")) {
		rets = append(rets, string(ret))
	}
	sort.Strings(rets)
	require.Equal(t, []string{"fu", "full"}, rets)
}
