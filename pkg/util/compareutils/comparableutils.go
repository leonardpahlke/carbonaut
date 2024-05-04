/*
Copyright 2023 CARBONAUT AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
func CompareLists[T comparable](newList, oldList []T) (sameItems, missingItems, newItems []T) {
	newSet := make(map[T]bool)
	oldSet := make(map[T]bool)

	sameItems = []T{}
	missingItems = []T{}
	newItems = []T{}

	// Fill set for newList
	for _, item := range newList {
		newSet[item] = true
	}

	// Fill set for oldList and determine same and missing items
	for _, item := range oldList {
		oldSet[item] = true
		if newSet[item] {
			sameItems = append(sameItems, item)
		} else {
			missingItems = append(missingItems, item)
		}
	}

	// Determine new items
	for _, item := range newList {
		if !oldSet[item] {
			newItems = append(newItems, item)
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
