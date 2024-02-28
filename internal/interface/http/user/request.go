package user

type GetUserReq struct {
	ID int `json:"id" validate:"required,gt=0"`
}

type GetUsersByRoleReq struct {
	Role   string `param:"role" json:"role" validate:"required"`
	Limit  int    `query:"limit" validate:"omitempty,gte=10,lte=100"`
	Offset int    `query:"offset" validate:"omitempty,gte=1"`
}

type GetUsersLikeUsernameReq struct {
	Username string `query:"username" validate:"required"`
	Limit    int    `query:"limit" validate:"omitempty,gte=10,lte=100"`
	Offset   int    `query:"offset" validate:"omitempty,gte=1"`
}

type DeleteUserReq struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type SignUpUserReq struct {
	Email    string `form:"email" validate:"required,email"`
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type SignUpAdminReq struct {
	Email    string `form:"email" validate:"required,email"`
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type UpdateUserReq struct {
	ID       int     `param:"id" json:"id" validate:"required,gte=1"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Username *string `json:"username" validate:"omitempty,alpha`
	Password *string `json:"password" validate:"omitempty"`
}
