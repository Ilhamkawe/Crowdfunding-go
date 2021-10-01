package user

type RegisterInputUser struct {
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" bindung:"required,email"`
}

type FormInputRegister struct {
	Name       string `form:"name" binding:"required"`
	Occupation string `form:"occupation" binding:"required"`
	Email      string `form:"email" binding:"required,email"`
	Password   string `form:"password" binding:"required"`
}

type FormUpdateRegister struct {
	ID         int
	Name       string `form:"name" binding:"required"`
	Occupation string `form:"occupation" binding:"required"`
	Email      string `form:"email" binding:"required,email"`
}
