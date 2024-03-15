package integration

import (
	"context"
	intrRepov1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intrRepo/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/integration/startup"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	grpcRepo "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/grpc"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type InteractiveRepoGrpcTestSuite struct {
	suite.Suite
	db     *gorm.DB
	rdb    redis.Cmdable
	server *grpcRepo.InteractiveGRPCRepositoryServer
}

func (s *InteractiveRepoGrpcTestSuite) SetupSuite() {
	s.db = startup.InitTestDB()
	s.rdb = startup.InitRedis()
	s.server = startup.InitInteractiveRepoGRPCServer()
}

func (s *InteractiveRepoGrpcTestSuite) TearDownTest() {
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

func (s *InteractiveRepoGrpcTestSuite) TestIncrReadCnt() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64

		wantErr  error
		wantResp *intrRepov1.IncrReadCntResponse
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
			wantResp: &intrRepov1.IncrReadCntResponse{},
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
			wantResp: &intrRepov1.IncrReadCntResponse{},
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
			wantResp: &intrRepov1.IncrReadCntResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveRepoGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			//err := svc.IncrReadCnt(context.Background(), tc.biz, tc.bizId)
			resp, err := svc.IncrReadCnt(context.Background(), &intrRepov1.IncrReadCntRequest{
				Biz: tc.biz, BizId: tc.bizId,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

			// 测试后归零
			tc.after(t)
		})
	}
}

func (s *InteractiveRepoGrpcTestSuite) TestLiked() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr  error
		wantResp *intrRepov1.LikedResponse
	}{
		{
			name: "查找是否点赞，成功",
			before: func(t *testing.T) {
				_, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.Create(dao.UserLikeBiz{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					Uid: 3,
					Ctime:      6,
					Utime:      7,
					Status: 1,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				_, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.UserLikeBiz
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.UserLikeBiz{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					Uid: 3,
					Ctime:      6,
					Utime:      7,
					Status: 1,
				}, data)
			},
			biz:      "test",
			bizId:    2,
			uid:      3,
			wantResp: &intrRepov1.LikedResponse{},
		},
		{
			name: "查找是否点赞，没有",
			before: func(t *testing.T) {
				_, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.Create(dao.UserLikeBiz{
					Id:         2,
					Biz:        "test",
					BizId:      22,
					Uid: 33,
					Ctime:      66,
					Utime:      77,
					Status: 0,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				_, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 测试 db 部分
				var data dao.UserLikeBiz
				err := s.db.Where("id = ?", 2).First(&data).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.UserLikeBiz{
					Id:         2,
					Biz:        "test",
					BizId:      22,
					Uid: 33,
					Ctime:      66,
					Utime:      77,
					Status: 0,
				}, data)
			},
			biz:      "test",
			bizId:    22,
			uid:      33,
			wantResp: &intrRepov1.LikedResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveRepoGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			resp, err := svc.Liked(context.Background(), &intrRepov1.LikedRequest{
				Biz: tc.biz, Id: tc.bizId, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

			// 测试后归零
			tc.after(t)
		})
	}
}

func (s *InteractiveRepoGrpcTestSuite) TestDecrLike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr  error
		wantResp *intrRepov1.DecrLikeResponse
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
			wantResp: &intrRepov1.DecrLikeResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveRepoGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			resp, err := svc.DecrLike(context.Background(), &intrRepov1.DecrLikeRequest{
				Biz: tc.biz, BizId: tc.bizId, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

			// 测试后归零
			tc.after(t)
		})
	}
}

func (s *InteractiveRepoGrpcTestSuite) TestAddCollectionItem() {
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
		wantResp *intrRepov1.AddCollectionItemResponse
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
			wantResp: &intrRepov1.AddCollectionItemResponse{},
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
			wantResp: &intrRepov1.AddCollectionItemResponse{},
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
			wantResp: &intrRepov1.AddCollectionItemResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveRepoGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			resp, err := svc.AddCollectionItem(context.Background(), &intrRepov1.AddCollectionItemRequest{
				Biz: tc.biz, BizId: tc.bizId, Cid: tc.cid, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)

			// 测试后归零
			tc.after(t)
		})
	}
}

func (s *InteractiveRepoGrpcTestSuite) TestGet() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)

		biz   string
		bizId int64

		wantErr error
		wantRes *intrRepov1.GetResponse
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
			wantRes: &intrRepov1.GetResponse{
				Intr: &intrRepov1.Interactive{
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
				err := s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test",
					BizId:      1,
					ReadCnt:    10,
					CollectCnt: 20,
					LikeCnt:    30,
					Ctime:      13,
					Utime:      24,
				}).Error
				assert.NoError(t, err)

				//err = s.rdb.HSet(ctx, "interactive:test:1",
				//	"read_cnt", 10, "collect_cnt", 20).Err()
				//assert.NoError(t, err)
			},
			biz:   "test",
			bizId: 1,
			wantRes: &intrRepov1.GetResponse{
				Intr: &intrRepov1.Interactive{
					Biz:        "test",
					BizId:      1,
					ReadCnt:    10,
					CollectCnt: 20,
					LikeCnt:    30,
				},
			},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveRepoGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			res, err := svc.Get(context.Background(), &intrRepov1.GetRequest{
				Biz: tc.biz, BizId: tc.bizId,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)

		})
	}
}

