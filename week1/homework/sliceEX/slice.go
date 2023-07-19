package sliceEX

func idxCheck(a []int, n int) bool {
	if n > len(a) {
		return false
	}
	return true
}

func SliceV01(a []int, n int) []int {
	if !idxCheck(a, n) {return []int{}}
	ans := []int{}
	ans = append(ans, a[:n]...)
	ans = append(ans, a[n + 1:]...)
	return ans
}

func SliceV02(a []int, n int) []int {
	if !idxCheck(a, n) {return []int{}}
	length := len(a) - 1
	ans := make([]int, length, length)
	copy(ans[:len(a[:n])], a[:n])
	copy(ans[len(a[:n]):], a[n + 1:])
	return ans
}

func idxCheckT[T any](a []T, n int) bool {
	if n > len(a) {
		return false
	}
	return true
}

func SliceV03[T any](a []T, n int) []T {
	if !idxCheckT(a, n) {
		return []T{}
	}
	length := len(a) - 1
	ans := make([]T, length, length)
	copy(ans[:len(a[:n])], a[:n])
	copy(ans[len(a[:n]):], a[n + 1:])
	return ans
}