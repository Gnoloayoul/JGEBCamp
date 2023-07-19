package main

import (
	"fmt"
	"reflect"
	"runtime"
)

// result: all 10
func func1() {
	for i := 0; i < 10; i++ {
		defer func() {
			fmt.Println(i)
		}()
	}
}

// result: 9-0
func func2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			fmt.Println(val)
		}(i)
	}
}

// result: 9-0
func func3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			fmt.Println(j)
		}()
	}
}

func printfFunc(funcX ...func()) {
	getName := func(f interface{}) string {
		return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	}
	for _, f := range funcX {
		fmt.Println("result from", getName(f))
		f()
		fmt.Println("++++++++++++++++++++++")
	}
}

func main() {
	printfFunc(func1, func2, func3)
}
