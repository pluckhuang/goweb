package domain

import (
	"time"
)

type FeedEvent struct {
	ID int64
	// 以 A 发表了一篇文章为例
	// 如果是 Pull Event，也就是拉模型，那么 Uid 是 A 的id
	// 如果是 Push Event，也就是推模型，那么 Uid 是 A 的某个粉丝的 id
	Uid int64
	// 表示是什么事件,  比如 follow_event
	Type  string
	Ctime time.Time
	Ext   ExtendFields
}
