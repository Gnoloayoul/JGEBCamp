package integration

import (
	"context"
	intrv1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intr/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/grpc"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/integration/startup"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type InteractiveGrpcTestSuite struct {
	suite.Suite
	db     *gorm.DB
	rdb    redis.Cmdable
	server *grpc.InteractiveServiceServer
}

func (s *InteractiveGrpcTestSuite) SetupSuite() {
	s.db = startup.InitTestDB()
	s.rdb = startup.InitRedis()
	s.server = startup.InitInteractiveGRPCServer()
}

func (s *InteractiveGrpcTestSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := s.db.Exec("TRUNCATE TABLE `interactives`").Error
	assert.NoError(s.T(), err)

	err = s.db.Exec("TRUNCATE TABLE `user_like_bizs`").Error
	assert.NoError(s.T(), err)

	err = s.db.Exec("TRUNCATE TABLE `user_collection_bizs`").Error
	assert.NoError(s.T(), err)

	// clear Redis
	err = s.rdb.FlushDB(ctx).Err()
	assert.NoError(s.T(), err)
}

func (s *InteractiveGrpcTestSuite) TestIncrReadCnt() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64

		wantErr  error
		wantResp *intrv1.IncrReadCntResponse
	}{
		{
			name: "增加成功， db 和 redis ",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.Create(dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
					Utime:      7,
				}).Error
				assert.NoError(t, err)

				// 写入 redis
				err = s.rdb.HSet(ctx, "interactive:test:2",
					"read_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.Interactive
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Id:    1,
					Biz:   "test",
					BizId: 2,
					// （查询）阅读 + 1
					ReadCnt:    4,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
				}, data)

				// 测试 redis 部分 + 清除数据
				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "read_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:      "test",
			bizId:    2,
			wantResp: &intrv1.IncrReadCntResponse{},
		},
		{
			name: "增加成功， db 有， redis（缓存）没有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.WithContext(ctx).Create(dao.Interactive{
					Id:         3,
					Biz:        "test",
					BizId:      3,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
					Utime:      7,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.Interactive
				err := s.db.Where("id = ?", 3).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Id:    3,
					Biz:   "test",
					BizId: 3,
					// （查询）阅读 + 1
					ReadCnt:    4,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
				}, data)

				// 测试 redis 部分 + 清除数据
				cnt, err := s.rdb.Exists(ctx, "interactive:test:3").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:      "test",
			bizId:    3,
			wantResp: &intrv1.IncrReadCntResponse{},
		},
		{
			name:   "增加成功，都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.Interactive
				err := s.db.Where("biz = ? AND biz_id = ?", "test", 4).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 0)
				assert.True(t, data.Ctime > 0)
				assert.True(t, data.Id > 0)
				data.Id = 0
				data.Utime = 0
				data.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Biz:     "test",
					BizId:   4,
					ReadCnt: 1,
				}, data)

				// 测试 redis 部分 + 清除数据
				cnt, err := s.rdb.Exists(ctx, "interactive:test:4").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:      "test",
			bizId:    4,
			wantResp: &intrv1.IncrReadCntResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			//err := svc.IncrReadCnt(context.Background(), tc.biz, tc.bizId)
			resp, err := svc.IncrReadCnt(context.Background(), &intrv1.IncrReadCntRequest{
				Biz: tc.biz, BizId: tc.bizId,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

			// 测试后归零
			tc.after(t)
		})
	}
}

