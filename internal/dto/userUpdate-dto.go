package dto

type UserUpdateDTO struct {
	Name     string `json:"name,omitempty" binding:"required,min=5,max=6"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
