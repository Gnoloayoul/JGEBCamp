package article

import (
	"context"
	"errors"
	// snowflake 雪花算法
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDAO struct {
	// 制作库
	col *mongo.Collection
	// 线上库
	liveCol *mongo.Collection
	node *snowflake.Node

	idGen IDGenerator
}

func (m *MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime, art.Utime = now, now
	// 没有自增主键
	// GLOBAL UNIFY ID (GUID, 全局唯一ID)
	// 这里使用雪花算法生成主键
	id := m.node.Generate().Int64()
	art.Id = id
 	_, err := m.col.InsertOne(ctx, art)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (m *MongoDBDAO) UpdateById(ctx context.Context, art Article) error {
	// 操作制作库
	filter := bson.M{"id": art.Id}
	update := bson.D{bson.E{"$set", bson.M{
		"title": art.Title,
		"comtent": art.Content,
		"utime": time.Now().UnixMilli(),
		"status": art.Status,
	}}}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

func (m *MongoDBDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Article, error) {
	panic("implement me")
}

func (m *MongoDBDAO) GetById(ctx context.Context, id int64) (Article, error) {
	panic("implement me")
}

func (m *MongoDBDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	panic("implement me")
}

func (m *MongoDBDAO) Sync(ctx context.Context, art Article) (int64, error) {

}

func (m *MongoDBDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	panic("implement me")
}


func InitCollections(db *mongo.Database) error {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()
	index := []mongo.IndexModel {
		{
			Keys: bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "ctime", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("articles").Indexes().
		CreateMany(ctx, index)
	return err
}

type IDGenerator func() int64

func NewMongoDBDAOV1(db *mongo.Database, idGen IDGenerator) ArticleDAO {
	return &MongoDBDAO{
		col: db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		idGen: idGen,
	}
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBDAO{
		col: db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		node: node,
	}
}