func (s *InteractiveGrpcTestSuite) TestLike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr  error
		wantResp *intrv1.LikeResponse
	}{
		{
			name: "点赞成功， db 和 redis 都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.Create(dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
					Utime:      7,
				}).Error
				assert.NoError(t, err)

				// 写入 redis
				err = s.rdb.HSet(ctx, "interactive:test:2",
					"like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.Interactive
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					// 点赞 + 1
					LikeCnt: 6,
					Ctime:   6,
				}, data)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("biz = ? AND biz_id = ? AND uid = ?",
					"test", 2, 123).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Id > 0)
				assert.True(t, likeBiz.Ctime > 0)
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Id = 0
				likeBiz.Ctime = 0
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserLikeBiz{
					Biz:    "test",
					BizId:  2,
					Uid:    123,
					Status: 1,
				}, likeBiz)

				// 测试 redis 部分 + 清除数据
				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:      "test",
			bizId:    2,
			uid:      123,
			wantResp: &intrv1.LikeResponse{},
		},
		{
			name:   "点赞成功， db 和 redis 都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.Interactive
				err := s.db.Where("biz = ? AND biz_id = ?", "test", 3).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Id > 0)
				assert.True(t, data.Ctime > 0)
				assert.True(t, data.Utime > 0)
				data.Id = 0
				data.Ctime = 0
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Biz:     "test",
					BizId:   3,
					LikeCnt: 1,
				}, data)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("biz = ? AND biz_id = ? AND uid = ?",
					"test", 3, 123).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Id > 0)
				assert.True(t, likeBiz.Ctime > 0)
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Id = 0
				likeBiz.Ctime = 0
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserLikeBiz{
					Biz:    "test",
					BizId:  3,
					Uid:    123,
					Status: 1,
				}, likeBiz)

				// 测试 redis 部分 + 清除数据
				cnt, err := s.rdb.Exists(ctx, "interactive:test:2").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:      "test",
			bizId:    3,
			uid:      123,
			wantResp: &intrv1.LikeResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			resp, err := svc.Like(context.Background(), &intrv1.LikeRequest{
				Biz: tc.biz, BizId: tc.bizId, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

			// 测试后归零
			tc.after(t)
		})
	}
}

func (s *InteractiveGrpcTestSuite) TestDisLike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr  error
		wantResp *intrv1.CancelLikeResponse
	}{
		{
			name: "取消点赞成功， db 和 redis 都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.Create(dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					Ctime:      6,
					Utime:      7,
				}).Error
				assert.NoError(t, err)

				err = s.db.Create(dao.UserLikeBiz{
					Id:     1,
					Biz:    "test",
					BizId:  2,
					Uid:    123,
					Ctime:  6,
					Utime:  7,
					Status: 1,
				}).Error
				assert.NoError(t, err)

				// 写入 redis
				err = s.rdb.HSet(ctx, "interactive:test:2",
					"like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.Interactive
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Utime > 7)
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					// 点赞 - 1
					LikeCnt: 4,
					Ctime:   6,
				}, data)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("id = ?", 1).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Utime > 7)
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserLikeBiz{
					Id:     1,
					Biz:    "test",
					BizId:  2,
					Uid:    123,
					Ctime:  6,
					Status: 0,
				}, likeBiz)

				// 测试 redis 部分 + 清除数据
				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 2, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:      "test",
			bizId:    2,
			uid:      123,
			wantResp: &intrv1.CancelLikeResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			resp, err := svc.CancelLike(context.Background(), &intrv1.CancelLikeRequest{
				Biz: tc.biz, BizId: tc.bizId, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

			// 测试后归零
			tc.after(t)
		})
	}
}

