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

package compareutils_test

import (
	"reflect"
	"testing"

	"carbonaut.dev/pkg/util/compareutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("compareutils", func() {
	Describe("CountValuesOfMap", func() {
		It("should return an empty map again", func() {
			Expect(compareutils.CountValuesOfMap(map[string]int{})).To(Equal(map[int]int{}))
		})

		It("should be counted up", func() {
			Expect(compareutils.CountValuesOfMap(map[string]int{"foo": 1, "bar": 1})).To(Equal(map[int]int{1: 2}))
		})
	})

	Describe("GetListDuplicates", func() {
		listWithDuplicates1 := []string{"A", "A"}
		listWithDuplicates2 := []string{"1", "A", "A1", "B", "C", "A", "A", "DDD", "BB", "BBB", "A1"}
		It("should return two a single value with two duplicate values in a list ", func() {
			Expect(*compareutils.GetListDuplicates(listWithDuplicates1)).To(Equal([]string{"A"}))
		})
		It("should return multiple values with mixed list of duplicate and none duplicates ", func() {
			Expect(*compareutils.GetListDuplicates(listWithDuplicates2)).To(Equal([]string{"A", "A1"}))
		})
	})

	Describe("CheckListContains", func() {
		list := []string{"A", "B", "C"}
		It("should return true", func() {
			Expect(compareutils.CheckListContains(list, "A")).To(BeTrue())
		})
		It("should return false", func() {
			Expect(compareutils.CheckListContains(list, "D")).To(BeFalse())
		})
	})

	Describe("Check if two maps are equal", func() {
		type S1 struct {
			A int
			B int
			C int
		}
		map1 := S1{A: 1, B: 2}
		map2 := S1{A: 1, B: 2}
		map3 := S1{A: 1, B: 3}
		map4 := S1{A: 1, B: 2, C: 3}
		It("should return true", func() {
			eql, err := compareutils.Equal(&map1, &map2)
			Expect(err).ToNot(HaveOccurred())
			Expect(eql).To(BeTrue())
		})
		It("should return false", func() {
			eql, err := compareutils.Equal(&map1, &map3)
			Expect(err).ToNot(HaveOccurred())
			Expect(eql).To(BeFalse())
		})
		It("should return false", func() {
			eql, err := compareutils.Equal(&map1, &map4)
			Expect(err).ToNot(HaveOccurred())
			Expect(eql).To(BeFalse())
		})
		It("should not manipulate the interfaces provided", func() {
			m1Tmp := map1
			m2Tmp := map2
			eql, err := compareutils.Equal(&m1Tmp, &m2Tmp)
			Expect(err).ToNot(HaveOccurred())
			Expect(eql).To(BeTrue())
			Expect(m1Tmp).To(Equal(map1))
			Expect(m2Tmp).To(Equal(map2))
		})
	})

	Describe("Filter", func() {
		It("basic filtering should return a filtered list", func() {
			list := []string{"A", "B", "C"}
			list2 := []string{"A", "B", "A", "C"}
			Expect(compareutils.Filter(list, "A")).To(Equal([]string{"A"}))
			Expect(compareutils.Filter(list2, "A")).To(Equal([]string{"A", "A"}))
		})

		It("deep filtering should return a filtered list", func() {
			type S1 struct {
				A int
				B int
				C int
			}
			list := []S1{{A: 1, B: 2}, {A: 1, B: 3}, {A: 1, B: 4}}
			Expect(compareutils.Filter(list, S1{A: 1, B: 2})).To(Equal([]S1{{A: 1, B: 2}}))

			type S2 struct {
				A  int
				B  int
				S1 S1
			}

			list2 := []S2{{A: 1, B: 2, S1: S1{A: 1, B: 2}}, {A: 1, B: 3, S1: S1{A: 1, B: 3}}, {A: 1, B: 4, S1: S1{A: 1, B: 4}}}
			Expect(compareutils.Filter(list2, S2{A: 1, B: 2, S1: S1{A: 1, B: 2}})).To(Equal([]S2{{A: 1, B: 2, S1: S1{A: 1, B: 2}}}))
		})
	})
})

func TestCompareLists(t *testing.T) {
	cases := []struct {
		name    string
		newList []int
		oldList []int
		same    []int
		missing []int
		new     []int
	}{
		{
			name:    "Test integers",
			newList: []int{1, 2, 3, 4},
			oldList: []int{2, 3, 5},
			same:    []int{2, 3},
			missing: []int{5},
			new:     []int{1, 4},
		},
		{
			name:    "Test empty lists",
			newList: []int{},
			oldList: []int{},
			same:    []int{},
			missing: []int{},
			new:     []int{},
		},
		{
			name:    "Test no overlap",
			newList: []int{1, 2},
			oldList: []int{3, 4},
			same:    []int{},
			missing: []int{3, 4},
			new:     []int{1, 2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			same, missing, new := compareutils.CompareLists(tc.newList, tc.oldList)
			if !reflect.DeepEqual(same, tc.same) || !reflect.DeepEqual(missing, tc.missing) || !reflect.DeepEqual(new, tc.new) {
				t.Errorf("%s failed: expected (%v, %v, %v), got (%v, %v, %v)", tc.name, tc.same, tc.missing, tc.new, same, missing, new)
			}
		})
	}
}
