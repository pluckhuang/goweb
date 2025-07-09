package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&FeedPullEvent{}, &FeedPushEvent{})
}

// 推拉模型、收发信箱、读写原理:

// a 获取 feed 流, 读流程:
// b 如果是大v,  a 如果关注b , b 产生一对多关系事件， b会把事件写到自己的发件箱中, 之后 a 会读取 b 的发件箱
// b 如果是小v,  a 如果关注b , b 产生一对多关系事件， b会把事件写到粉丝的 收件箱中比如a的, 之后 a 只读自己的收件箱
// b 不管是大v 还是小v, b 产生对 a 事件都会把事件写到 a 的收件箱中
// 这样 a 后取 feed 事件, 需要2部分: 一是获取自己的收件箱事件, 而是获取自己的关注者, 再获取这些关注者的发件箱事件, 最终聚合排序裁剪

// 写流程
// b 对 a 产生 一对一 的事件, 事件写入 a 的收件箱.
// b 产生 一对多 的事件, 看 b 是大v 还是小v, 如果是大v, 则写入 b 的发件箱, 如果是小v, 则写入粉丝的收件箱中.
