package arraytrie

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
type Trie [257]*Trie

var hit = &Trie{}

func New() *Trie {
	return nil
}

// Load adds a string to the set
func Load(at *Trie, word []byte) *Trie {
	if at == nil {
		at = &Trie{}
	}

	if len(word) == 0 {
		at[256] = hit
	} else {
		at[word[0]] = Load(at[word[0]], word[1:])
	}

	return at
}

// LoadString adds a string to the set
func LoadString(at *Trie, word string) *Trie {
	return Load(at, []byte(word))
}

// ContainsString returns whether or not the provided word is in the set of
// loaded words
func ContainsString(at *Trie, word string) bool {
	return Contains(at, []byte(word))
}

// Contains returns whether or not the provided word is in the set of loaded
// words
func Contains(at *Trie, word []byte) bool {
	if at == nil {
		return false
	}

	if len(word) == 0 {
		return at[256] == hit
	}

	return Contains(at[word[0]], word[1:])
}

// PermuteAll returns all permutations of the provided letters, allowing
// duplicates, that appear in the set of loaded words.
func PermuteAll(at *Trie, letters []byte) [][]byte {
	retval := permuteInner(at, letters)
	for i := 0; i < len(retval); i++ {
		// reverse by starting at each end and swapping as we go
		for a, b := 0, len(retval[i])-1; a < b; a, b = a+1, b-1 {
			retval[i][a], retval[i][b] = retval[i][b], retval[i][a]
		}
	}
	return retval
}

// permuteInner builds up in reverse, which lets us avoid copying
func permuteInner(at *Trie, letters []byte) [][]byte {
	if len(letters) == 0 || at == nil {
		return nil
	}

	retval := [][]byte{}
	if at[256] == hit {
		retval = append(retval, []byte{})
	}

	for _, letter := range letters {
		if at[letter] != nil {
			for _, child := range permuteInner(at[letter], letters) {
				retval = append(retval, append(child, letter))
			}
		}
	}

	return retval
}
