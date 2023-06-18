package models

type User struct {
	Name     string  `json:"name,omitempty" validate:"required"`
	Username string  `json:"username" validate:"required"`
	Email    string  `json:"email,omitempty" validate:"required"`
	Password string  `json:"password,omitempty" validate:"required"`
	Type     string  `json:"title,omitempty"`
	Balance  float32 `json:"balance"`
	Status   string  `json:"status"`
	IsAdmin  bool    `json:"isAdmin"`
}