func (s *InteractiveRepoGrpcTestSuite) TestGetByIds() {
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
		wantRes *intrRepov1.GetByIdsResponse
	}{
		{
			name: "查找成功",
			biz:  "test",
			ids:  []int64{1, 2},
			wantRes: &intrRepov1.GetByIdsResponse{
				Intrs: []*intrRepov1.Interactive{
					{
						Biz:        "test",
						BizId:      1,
						ReadCnt:    1,
						CollectCnt: 2,
						LikeCnt:    3,
					},
					{
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
			wantRes: &intrRepov1.GetByIdsResponse{
				Intrs: []*intrRepov1.Interactive{},
			},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveRepoGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// 运行测试
			res, err := svc.GetByIds(context.Background(), &intrRepov1.GetByIdsRequest{
				Biz: tc.biz, Ids: tc.ids,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func (s *InteractiveRepoGrpcTestSuite) TestBatchIncrReadCnt() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)

		biz   []string
		bizId []int64

		wantErr error
		wantRes *intrRepov1.BatchIncrReadCntResponse
	}{
		{
			name: "批量添加阅读量成功",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test1",
					BizId:      13,
					ReadCnt:    100,
				}).Error
				assert.NoError(t, err)
				err = s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test2",
					BizId:      14,
					ReadCnt:    101,
				}).Error
				assert.NoError(t, err)
				err = s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test3",
					BizId:      14,
					ReadCnt:    102,
				}).Error
				assert.NoError(t, err)
				err = s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test4",
					BizId:      15,
					ReadCnt:    103,
				}).Error
				assert.NoError(t, err)
			},
			biz: []string{"test1", "test2", "test3", "test4"},
			bizId: []int64{13, 14, 15, 16},
			wantErr: nil,
			wantRes: &intrRepov1.BatchIncrReadCntResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveRepoGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			res, err := svc.BatchIncrReadCnt(context.Background(), &intrRepov1.BatchIncrReadCntRequest{
				Biz: tc.biz, BizId: tc.bizId,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)

		})
	}
}

func (s *InteractiveRepoGrpcTestSuite) TestIncrLike(){
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after func(t *testing.T)

		biz   string
		bizId int64
		uid int64

		wantErr error
		wantRes *intrRepov1.IncrLikeResponse
	}{
		{
			name: "添加点赞成功",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				// 写入 db
				err := s.db.WithContext(ctx).Create(&dao.Interactive{
					Id: 1,
					Biz:        "test1",
					BizId:      23,
					LikeCnt: 1,
				}).Error
				assert.NoError(t, err)
				err = s.db.WithContext(ctx).Create(&dao.UserLikeBiz{
					Id: 1,
					Biz:        "test1",
					BizId:      23,
					Uid: 33,
					Status: 0,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				var data dao.Interactive
				err := s.db.WithContext(ctx).Where("biz = ? AND biz_id = ?", "test1", 23).First(&data).Error
				assert.NoError(t, err)
				// 强制归零
				data.Utime = 0
				assert.Equal(t, dao.Interactive{
					Id: 1,
					Biz:        "test1",
					BizId:      23,
					LikeCnt: 2,
				}, data)

				var data2 dao.UserLikeBiz
				err = s.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ?", "test1", 23, 33).First(&data2).Error
				// 强制归零
				data2.Utime = 0
				assert.NoError(t, err)
				assert.Equal(t, dao.UserLikeBiz{
					Id: 1,
					Biz:        "test1",
					BizId:      23,
					Uid: 33,
					Status: 1,
				}, data2)

			},
			biz: "test1",
			bizId: 23,
			uid: 33,
			wantErr: nil,
			wantRes: &intrRepov1.IncrLikeResponse{},
		},
	}

	// 启动服务，准备测试
	svc := startup.InitInteractiveRepoGRPCServer()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化测试
			tc.before(t)

			// 运行测试
			res, err := svc.IncrLike(context.Background(), &intrRepov1.IncrLikeRequest{
				Biz: tc.biz, BizId: tc.bizId, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)

			tc.after(t)
		})
	}
}

func (s *InteractiveRepoGrpcTestSuite) Collected(){}




func TestInteractiveRepoGrpcService(t *testing.T) {
	suite.Run(t, &InteractiveRepoGrpcTestSuite{})
}
