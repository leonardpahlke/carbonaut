package compareutils_test

import (
	"reflect"
	"testing"

	"carbonaut.dev/pkg/util/compareutils"
)

func TestCountValuesOfMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected map[int]int
	}{
		{"Empty Map", map[string]int{}, map[int]int{}},
		{"Map with Values", map[string]int{"foo": 1, "bar": 1}, map[int]int{1: 2}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := compareutils.CountValuesOfMap(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("CountValuesOfMap() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestGetListDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"List with Single Duplicate", []string{"A", "A"}, []string{"A"}},
		{"List with Multiple Duplicates", []string{"1", "A", "A1", "B", "C", "A", "A", "DDD", "BB", "BBB", "A1"}, []string{"A", "A1"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := compareutils.GetListDuplicates(tc.input)
			if !reflect.DeepEqual(*result, tc.expected) {
				t.Errorf("GetListDuplicates() = %v, want %v", *result, tc.expected)
			}
		})
	}
}

func TestCheckListContains(t *testing.T) {
	list := []string{"A", "B", "C"}
	tests := []struct {
		name     string
		element  string
		expected bool
	}{
		{"Element Exists", "A", true},
		{"Element Does Not Exist", "D", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := compareutils.CheckListContains(list, tc.element)
			if result != tc.expected {
				t.Errorf("CheckListContains() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Run("basic filtering should return a filtered list", func(t *testing.T) {
		list := []string{"A", "B", "C"}
		list2 := []string{"A", "B", "A", "C"}
		expected1 := []string{"A"}
		expected2 := []string{"A", "A"}

		result1 := compareutils.Filter(list, "A")
		if !reflect.DeepEqual(result1, expected1) {
			t.Errorf("Filter() = %v, want %v", result1, expected1)
		}

		result2 := compareutils.Filter(list2, "A")
		if !reflect.DeepEqual(result2, expected2) {
			t.Errorf("Filter() = %v, want %v", result2, expected2)
		}
	})

	t.Run("deep filtering should return a filtered list", func(t *testing.T) {
		type S1 struct {
			A int
			B int
			C int
		}
		list := []S1{{A: 1, B: 2}, {A: 1, B: 3}, {A: 1, B: 4}}
		expected1 := []S1{{A: 1, B: 2}}

		type S2 struct {
			A  int
			B  int
			S1 S1
		}
		list2 := []S2{{A: 1, B: 2, S1: S1{A: 1, B: 2}}, {A: 1, B: 3, S1: S1{A: 1, B: 3}}, {A: 1, B: 4, S1: S1{A: 1, B: 4}}}
		expected2 := []S2{{A: 1, B: 2, S1: S1{A: 1, B: 2}}}

		result1 := compareutils.Filter(list, S1{A: 1, B: 2})
		if !reflect.DeepEqual(result1, expected1) {
			t.Errorf("Filter() = %v, want %v", result1, expected1)
		}

		result2 := compareutils.Filter(list2, S2{A: 1, B: 2, S1: S1{A: 1, B: 2}})
		if !reflect.DeepEqual(result2, expected2) {
			t.Errorf("Filter() = %v, want %v", result2, expected2)
		}
	})
}

func TestEqual(t *testing.T) {
	type S1 struct {
		A int
		B int
		C int
	}
	map1 := S1{A: 1, B: 2}
	map2 := S1{A: 1, B: 2}
	map3 := S1{A: 1, B: 3}
	map4 := S1{A: 1, B: 2, C: 3}

	tests := []struct {
		name     string
		map1     *S1
		map2     *S1
		expected bool
	}{
		{"Maps are Equal", &map1, &map2, true},
		{"Maps are Not Equal", &map1, &map3, false},
		{"Maps with Different Fields", &map1, &map4, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			equal, err := compareutils.Equal(tc.map1, tc.map2)
			if err != nil {
				t.Fatal("Unexpected error:", err)
			}
			if equal != tc.expected {
				t.Errorf("Equal() = %v, want %v", equal, tc.expected)
			}
		})
	}
}

func TestCompareLists(t *testing.T) {
	cases := []struct {
		name    string
		newList []*int
		oldList []*int
		same    []*int
		missing []*int
		new     []*int
	}{
		{
			name:    "Test integers",
			newList: convertToPointerSlice([]int{1, 2, 3, 4}),
			oldList: convertToPointerSlice([]int{2, 3, 5}),
			same:    convertToPointerSlice([]int{2, 3}),
			missing: convertToPointerSlice([]int{5}),
			new:     convertToPointerSlice([]int{1, 4}),
		},
		{
			name:    "Test empty lists",
			newList: []*int{},
			oldList: []*int{},
			same:    []*int{},
			missing: []*int{},
			new:     []*int{},
		},
		{
			name:    "Test no overlap",
			newList: convertToPointerSlice([]int{1, 2}),
			oldList: convertToPointerSlice([]int{3, 4}),
			same:    []*int{},
			missing: convertToPointerSlice([]int{3, 4}),
			new:     convertToPointerSlice([]int{1, 2}),
		},
		{
			name:    "Test no new items",
			newList: convertToPointerSlice([]int{}),
			oldList: convertToPointerSlice([]int{3, 4}),
			same:    []*int{},
			missing: convertToPointerSlice([]int{3, 4}),
			new:     convertToPointerSlice([]int{}),
		},
		{
			name:    "Test only new items",
			newList: convertToPointerSlice([]int{3, 4}),
			oldList: convertToPointerSlice([]int{}),
			same:    []*int{},
			missing: convertToPointerSlice([]int{}),
			new:     convertToPointerSlice([]int{3, 4}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sameL, missingL, newL := compareutils.CompareLists(tc.newList, tc.oldList)
			if !reflect.DeepEqual(sameL, tc.same) || !reflect.DeepEqual(missingL, tc.missing) || !reflect.DeepEqual(newL, tc.new) {
				t.Errorf("%s failed: expected (%v, %v, %v), got (%v, %v, %v)", tc.name, pointersToStrings(tc.same), pointersToStrings(tc.missing), pointersToStrings(tc.new), pointersToStrings(sameL), pointersToStrings(missingL), pointersToStrings(newL))
			}
		})
	}
}

func convertToPointerSlice(slice []int) []*int {
	ptrSlice := make([]*int, len(slice))
	for i, v := range slice {
		value := v
		ptrSlice[i] = &value
	}
	return ptrSlice
}

func pointersToStrings(pointers []*int) []int {
	slice := make([]int, len(pointers))
	for i, ptr := range pointers {
		if ptr != nil {
			slice[i] = *ptr
		}
	}
	return slice
}
