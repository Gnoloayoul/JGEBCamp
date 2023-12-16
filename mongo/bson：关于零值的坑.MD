# bson 查找数据时由于零值引发的坑
假设要在 mongoDB 中查找一个 `id = 123` 的 Article，可以这么写：
````
// 查找 id = 123
filter := bson.D{bson.E{Key: "id", Value: "123"}}
var art Article
err = col.FindOne(ctx, filter).Decode(&art)
assert.NoError(t, err)
// 打印查找结果
fmt.Printf("%v \n", art)
````
但看到了 `FindOne` 的源码是这样子的：
````
func (coll *Collection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *SingleResult {

	if ctx == nil {
		ctx = context.Background()
	}

	findOpts := make([]*options.FindOptions, 0, len(opts))
	for _, opt := range opts {
		if opt == nil {
			continue
......			
````
尤其是注意到其中的 `filter interface{}` ，也许会想到这样的骚操作：  
用 `Article{Id: 123}` 替换掉上面的 `filter := bson.D{bson.E{Key: "id", Value: "123"}}`  
````
var art Article
err = col.FindOne(ctx, Article{Id: 123}).Decode(&art)
assert.NoError(t, err)
fmt.Printf("%v \n", art)
```` 
结果，翻车了，返回了 `no document in result`，即是 `mongo.ErrNoDocuments` 错误  
**按理来说，`filter` 既然是 `interface` ，传啥都行， 而两者的语义也差不多啊，为什么查不到呢？**  
当通过返回信息中，能找到线索：当查询条件是 `Article{Id: 123}` ，实际上在 `MongoDB` 里是这样子的：  
![](C:\Users\Administrator\Desktop\JGEBCamp\mongo\bson查询传零值.png)  
`Article` 中没被赋值的部分，都被当做**零值**传进去查找  
要避免这个坑，需要对这个结构体 `Article` 做点改造：  
````
// 改造前
type Article struct {
	Id       int64
	Title    string
	Content  string
	AuthorId int64
	Status   uint8
	Ctime    int64
	Utime    int64
}

// 改造后
// `bson:"XXXX,omitempty"`
type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}

````

这里说的是用 `bson` **查询**的情况，在**更新**时也要注意该问题  
进一步来说， **bson** 是类似 **JSON** 的一种序列化与反序列化的协议  
在操作 bson 时，建议使用 _**结构体 + omitempty**_ 这种写法

