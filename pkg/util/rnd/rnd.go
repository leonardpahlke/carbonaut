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

package rnd

import (
	"crypto/rand"
	"math/big"
)

// GetNumber selects a random number in a given range
// return -1 if provided input is invalid
func GetNumber(min, max int) int {
	if min > max {
		return -1
	}
	if max == 0 {
		return 0
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		return -1
	}
	n := nBig.Int64() + int64(min)
	return int(n)
}

// GetRandomListSubset creates a random shuffled subset of a provided list with a random number of elements
func GetRandomListSubset[L any](l []L) []L {
	subset := []L{}
	if len(l) == 0 {
		return subset
	}
	subsetSize := GetNumber(1, len(l)+1)
	availableIndexes := []int{}
	for i := 0; i < len(l); i++ {
		availableIndexes = append(availableIndexes, i)
	}
	for i := 0; i < subsetSize; i++ {
		selectedElem := GetNumber(0, len(availableIndexes)-1)
		subset = append(subset, l[availableIndexes[selectedElem]])
		availableIndexes = append(availableIndexes[:selectedElem], availableIndexes[selectedElem+1:]...)
	}
	return subset
}
