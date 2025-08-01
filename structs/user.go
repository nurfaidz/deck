package structs

type UserResponse struct {
	Id        uint    `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Token     *string `json:"token,omitempty"`
}

type UserCreateRequest struct {
	Username string `json:"username" binding:"required" gorm:"unique; not null"`
	Email    string `json:"email" binding:"required,email" gorm:"unique; not null"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Username string `json:"username" binding:"required" gorm:"unique; not null"`
	Email    string `json:"email" binding:"required,email" gorm:"unique; not null"`
	Password string `json:"password,omitempty"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
