package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 中间件 gin.HandleFunc
// gin.Default默认使用了两个中间件，一个打日志，另一个恢复panic
// 在中间件里使用goroutine，不能c *gin.Context, 只能是c的复制，不然c.Next()等操作出来的结果就是不可控的

// 原生的"中间件"
func indexHandleFunc(c *gin.Context) {
	fmt.Println("index in ...")
	name, ok := c.Get("name") // 从上下文取值
	if !ok {
		name = "nobody"
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": name,
	})
}

// 新写的m1 统计耗时
func m1(c *gin.Context) {
	fmt.Println("m1 in ...")

	start := time.Now()
	c.Next() // 继续处理后面逻辑
	//c.Abort() // 阻止处理后面逻辑
	cost := time.Since(start)
	fmt.Printf("cost: %v\n", cost)

	fmt.Println("m1 out ...")
}

func m2(c *gin.Context) {
	fmt.Println("m2 in ...")
	////c.Next() // 继续处理后面逻辑
	//c.Abort() // 阻止处理后面逻辑 阻断了响应函数，自然是一片“空白”
	c.Set("name", "空花首席") // 在上下文设置值 （跨中间件存取值）
	fmt.Println("m2 out ...")
}

// 登陆中间件 [不推荐写法]
func authMiddleware0(c *gin.Context) {
	// if 是登录用户
	// c.Next()
	// else
	// c.Abort()
}

// 登录中间件 闭包写法
func authMiddleware(docheck bool) gin.HandlerFunc{
	// 链接数据库
	// 或者其他准备工作
	return func(c *gin.Context) {
		if docheck {
			// 总之就是具体检验逻辑
			// if 是登录用户
			// c.Next()
			// else
			// c.Abort()
		} else {
			c.Next()
		}

	}
}

func main() {
	r :=  gin.Default()

	r.Use(m1, m2) // 将自定义的中间件m1注册上
	r.GET("./index", indexHandleFunc)
	r.GET("./book", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "book",
		})
	})
	r.GET("./shop", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "shop",
		})
	})

	// 路由组注册中间件方法1：
	//xxGroup := r.Group("/xx", authMiddleware(true))
	//{
	//	xxGroup.GET("./index", func(c *gin.Context) {
	//		c.JSON(http.StatusOK, gin.H{
	//			"msg": "xx",
	//		})
	//	})
	//}

	// 路由组注册中间件方法2：
	//xx2Group := r.Group("/xx")
	//xx2Group.Use(authMiddleware(true))
	//{
	//	xxGroup.GET("./index", func(c *gin.Context) {
	//		c.JSON(http.StatusOK, gin.H{
	//			"msg": "xx",
	//		})
	//	})
	//}

	r.Run(":9000")
}