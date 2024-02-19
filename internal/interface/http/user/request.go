package user

type GetUserReq struct {
	ID int `json:"id"`
}

type GetUsersByRoleReq struct {
	Role   string `param:"role" json:"role"`
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
}

type GetUsersLikeUsernameReq struct {
	Username string `query:"username"`
	Limit    int    `query:"limit"`
	Offset   int    `query:"offset"`
}

type DeleteUserReq struct {
	ID int `json:"id"`
}
