package model

type RegisterUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	RoleName  string `json:"role_name" validate:"required,oneof=admin user guest"`
	Age       int    `json:"age" validate:"required,min=13,max=120"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type EmailVerificationReq struct {
	Code  string `json:"code" validate:"required,len=6"`
	Email string `json:"email" validate:"required,email"`
}
