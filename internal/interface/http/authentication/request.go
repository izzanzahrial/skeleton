package authentication

type LoginReq struct {
	Email    string `form:"email" validate:"email,required_if=Username"`
	Username string `form:"username" validate:"required_if=Email"`
	Password string `form:"password" validate:"required"`
}
