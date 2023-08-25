package main

import (
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/week1/homework/sliceEX"
)

func main() {
	arr := []int{1, 2, 3, 4, 5}
	arr1 := []string{"a", "b", "c", "d", "e"}
	fmt.Println(sliceEX.SliceV01(arr, 3))
	fmt.Println(sliceEX.SliceV02(arr, 3))
	fmt.Println(sliceEX.SliceV03(arr, 3))
	fmt.Println(sliceEX.SliceV03(arr1, 3))
}
