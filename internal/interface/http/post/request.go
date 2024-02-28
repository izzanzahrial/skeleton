package post

type CreatPostReq struct {
	UserID  int64  `form:"id" json:"user_id" validate:"required"`
	Title   string `form:"title" json:"title" validate:"required"`
	Content string `form:"content" json:"content" validate:"required"`
}

type GetPostByUserIDReq struct {
	UserID int64 `param:"id" json:"user_id" validate:"required"`
}

type GetPostsFullTextReq struct {
	Keyword string `query:"keyword" json:"keyword"`
	Limit   int    `query:"limit" json:"limit" validate:"omitempty,gte=10"`
	Offset  int    `query:"offset" json:"offset" validate:"omitempty,gte=1"`
}
