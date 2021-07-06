package maptrie_test

import (
	"sort"
	"testing"

	"github.com/dangermike/wordjumble/maptrie"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	mt := maptrie.New()
	require.False(t, mt.ContainsString(""))
	require.False(t, mt.ContainsString("x"))
	mt = nil
	require.False(t, mt.ContainsString(""))
	require.False(t, mt.ContainsString("x"))
}

func TestFull(t *testing.T) {
	mt := maptrie.New()
	mt.LoadString("full")
	mt.LoadString("fu")
	require.False(t, mt.ContainsString(""))
	require.False(t, mt.ContainsString("x"))
	require.False(t, mt.ContainsString("fullx"))
	require.True(t, mt.ContainsString("full"))
	require.True(t, mt.ContainsString("fu"))
}

func TestPermuteAll(t *testing.T) {
	mt := maptrie.New()
	mt.LoadString("full")
	mt.LoadString("fu")
	rets := []string{}
	for _, ret := range mt.PermuteAll([]byte("fulxd")) {
		rets = append(rets, string(ret))
	}
	sort.Strings(rets)
	require.Equal(t, []string{"fu", "full"}, rets)
}
