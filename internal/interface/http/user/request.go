package user

type GetUserReq struct {
	ID int `json:"id" validate:"gt=0,required"`
}

type GetUsersByRoleReq struct {
	Role   string `param:"role" json:"role,required"`
	Limit  int    `query:"limit" validate:"gt=0,lte=100"`
	Offset int    `query:"offset" validate:"gt=0"`
}

type GetUsersLikeUsernameReq struct {
	Username string `query:"username" validate:"required"`
	Limit    int    `query:"limit" validate:"gt=0,lte=100"`
	Offset   int    `query:"offset" validate:"gt=0"`
}

type DeleteUserReq struct {
	ID int `json:"id" validate:"gt=0,required"`
}

type SignUpUserReq struct {
	Email    string `form:"email" validate:"email,required"`
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type SignUpAdminReq struct {
	Email    string `form:"email" validate:"email,required"`
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}
