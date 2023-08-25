package main

import "fmt"

func printfSlice(s []int, sn string) {
	fmt.Printf("%s: %v, len %d, cap: %d \n", sn, s, len(s), cap(s))
}

func shareSlice() {
	s1 := []int{1, 2, 3, 4}
	s2 := s1[2:]
	printfSlice(s1, "s1")
	printfSlice(s2, "s2")

	// s2的改动能影响到s1
	s2[0] = 99
	printfSlice(s1, "s1")
	printfSlice(s2, "s2")

	// 这里可以看到，因为s2追加了199，破坏了结构，s2与s1不能算同源了，因此s1不会反应s2的修改
	s2 = append(s2, 199)
	printfSlice(s1, "s1")
	printfSlice(s2, "s2")

	// 同理，这里的s1与s2是两个不同的东西了
	s2[1] = 1999
	printfSlice(s1, "s1")
	printfSlice(s2, "s2")
}

func main() {
	shareSlice()
}
