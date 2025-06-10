package domain

// Interactive 这个是总体交互的计数
type Interactive struct {
	Biz   string
	BizId int64

	ReadCnt    int64 `json:"read_cnt"`
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`

	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}
