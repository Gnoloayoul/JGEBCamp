package sliceEX_test

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

func TestSliceV03(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputSlice []interface{}
		inputN     int
		expected   []interface{}
	}{
		{[]interface{}{1, 2, 3, 4, 5}, 2, []interface{}{1, 2, 4, 5}},
		{[]interface{}{10, 20, 30, 40, 50}, 3, []interface{}{10, 20, 30, 50}},
		{[]interface{}{100}, 0, []interface{}{}},
		{[]interface{}{"a", "b", "c", "d", "e", "f", "g"}, 3, []interface{}{"a", "b", "c", "e", "f", "g"}},
		{[]interface{}{}, 0, []interface{}{}},
	}

	for _, test := range tests {
		result := sliceEX.SliceV03(test.inputSlice, test.inputN)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %v, %d\nExpected: %v\nGot: %v", test.inputSlice, test.inputN, test.expected, result)
		}
	}
}

func BenchmarkSliceV01(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5}
	index := 2
	for i := 0; i < b.N; i++ {
		sliceEX.SliceV01(slice, index)
	}
}

func BenchmarkSliceV02(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5}
	index := 2
	for i := 0; i < b.N; i++ {
		sliceEX.SliceV02(slice, index)
	}
}

func BenchmarkSliceV03I(b *testing.B) {
	slice := []interface{}{1, 2, 3, 4, 5}
	index := 2
	for i := 0; i < b.N; i++ {
		sliceEX.SliceV03(slice, index)
	}
}

func BenchmarkSliceV03A(b *testing.B) {
	slice := []interface{}{"a", "b", "c", "d", "e"}
	index := 2
	for i := 0; i < b.N; i++ {
		sliceEX.SliceV03(slice, index)
	}
}