package sliceEX

import (
	"github.com/Gnoloayoul/JGEBCamp/week1/homework/sliceEX"
	"reflect"
	"testing"
)

func TestSliceV01(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputSlice []int
		inputN     int
		expected   []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{1, 2, 4, 5}},
		{[]int{10, 20, 30, 40, 50}, 3, []int{10, 20, 30, 50}},
		{[]int{100}, 0, []int{}},
		{[]int{}, 0, []int{}},
	}

	for _, test := range tests {
		result := sliceEX.SliceV01(test.inputSlice, test.inputN)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %v, %d\nExpected: %v\nGot: %v", test.inputSlice, test.inputN, test.expected, result)
		}
	}
}

func TestSliceV02(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputSlice []int
		inputN     int
		expected   []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{1, 2, 4, 5}},
		{[]int{10, 20, 30, 40, 50}, 3, []int{10, 20, 30, 50}},
		{[]int{100}, 0, []int{}},
		{[]int{}, 0, []int{}},
	}

	for _, test := range tests {
		result := sliceEX.SliceV02(test.inputSlice, test.inputN)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %v, %d\nExpected: %v\nGot: %v", test.inputSlice, test.inputN, test.expected, result)
		}
	}
}
