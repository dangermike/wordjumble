package arraytrie

import helper "github.com/dangermike/wordjumble/trie"

type Trie struct {
	trie  *trie
	count int
}

func New() helper.Trie {
	return &Trie{trie: &trie{}}
}

func (a *Trie) Load(word []byte) bool {
	a.trie = load(a.trie, word)
	a.count++
	return true
}

func (a *Trie) LoadString(word string) bool {
	return a.Load([]byte(word))
}

func (a *Trie) PermuteAll(letters []byte, consume bool) [][]byte {
	return permuteAll(a.trie, []byte(letters), consume)
}

func (a *Trie) Count() int {
	return a.count
}
