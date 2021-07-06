package arraytrie

type Trie struct {
	trie  *trie
	count int
}

func (a *Trie) Load(word []byte) bool {
	a.trie = load(a.trie, word)
	a.count++
	return true
}

func (a *Trie) LoadString(word string) bool {
	return a.Load([]byte(word))
}

func (a *Trie) PermuteAll(letters []byte) [][]byte {
	return permuteAll(a.trie, []byte(letters))
}

func (a *Trie) Count() int {
	return a.count
}

func New() *Trie {
	return &Trie{trie: &trie{}}
}
