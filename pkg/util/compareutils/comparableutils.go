package compareutils

import (
	"bytes"
	"encoding/json"
)

func CountValuesOfMap[K comparable, V comparable](m map[K]V) map[V]int {
	r := map[V]int{}
	for _, v := range m {
		r[v]++
	}
	return r
}

// GetListDuplicates iterates over a provided list and returns all elements that appear more than once
func GetListDuplicates[E comparable](l []E) *[]E {
	c := map[E]int{}
	for i := range l {
		c[l[i]]++
	}
	duplicates := []E{}
	for j := range c {
		if c[j] > 1 {
			duplicates = append(duplicates, j)
		}
	}
	return &duplicates
}

// CompareLists takes two slices of any comparable type and returns three slices:
// sameItems: items present in both newList and oldList
// missingItems: items present in oldList but not in newList
// newItems: items present in newList but not in oldList
func CompareLists[T comparable](newList, oldList []*T) (sameItems, missingItems, newItems []*T) {
	newSet := make(map[T]bool)
	oldSet := make(map[T]bool)

	sameItems = []*T{}
	missingItems = []*T{}
	newItems = []*T{}

	// Fill set for newList
	for i := range newList {
		newSet[*newList[i]] = true
	}

	// Fill set for oldList and determine same and missing items
	for i := range oldList {
		oldSet[*oldList[i]] = true
		if newSet[*oldList[i]] {
			sameItems = append(sameItems, oldList[i])
		} else {
			missingItems = append(missingItems, oldList[i])
		}
	}

	// Determine new items
	for i := range newList {
		if !oldSet[*newList[i]] {
			newItems = append(newItems, newList[i])
		}
	}

	return sameItems, missingItems, newItems
}

// CheckListContains checks if a list contains a specific element
func CheckListContains[E comparable](l []E, e E) bool {
	for i := range l {
		if l[i] == e {
			return true
		}
	}
	return false
}

// Equal checks if two structs are equal
// this errors if the json marshal fails
func Equal[E any](e1, e2 *E) (bool, error) {
	b1, err := json.Marshal(e1)
	if err != nil {
		return false, err
	}
	b2, err := json.Marshal(e2)
	if err != nil {
		return false, err
	}
	return bytes.Equal(b1, b2), nil
}

// Filter retrieves all elements from a list that match an element
func Filter[E comparable](l []E, e E) []E {
	var r []E
	for i := range l {
		if l[i] == e {
			r = append(r, l[i])
		}
	}
	return r
}

// CollectKeys retrieves all map keys
func CollectKeys[E comparable, V any](m map[E]V) []E {
	keys := make([]E, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
