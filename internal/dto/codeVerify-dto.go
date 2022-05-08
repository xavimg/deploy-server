package dto

type CodeVerifyDTO struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}
