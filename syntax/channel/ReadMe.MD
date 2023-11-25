# 关于 Channel

## TestChannel  
- 如果输入超过了 Channel 的容量
```
	ch2 := make(chan int, 2)
	ch2 <- 123
	ch2 <- 234
	ch2 <- 333
```
运行会返回：  
```
=== RUN   TestChannel
fatal error: all goroutines are asleep - deadlock!
```
- 在一般情况下， Channel 是先进先出的
```
	ch2 := make(chan int, 2)
	ch2 <- 123
	ch2 <- 234
	
    val, ok := <-ch2
	if !ok {
		// ch2 已经被人关了
	}
	println(val, ok)
```
运行会返回：
```
=== RUN   TestChannel
123 true
--- PASS: TestChannel (0.00s)
```
先取一次
```
	ch2 := make(chan int, 2)
	ch2 <- 123
	ch2 <- 234
	
	<-ch // 在这里先把 123 给取掉
    val, ok := <-ch2
	if !ok {
		// ch2 已经被人关了
	}
	println(val, ok)
```
运行会返回：
```
=== RUN   TestChannel
234 true
--- PASS: TestChannel (0.00s)
```

## TestChannelClose
这里演示了，关闭了 Channel， 还是能取完里面的元素  
但在这里还要给这 Channel 发数据，程序就会卡在那里不会往下走  

## TestLoopChannel
这里演示了，唤起一个 go 里写好一个 Channel，然后可用两个方法输出里面的元素  
```
	// 写法1： 直接遍历 ch
	for val := range ch {
		t.Log(val)
	}

	// 写法2： 依次从 Channel 中取
	for {
		val, ok := <-ch
		if !ok {
			break
		}
		t.Log(val)
	}
```