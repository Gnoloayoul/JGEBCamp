package sliceEX

import "fmt"

func SliceV01(a []int, n int) ([]int, int, error) {
	length := len(a)
	if n >= length || n < 0 {
		var zero int
		return a, zero, fmt.Errorf("index out of range, length %d, index %d", length, n)
	}

	ans := []int{}
	val := a[n]
	ans = append(ans, a[:n]...)
	ans = append(ans, a[n + 1:]...)
	return ans, val, nil
}

func SliceV02(a []int, n int) ([]int, int, error) {
	length := len(a)
	if n >= length || n < 0 {
		var zero int
		return a, zero, fmt.Errorf("index out of range, length %d, index %d", length, n)
	}

	ans := make([]int, length - 1, length - 1)
	val := a[n]
	copy(ans[:len(a[:n])], a[:n])
	copy(ans[len(a[:n]):], a[n + 1:])
	return ans, val, nil
}

func SliceV03[T any](a []T, n int) ([]T, T, error) {
	length := len(a)
	if n >= length || n < 0 {
		var zero T
		return a, zero, fmt.Errorf("index out of range, length %d, index %d", length, n)
	}

	val := a[n]
	copy(a[n:], a[n + 1:])
	a = a[:length - 1]
	return a, val, nil
}