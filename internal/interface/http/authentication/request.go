package authentication

type LoginReq struct {
	Email    string `form:"email" validate:"required_without=Username"`
	Username string `form:"username" validate:"required_without=Email"`
	Password string `form:"password" validate:"required"`
}
