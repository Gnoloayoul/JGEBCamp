package sliceEX

import "fmt"

// SliceV03
// 全要求满足
func SliceV03[T any](a []T, n int) ([]T, T, error) {
	length := len(a)
	if n >= length || n < 0 {
		var zero T
		return a, zero, fmt.Errorf("index out of range, length %d, index %d", length, n)
	}

	val := a[n]
	copy(a[n:], a[n + 1:])
	a = shrink(a[:length - 1])
	return a, val, nil
}

// shrink
// 依据判断，自行缩容
func shrink[T any](src []T) []T {
	capt, length := cap(src), len(src)
	switch {
	case capt <= 64:
		return src
	case capt >= 2048 && (capt / length >= 2):
		res := make([]T, 0, capt >> 1)
		res = append(res, src...)
		return res
	case capt < 2048 && (capt / length >= 4):
		factor := 0.625
		res := make([]T, 0, int(float32(factor) * float32(capt)))
		res = append(res, src...)
		return res
	default:
		return src
	}
}

// SliceV01
// 能用就行的
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

// SliceV02
// 优化了一些速度
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



