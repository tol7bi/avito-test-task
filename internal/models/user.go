package models

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // employee, moderator
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
