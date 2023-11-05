package mongo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMongo(t *testing.T) {
	// 控制初始化超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 每个命令（查询）执行之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		// 执行成功
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {},
		// 执行失败
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {},
	}
	ops := options.Client().ApplyURI("mongodb://root:example@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, ops)
	assert.NoError(t, err)

	mdb := client.Database("webook")
	col := mdb.Collection("articles")
	defer func() {
		// 全清数据
		_, err = col.DeleteMany(ctx, bson.D{})
	}()

	res, err := col.InsertOne(ctx, Article{
		Id:      123,
		Title:   "我的标题",
		Content: "我的内容",
	})
	fmt.Printf("id: %s", res.InsertedID)

	// bson
	// 找 id = 123
	filter := bson.D{bson.E{Key: "id", Value: "123"}}
	var art Article
	err = col.FindOne(ctx, filter).Decode(&art)
	assert.NoError(t, err)
	fmt.Printf("%v \n", art)

	// 绝对会翻车的写法2
	art = Article{}
	err = col.FindOne(ctx, Article{Id: 123}).Decode(&art)
	if err == mongo.ErrNoDocuments {
		fmt.Println("没有数据")
	}
	assert.NoError(t, err)
	fmt.Printf("%v \n", art)

	// bson：更新
	set := bson.D{bson.E{Key: "$set",
		// 这里只更新一个 title 字段，就用 bson.E
		// 但要更新多个字段，要用 bson.D
		Value: bson.E{Key: "title", Value: "新的标题"}}}
	updateRes, err := col.UpdateMany(ctx, filter, set)
	assert.NoError(t, err)
	// *UpdateResult.ModifiedCount 更新了多少行数据（计数）
	fmt.Println(updateRes.ModifiedCount)

	// bson: 直接用结构体来更新
	updateRes, err = col.UpdateMany(ctx, filter, bson.D{
		bson.E{Key: "$set", Value: Article{Title: "新的标题2", AuthorId: 123456}}})
	assert.NoError(t, err)
	// *UpdateResult.ModifiedCount 更新了多少行数据（计数）
	fmt.Println(updateRes.ModifiedCount)

	// bson: 复合条件查询 or
	// 写法1
	or := bson.A{bson.D{bson.E{"id", 123}},
		bson.D{bson.E{"id", 456}}}
	// 写法2
	//or := bson.A{bson.M{"id": 123}, bson.M{"id": 456}}
	orRes, err := col.Find(ctx, bson.D{bson.E{"$or", or}})
	assert.NoError(t, err)
	var resArt []Article
	err = orRes.All(ctx, &resArt)
	assert.NoError(t, err)

	// bson: 复合条件查询 and
	// 写法1
	and := bson.A{bson.D{bson.E{"id", 123}},
		bson.D{bson.E{"title", "我的标题1"}}}
	// 写法2
	//and := bson.A{bson.M{"id": 123}, bson.M{"title": 我的标题1}}
	andRes, err := col.Find(ctx, bson.D{bson.E{"$and", and}})
	assert.NoError(t, err)
	var aresArt []Article
	err = andRes.All(ctx, &aresArt)
	assert.NoError(t, err)

	// bson: in 查询
	in := bson.D{bson.E{"id", bson.D{bson.E{"$in", []any{123, 456}}}}}
	inRes, err := col.Find(ctx, in)
	resArt = []Article{}
	err = inRes.All(ctx, &resArt)
	assert.NoError(t, err)

	// bson: 建索引
	col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	})

	// bson: 建索引 写法2
	idxRex, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"author_id": 1},
		},
	})
	assert.NoError(t, err)

	// bson: 删除
	delRes, err := col.DeleteOne(ctx, filter)
	assert.NoError(t, err)
	// *DeleteResult.DeletedCount: The number of documents deleted 删除了多少行
	fmt.Println(delRes.DeletedCount)

}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
