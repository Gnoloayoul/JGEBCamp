package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

func (a Article) Abstract() string {
	// 摘要：取前几句
	// 考虑中文问题
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return a.Content
	}
	// 不用纠结会不会截取到一个完整的英文单词
	// 词组、介词，往后找标点符号
	return string(cs[:100])
}

type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

type ArticleStatusV2 string

const (
	// ArticleStatusUnknown
	// 用来避免零值问题（当收到一个零，是用户没传呢，还是传的就是0）
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

func (s ArticleStatus) NonPublished() bool {
	return s != ArticleStatusUnpublished
}

func (s ArticleStatus) Valid() bool {
	return s.ToUint8() > 0
}

func (s ArticleStatus) String() string {
	switch s {
	case ArticleStatusPrivate:
		return "private"
	case ArticleStatusUnpublished:
		return "unpublished"
	case ArticleStatusPublished:
		return "published"
	default:
		return "unknown"
	}
}

// 如果状态很复杂，有很多行为（就是要搞很多方法），状态里面需要一些额外字段
// 就用这个（V1）版本

type ArticleStatusV1 struct {
	Val  uint8
	Name string
}

var (
	ArticleStatusV1Unknown = ArticleStatusV1{Val: 0, Name: "unknown"}
)
