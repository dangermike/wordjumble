package maptrie

import helper "github.com/dangermike/wordjumble/trie"

type trie map[byte]trie

var hit = trie{}

func load(mt trie, word []byte) trie {
	if mt == nil {
		mt = trie{}
	}
	if len(word) == 0 {
		mt[0] = hit
	} else {
		mt[word[0]] = load(mt[word[0]], word[1:])
	}

	return mt
}

func contains(mt trie, word []byte) bool {
	if mt == nil {
		return false
	}

	if len(word) == 0 {
		_, ok := mt[0]
		return ok
	}

	return contains(mt[word[0]], word[1:])
}

// PermuteAll returns all permutations of the provided letters, allowing
// duplicates, that appear in the set of loaded words.
func permuteAll(mt trie, letters []byte, consume bool) [][]byte {
	retval := permuteInner(mt, letters, consume)
	for i := 0; i < len(retval); i++ {
		// reverse by starting at each end and swapping as we go
		for a, b := 0, len(retval[i])-1; a < b; a, b = a+1, b-1 {
			retval[i][a], retval[i][b] = retval[i][b], retval[i][a]
		}
	}
	return helper.UniqueifyResults(retval)
}

// permuteInner builds up in reverse, which lets us avoid copying
func permuteInner(mt trie, letters []byte, consume bool) [][]byte {
	if mt == nil {
		return nil
	}

	retval := [][]byte{}
	if mt[0] != nil {
		retval = append(retval, []byte{})
	}

	for i := 0; i < len(letters); i++ {
		letter := letters[i]
		if mt[letter] != nil {
			next := letters
			if consume {
				letters[0], letters[i] = letters[i], letters[0]
				next = next[1:]
			}
			for _, child := range permuteInner(mt[letter], next, consume) {
				retval = append(retval, append(child, letter))
			}
		}
	}

	return retval
}
