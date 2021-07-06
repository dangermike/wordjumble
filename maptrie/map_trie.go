package maptrie

type Trie map[byte]Trie

var hit = Trie{}

func New() Trie {
	return Trie{}
}

func Load(mt Trie, word []byte) Trie {
	if mt == nil {
		mt = Trie{}
	}
	if len(word) == 0 {
		mt[0] = hit
	} else {
		mt[word[0]] = Load(mt[word[0]], word[1:])
	}

	return mt
}

func LoadString(mt Trie, word string) Trie {
	return Load(mt, []byte(word))
}

func ContainsString(mt Trie, word string) bool {
	return Contains(mt, []byte(word))
}

func Contains(mt Trie, word []byte) bool {
	if mt == nil {
		return false
	}

	if len(word) == 0 {
		_, ok := mt[0]
		return ok
	}

	return Contains(mt[word[0]], word[1:])
}

// PermuteAll returns all permutations of the provided letters, allowing
// duplicates, that appear in the set of loaded words.
func PermuteAll(mt Trie, letters []byte) [][]byte {
	retval := permuteInner(mt, letters)
	for i := 0; i < len(retval); i++ {
		// reverse by starting at each end and swapping as we go
		for a, b := 0, len(retval[i])-1; a < b; a, b = a+1, b-1 {
			retval[i][a], retval[i][b] = retval[i][b], retval[i][a]
		}
	}
	return retval
}

// permuteInner builds up in reverse, which lets us avoid copying
func permuteInner(mt Trie, letters []byte) [][]byte {
	if len(letters) == 0 || mt == nil {
		return nil
	}

	retval := [][]byte{}
	if mt[0] != nil {
		retval = append(retval, []byte{})
	}

	for _, letter := range letters {
		if mt[letter] != nil {
			for _, child := range permuteInner(mt[letter], letters) {
				retval = append(retval, append(child, letter))
			}
		}
	}

	return retval
}
