// 本例是群内讨论出现的一个例子

package main

import "fmt"

// 这里的result肯定是42，但我在思考的时候还担心这个result能不能传出来，有可能是6，还有可能是0
// 事实是，还是42
// return的6，相当于上下文，被闭包给扑捉到
// 虽然是有return，在有defer的存在下，还是得执行完defer的闭包才能退出函数
// 这样结果就是6*7=42
func f() (result int) {
	defer func() {
		// result is accessed after it was set to 6 by the return statement
		result *= 7
	}()
	return 6
}

func main() {
	fmt.Println(f())
}
