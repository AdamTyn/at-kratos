package util

import "sort"

type ArrayString []string

func (arr ArrayString) Has(target string) bool {
	sort.Strings(arr)
	index := sort.SearchStrings(arr, target)
	if index < len(arr) && arr[index] == target {
		return true
	}
	return false
}

type ArrayInt []int

func (arr ArrayInt) Has(target int) bool {
	sort.Ints(arr)
	index := sort.SearchInts(arr, target)
	if index < len(arr) && arr[index] == target {
		return true
	}
	return false
}
