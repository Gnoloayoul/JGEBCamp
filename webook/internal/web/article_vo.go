package web

import "github.com/Gnoloayoul/JGEBCamp/webook/internal/domain"

// VO view object， 就是对标前端的

type ArticleVO struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract string `json:"abstract"`
	Content  string `json:"content"`
	// 注意一点，状态这个东西，可以是前端来处理，也可以是后端来处理
	// 0 -> unknown -> 未知状态
	// 1 -> 未发表 ->
	// 后端处理： 涉及发版、国际化时
	Status uint8  `json:"status"`
	Author string `json:"author"`
	Ctime  string `json:"ctime"`
	Utime  string `json:"utime"`
}

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}