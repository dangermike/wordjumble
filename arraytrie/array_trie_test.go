package arraytrie_test

import (
	"sort"
	"testing"

	"github.com/dangermike/wordjumble/arraytrie"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	at := arraytrie.New()
	require.False(t, arraytrie.ContainsString(at, ""))
	require.False(t, arraytrie.ContainsString(at, "x"))
	at = nil
	require.False(t, arraytrie.ContainsString(at, ""))
	require.False(t, arraytrie.ContainsString(at, "x"))
}

func TestFull(t *testing.T) {
	at := arraytrie.LoadString(arraytrie.LoadString(nil, "full"), "fu")
	require.False(t, arraytrie.ContainsString(at, ""))
	require.False(t, arraytrie.ContainsString(at, "x"))
	require.False(t, arraytrie.ContainsString(at, "fullx"))
	require.True(t, arraytrie.ContainsString(at, "full"))
	require.True(t, arraytrie.ContainsString(at, "fu"))
}

func TestPermuteAll(t *testing.T) {
	at := arraytrie.LoadString(arraytrie.LoadString(nil, "full"), "fu")
	rets := []string{}
	for _, ret := range arraytrie.PermuteAll(at, []byte("fulxd")) {
		rets = append(rets, string(ret))
	}
	sort.Strings(rets)
	require.Equal(t, []string{"fu", "full"}, rets)
}
