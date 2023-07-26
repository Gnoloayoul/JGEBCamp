package sliceEX

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSliceV01(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputSlice []int
		inputN     int
		expected   []int
		expectedVal int
		expErr error
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{1, 2, 4, 5}, 3, nil},
		{[]int{10, 20, 30, 40, 50}, 3, []int{10, 20, 30, 50}, 40, nil},
		{[]int{100}, 0, []int{}, 100, nil},
		{[]int{}, 0, []int{}, 0, fmt.Errorf("index out of range, length 0, index 0")},
	}

	for _, test := range tests {
		result, val, err := SliceV01(test.inputSlice, test.inputN)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %v, %d\nExpected: %v, %d\nGot: %v, %d", test.inputSlice, test.inputN, test.expected, test.expectedVal, result, val)
		}
		if !reflect.DeepEqual(err, test.expErr) {
			t.Errorf("want: %v\nbut get: %v", test.expErr, err)
		}
	}
}

func TestSliceV02(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputSlice []int
		inputN     int
		expected   []int
		expectedVal int
		expErr error
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{1, 2, 4, 5}, 3, nil},
		{[]int{10, 20, 30, 40, 50}, 3, []int{10, 20, 30, 50}, 40, nil},
		{[]int{100}, 0, []int{}, 100, nil},
		{[]int{}, 0, []int{}, 0, fmt.Errorf("index out of range, length 0, index 0")},
	}

	for _, test := range tests {
		result, val, err := SliceV02(test.inputSlice, test.inputN)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %v, %d\nExpected: %v, %d\nGot: %v, %d", test.inputSlice, test.inputN, test.expected, test.expectedVal, result, val)
		}
		if !reflect.DeepEqual(err, test.expErr) {
			t.Errorf("want: %v\nbut get: %v", test.expErr, err)
		}
	}
}

func TestSliceV03(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inputSlice []interface{}
		inputN     int
		expected   []interface{}
		expectedVal interface{}
		expErr error
	}{
		{[]interface{}{1, 2, 3, 4, 5}, 2, []interface{}{1, 2, 4, 5}, 3, nil},
		{[]interface{}{10, 20, 30, 40, 50}, 3, []interface{}{10, 20, 30, 50}, 40, nil},
		{[]interface{}{100}, 0, []interface{}{}, 100, nil},
		{[]interface{}{"a", "b", "c", "d", "e", "f", "g"}, 3, []interface{}{"a", "b", "c", "e", "f", "g"}, "d", nil},
		{[]interface{}{}, 0, []interface{}{}, nil, fmt.Errorf("index out of range, length 0, index 0")},
	}

	for _, test := range tests {
		result, val, err := SliceV03(test.inputSlice, test.inputN)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %v, %d\nExpected: %v, %v\nGot: %v, %v", test.inputSlice, test.inputN, test.expected, test.expectedVal, result, val)
		}
		if !reflect.DeepEqual(err, test.expErr) {
			t.Errorf("want: %v\nbut get: %v", test.expErr, err)
		}
	}
}

func BenchmarkSliceV01(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5}
	index := 2
	for i := 0; i < b.N; i++ {
		SliceV01(slice, index)
	}
}

func BenchmarkSliceV02(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5}
	index := 2
	for i := 0; i < b.N; i++ {
		SliceV02(slice, index)
	}
}

func BenchmarkSliceV03I(b *testing.B) {
	slice := []interface{}{1, 2, 3, 4, 5}
	index := 2
	for i := 0; i < b.N; i++ {
		SliceV03(slice, index)
	}
}

func BenchmarkSliceV03A(b *testing.B) {
	slice := []interface{}{"a", "b", "c", "d", "e"}
	index := 2
	for i := 0; i < b.N; i++ {
		SliceV03(slice, index)
	}
}