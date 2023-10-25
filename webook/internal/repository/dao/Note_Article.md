# mysql 的制作库如何设定索引？
这里有这么一个表
````
type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 限定长度：1024
	Title string `gorm:"type=varchar(1024)"`
	// BLOB：mysql 中适合存大文本数据的数据类型
	content string `gorm:"type=BLOB"`
	AuthorId int64
	Ctime int64
	Utime int64
}
````
如何给这表设定索引？
换个问题，就是：在帖子这里，是什么样的**查询**场景？  
**【where】**  
- 对于创作者来说，是不是看草稿箱，看所有自己的文章？  
`SELECT * FROM articles WHERE author_id = XXXX[作者ID]`  
- 单独查询某一篇  
`SELECT * FROM articles WHERE id = XXXX[文章ID]`    
产品经理会告诉你，要安装 **创建时间** 的倒序排序  
````
SELECT * FROM articles WHERE author_id = XXXX[作者ID] ORDER BY `ctime` DESC;  
````    
最佳选择，就是要在 Author_id 和 ctime 上创建联合索引  
写法：  
````
type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	Title string `gorm:"type=varchar(1024)"`
	content string `gorm:"type=BLOB"`
	AuthorId int64 `gorm:"index=aid_ctime"`
	Ctime int64 `gorm:"ndex=aid_ctime"`
	Utime int64
}
````
## 问题，这两哪个好？
- Author_id 和 ctime 上创建联合索引    
- Author_id 单建索引


