package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取path参数
// param返回的都是string
// 留心路径冲突

func main() {
	r := gin.Default()

	r.GET("/user/:name/:age", func(c *gin.Context) {
		// 获取路劲参数
		name := c.Param("name")
		age := c.Param("age")
		c.JSON(http.StatusOK, gin.H{
			"name": name,
			"age":  age,
		})
	})
	r.GET("/blog/:year/:month", func(c *gin.Context) {
		// 获取路劲参数
		year := c.Param("year")
		month := c.Param("month")
		c.JSON(http.StatusOK, gin.H{
			"year":  year,
			"month": month,
		})
	})

	r.Run(":9090")
}