func (s *InteractiveGrpcTestSuite) TestCollect() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		cid   int64
		uid   int64

		wantErr  error
		wantResp *intrv1.CollectResponse
	}{
		{
			name: "收藏成功， db 和 redis ",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test",
					BizId:      3,
					CollectCnt: 10,
					Ctime:      123,
					Utime:      234,
				}).Error
				assert.NoError(t, err)

				// 写入 redis
				err = s.rdb.HSet(ctx, "interactive:test:3",
					"collect_cnt", 10).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.Interactive
				err := s.db.WithContext(ctx).
					Where("biz = ? AND biz_id = ?", "test", 3).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Id > 0)
				assert.True(t, data.Ctime > 0)
				assert.True(t, data.Utime > 0)
				data.Id = 0
				data.Ctime = 0
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Biz:        "test",
					BizId:      3,
					CollectCnt: 11,
				}, data)

				var cdata dao.UserCollectionBiz
				err = s.db.WithContext(ctx).
					Where("uid = ? AND biz = ? AND biz_id = ?", 1, "test", 3).
					First(&cdata).Error
				assert.NoError(t, err)
				assert.True(t, cdata.Id > 0)
				assert.True(t, cdata.Ctime > 0)
				assert.True(t, cdata.Utime > 0)
				cdata.Id = 0
				cdata.Ctime = 0
				cdata.Utime = 0
				assert.Equal(t, dao.UserCollectionBiz{
					Biz:   "test",
					BizId: 3,
					Cid:   1,
					Uid:   1,
				}, cdata)

				// 测试 redis 部分 + 清除数据
				cnt, err := s.rdb.HGet(ctx, "interactive:test:3", "collect_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 11, cnt)
			},
			biz:      "test",
			bizId:    3,
			cid:      1,
			uid:      1,
			wantResp: &intrv1.CollectResponse{},
		},
		{
			name: "收藏成功， db 有 redis(缓存) 没有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test",
					BizId:      2,
					CollectCnt: 10,
					Ctime:      123,
					Utime:      234,
				}).Error
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.Interactive
				err := s.db.WithContext(ctx).
					Where("biz = ? AND biz_id = ?", "test", 2).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Id > 0)
				assert.True(t, data.Ctime > 0)
				assert.True(t, data.Utime > 0)
				data.Id = 0
				data.Ctime = 0
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Biz:        "test",
					BizId:      2,
					CollectCnt: 11,
				}, data)

				var cdata dao.UserCollectionBiz
				err = s.db.WithContext(ctx).
					Where("uid = ? AND biz = ? AND biz_id = ?", 1, "test", 2).
					First(&cdata).Error
				assert.NoError(t, err)
				assert.True(t, cdata.Id > 0)
				assert.True(t, cdata.Ctime > 0)
				assert.True(t, cdata.Utime > 0)
				cdata.Id = 0
				cdata.Ctime = 0
				cdata.Utime = 0
				assert.Equal(t, dao.UserCollectionBiz{
					Biz:   "test",
					BizId: 2,
					Cid:   1,
					Uid:   1,
				}, cdata)
			},
			biz:      "test",
			bizId:    2,
			cid:      1,
			uid:      1,
			wantResp: &intrv1.CollectResponse{},
		},
		{
			name:   "收藏成功， db 和 redis(缓存) 都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.Interactive
				err := s.db.WithContext(ctx).
					Where("biz = ? AND biz_id = ?", "test", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.Id > 0)
				assert.True(t, data.Ctime > 0)
				assert.True(t, data.Utime > 0)
				data.Id = 0
				data.Ctime = 0
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Biz:        "test",
					BizId:      1,
					CollectCnt: 1,
				}, data)

				var cdata dao.UserCollectionBiz
				err = s.db.WithContext(ctx).
					Where("uid = ? AND biz = ? AND biz_id = ?", 1, "test", 1).
					First(&cdata).Error
				assert.NoError(t, err)
				assert.True(t, cdata.Id > 0)
				assert.True(t, cdata.Ctime > 0)
				assert.True(t, cdata.Utime > 0)
				cdata.Id = 0
				cdata.Ctime = 0
				cdata.Utime = 0
				assert.Equal(t, dao.UserCollectionBiz{
					Biz:   "test",
					BizId: 1,
					Cid:   1,
					Uid:   1,
				}, cdata)

				cnt, err := s.rdb.Exists(ctx, "interactive:test:1").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:      "test",
			bizId:    1,
			cid:      1,
			uid:      1,
			wantResp: &intrv1.CollectResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			resp, err := svc.Collect(context.Background(), &intrv1.CollectRequest{
				Biz: tc.biz, BizId: tc.bizId, Cid: tc.cid, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

			// 测试后归零
			tc.after(t)
		})
	}
}

