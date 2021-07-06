package arraytrie

import helper "github.com/dangermike/wordjumble/trie"

// Trie is a time-efficient method for storing sets of strings where the
// search time is proportional to the length of the string, not the size of the
// set of stored values.
//
// This implementation is space-inefficient in that we can store any arbitrary
// byte string. In practice, we are probably only interested in lower-case,
// unaccented Latin characters. That means each layer of the array is 257
// pointers (2056 bytes on a 64-bit machine) instead of the 27 bytes per array
// we would need (216 bytes) if we only tried to hold the minimum character set.
// If we really wanted to get fancy, we could do some kind of custom base-32
// encoding that ensures that lowercase, unaccented Latin characters are 1 byte.
type trie [257]*trie

var hit = &trie{}

// Load adds a string to the set
func load(at *trie, word []byte) *trie {
	if at == nil {
		at = &trie{}
	}

	if len(word) == 0 {
		at[256] = hit
	} else {
		at[word[0]] = load(at[word[0]], word[1:])
	}

	return at
}

// ContainsString returns whether or not the provided word is in the set of
// loaded words
func (a *Trie) ContainsString(word string) bool {
	return a.Contains([]byte(word))
}

// Contains returns whether or not the provided word is in the set of loaded
// words
func (a *Trie) Contains(word []byte) bool {
	if a == nil {
		return false
	}
	return contains(a.trie, word)
}

func contains(at *trie, word []byte) bool {
	if at == nil {
		return false
	}

	if len(word) == 0 {
		return at[256] == hit
	}

	return contains(at[word[0]], word[1:])
}

// PermuteAll returns all permutations of the provided letters, allowing
// duplicates, that appear in the set of loaded words.
func permuteAll(at *trie, letters []byte, consume bool) [][]byte {
	retval := permuteInner(at, letters, consume)
	for i := 0; i < len(retval); i++ {
		// reverse by starting at each end and swapping as we go
		for a, b := 0, len(retval[i])-1; a < b; a, b = a+1, b-1 {
			retval[i][a], retval[i][b] = retval[i][b], retval[i][a]
		}
	}

	return helper.UniqueifyResults(retval)
}

// permuteInner builds up in reverse, which lets us avoid copying
func permuteInner(at *trie, letters []byte, consume bool) [][]byte {
	if at == nil {
		return nil
	}

	retval := [][]byte{}
	if at[256] == hit {
		retval = append(retval, []byte{})
	}

	for i := 0; i < len(letters); i++ {
		letter := letters[i]
		if at[letter] != nil {
			next := letters
			if consume {
				letters[0], letters[i] = letters[i], letters[0]
				next = next[1:]
			}
			for _, child := range permuteInner(at[letter], next, consume) {
				retval = append(retval, append(child, letter))
			}
		}
	}

	return retval
}
