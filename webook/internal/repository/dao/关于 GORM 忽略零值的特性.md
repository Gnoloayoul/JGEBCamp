# 关于 GORM 忽略零值的特性
`err := dao.db.WithContext(ctx).Update(&art).Error`   
在 GORM 里更新数据库表时会很容易写成这句，这句就是运用了 GORM 忽略零值的特性  
所谓的忽略零值，就是当输入的字段为零值，那么不会更新  
像上面这一句，就会用主键进行更新  
不是不可以用，而是**不建议**用。因为这样的可读性会很差：一眼上去，你根本不知道哪些字段是更新了（换句话说不清楚这句代码有什么行为）  
这里有个可行的写法：直接显性写出哪些字段会更新
````
err := dao.db.WithContext(ctx).Model(&art).
    Where("id=?", art.Id).
    Updates(map[string]any{
        "title": art.Title,
        "content": art.Content,
        "utime": art.Utime,
}).Error
````
其中 art 的构造：
````
type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 限定长度：1024
	Title string `gorm:"type=varchar(1024)"`
	// BLOB：mysql 中适合存大文本数据的数据类型
	Content string `gorm:"type=BLOB"`
	// 仅仅给 AuthorId 上索引
	AuthorId int64 `gorm:"index"`
	Ctime int64
	Utime int64
}
````
上面的代码的动作就是：一次更新，会更新到 `title` 、 `content` 、 `utime` 这三个字段
