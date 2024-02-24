package post

type CreatPostReq struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Title   string `form:"title" json:"title" validate:"required"`
	Content string `form:"content" json:"content" validate:"required"`
}

type GetPostByUserIDReq struct {
	UserID int64 `json:"user_id" validate:"required"`
}

type GetPostsFullTextReq struct {
	Title   string `form:"title" json:"title" validate:"required"`
	Content string `form:"content" json:"content" validate:"required"`
	Limit   int    `query:"limit" json:"limit"`
	Offset  int    `query:"offset" json:"offset"`
}
