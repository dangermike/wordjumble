package trie_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/dangermike/wordjumble/arraytrie"
	"github.com/dangermike/wordjumble/maptrie"
	"github.com/dangermike/wordjumble/trie"
	"github.com/stretchr/testify/require"
)

// ensure we implement the interface
var (
	_ trie.Trie = maptrie.New()
	_ trie.Trie = arraytrie.New()
)

func TestEmpty(t *testing.T) {
	for _, test := range []struct {
		name string
		new  func() trie.Trie
	}{
		{"arraytrie", arraytrie.New},
		{"maptrie", maptrie.New},
	} {
		t.Run(test.name, func(t *testing.T) {
			tr := test.new()
			require.False(t, tr.ContainsString(""))
			require.False(t, tr.ContainsString("x"))
		})
	}
}

func TestFull(t *testing.T) {
	for _, test := range []struct {
		name string
		new  func() trie.Trie
	}{
		{"arraytrie", arraytrie.New},
		{"maptrie", maptrie.New},
	} {
		tr := test.new()
		tr.LoadString("full")
		tr.LoadString("fu")
		require.False(t, tr.ContainsString(""))
		require.False(t, tr.ContainsString("x"))
		require.False(t, tr.ContainsString("fullx"))
		require.True(t, tr.ContainsString("full"))
		require.True(t, tr.ContainsString("fu"))
	}
}

func TestPermuteAll(t *testing.T) {
	for _, test := range []struct {
		name string
		new  func() trie.Trie
	}{
		{"arraytrie", arraytrie.New},
		{"maptrie", maptrie.New},
	} {
		tr := test.new()
		tr.LoadString("full")
		tr.LoadString("fu")

		for _, test := range []struct {
			word     string
			consume  bool
			words    []string
			expected []string
		}{
			{"fulxd", false, []string{"fu", "full"}, []string{"fu", "full"}},
			{"fulxd", true, []string{"fu", "full"}, []string{"fu"}},
			{"fulxdll", true, []string{"fu", "full"}, []string{"fu", "full"}},
			{"fulxdll", false, []string{"fu", "full"}, []string{"fu", "full"}},
			{"cab", false, []string{"baa", "cab"}, []string{"baa", "cab"}},
			{"cab", true, []string{"baa", "cab"}, []string{"cab"}},
		} {
			t.Run(fmt.Sprintf("%s,%v", test.word, test.consume), func(t *testing.T) {
				mt := maptrie.New()
				for _, word := range test.words {
					mt.LoadString(word)
				}
				rets := []string{}
				for _, ret := range mt.PermuteAll([]byte(test.word), test.consume) {
					rets = append(rets, string(ret))
				}
				sort.Strings(rets)
				require.Equal(t, test.expected, rets)
			})
		}
	}
}
