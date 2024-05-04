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
