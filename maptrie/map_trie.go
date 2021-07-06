package maptrie

import helper "github.com/dangermike/wordjumble/trie"

type Trie struct {
	trie  trie
	count int
}

func New() helper.Trie {
	return &Trie{
		trie:  trie{},
		count: 0,
	}
}

func (m *Trie) Load(word []byte) bool {
	m.trie = load(m.trie, word)
	m.count++
	return true
}

func (m *Trie) LoadString(word string) bool {
	return m.Load([]byte(word))
}

func (m *Trie) PermuteAll(letters []byte, consume bool) [][]byte {
	return permuteAll(m.trie, []byte(letters), consume)
}

func (m *Trie) Count() int {
	return m.count
}

func (m *Trie) Contains(word []byte) bool {
	if m == nil {
		return false
	}
	return contains(m.trie, word)
}

func (m *Trie) ContainsString(word string) bool {
	return m.Contains([]byte(word))
}