func (s *InteractiveGrpcTestSuite) TestGet() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr error
		wantRes *intrv1.GetResponse
	}{
		{
			name: "全取出来了，没缓存",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test",
					BizId:      12,
					ReadCnt:    100,
					CollectCnt: 200,
					LikeCnt:    300,
					Ctime:      123,
					Utime:      234,
				}).Error
				assert.NoError(t, err)
			},
			biz:   "test",
			bizId: 12,
			uid:   123,
			wantRes: &intrv1.GetResponse{
				Intr: &intrv1.Interactive{
					Biz:        "test",
					BizId:      12,
					ReadCnt:    100,
					CollectCnt: 200,
					LikeCnt:    300,
				},
			},
		},
		{
			name: "全取出来了，命中缓存， 用户已经点赞收藏",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.WithContext(ctx).Create(&dao.UserLikeBiz{
					Biz:    "test",
					BizId:  3,
					Uid:    123,
					Ctime:  123,
					Utime:  124,
					Status: 1,
				}).Error
				assert.NoError(t, err)

				err = s.db.WithContext(ctx).Create(&dao.UserCollectionBiz{
					Cid:   1,
					Biz:   "test",
					BizId: 3,
					Uid:   123,
					Ctime: 123,
					Utime: 234,
				}).Error
				assert.NoError(t, err)

				err = s.rdb.HSet(ctx, "interactive:test:3",
					"read_cnt", 0, "collect_cnt", 1).Err()
				assert.NoError(t, err)
			},
			biz:   "test",
			bizId: 3,
			uid:   123,
			wantRes: &intrv1.GetResponse{
				Intr: &intrv1.Interactive{
					BizId:      3,
					CollectCnt: 1,
					Collected:  true,
					Liked:      true,
				},
			},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			res, err := svc.Get(context.Background(), &intrv1.GetRequest{
				Biz: tc.biz, BizId: tc.bizId, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)

		})
	}
}

func (s *InteractiveGrpcTestSuite) TestGetByIds() {
	t := s.T()
	perCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 准备数据
	for i := 1; i < 5; i++ {
		i := int64(i)
		err := s.db.WithContext(perCtx).
			Create(&dao.Interactive{
				Id:         i,
				Biz:        "test",
				BizId:      i,
				ReadCnt:    i,
				CollectCnt: i + 1,
				LikeCnt:    i + 2,
			}).Error
		assert.NoError(t, err)
	}

	testCases := []struct {
		name   string
		before func(t *testing.T)

		biz string
		ids []int64

		wantErr error
		wantRes *intrv1.GetByIdsResponse
	}{
		{
			name: "查找成功",
			biz:  "test",
			ids:  []int64{1, 2},
			wantRes: &intrv1.GetByIdsResponse{
				Intrs: map[int64]*intrv1.Interactive{
					1: {
						Biz:        "test",
						BizId:      1,
						ReadCnt:    1,
						CollectCnt: 2,
						LikeCnt:    3,
					},
					2: {
						Biz:        "test",
						BizId:      2,
						ReadCnt:    2,
						CollectCnt: 3,
						LikeCnt:    4,
					},
				},
			},
		},
		{
			name: "没有对应数据",
			biz:  "test",
			ids:  []int64{100, 200},
			wantRes: &intrv1.GetByIdsResponse{
				Intrs: map[int64]*intrv1.Interactive{},
			},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// 运行测试
			res, err := svc.GetByIds(context.Background(), &intrv1.GetByIdsRequest{
				Biz: tc.biz, BizIds: tc.ids,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestInteractiveGrpcService(t *testing.T) {
	suite.Run(t, &InteractiveGrpcTestSuite{})
}
