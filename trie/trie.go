package trie

import "sort"

type Trie interface {
	Load(word []byte) bool
	LoadString(word string) bool
	Count() int
	PermuteAll(letters []byte, consume bool) [][]byte
	Contains(letters []byte) bool
	ContainsString(word string) bool
}

func UniqueifyResults(retval [][]byte) [][]byte {
	sort.Slice(retval, func(i, j int) bool {
		return bytesLessThan(retval[i], retval[j])
	})
	shift := 0
	for i := 1; i < len(retval); i++ {
		if bytesEqual(retval[i-1], retval[i]) {
			shift++
		}
		if shift > 0 {
			retval[i-shift] = retval[i]
		}
	}
	retval = retval[0 : len(retval)-shift]

	return retval
}

func bytesLessThan(a, b []byte) bool {
	for len(a) > 0 && len(b) > 0 && a[0] == b[0] {
		a = a[1:]
		b = b[1:]
	}
	if len(b) == 0 {
		return false
	}
	return len(a) == 0 || a[0] < b[0]
}

func bytesEqual(a, b []byte) bool {
	for len(a) > 0 && len(b) > 0 && a[0] == b[0] {
		a = a[1:]
		b = b[1:]
	}
	return len(a) == 0 && len(b) == 0
}
