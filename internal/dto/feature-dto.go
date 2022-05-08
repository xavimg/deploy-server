package dto

type FeatureDTO struct {
	Title string `json:"title" binding:"required,min=5,max=10"`
	Body  string `json:"body" binding:"required,min=10,max=500"`
}
