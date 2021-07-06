package maptrie_test

import (
	"sort"
	"testing"

	"github.com/dangermike/wordjumble/maptrie"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	mt := maptrie.New()
	require.False(t, maptrie.ContainsString(mt, ""))
	require.False(t, maptrie.ContainsString(mt, "x"))
	mt = nil
	require.False(t, maptrie.ContainsString(mt, ""))
	require.False(t, maptrie.ContainsString(mt, "x"))
}

func TestFull(t *testing.T) {
	mt := maptrie.LoadString(maptrie.LoadString(nil, "full"), "fu")
	require.False(t, maptrie.ContainsString(mt, ""))
	require.False(t, maptrie.ContainsString(mt, "x"))
	require.False(t, maptrie.ContainsString(mt, "fullx"))
	require.True(t, maptrie.ContainsString(mt, "full"))
	require.True(t, maptrie.ContainsString(mt, "fu"))
}

func TestPermuteAll(t *testing.T) {
	mt := maptrie.LoadString(maptrie.LoadString(nil, "full"), "fu")
	rets := []string{}
	for _, ret := range maptrie.PermuteAll(mt, []byte("fulxd")) {
		rets = append(rets, string(ret))
	}
	sort.Strings(rets)
	require.Equal(t, []string{"fu", "full"}, rets)
}
