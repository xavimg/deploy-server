package dto

type MessageDTO struct {
	From    float64
	To      float64 `json:"to"`
	Tittle  string  `json:"tittle"`
	Message string  `json:"detail"`
}

type NotificationMessageDTO struct {
	Remitent float64 `json:"from"`
	Tittle   string  `json:"tittle"`
}

type BodyMessageDTO struct {
	Detail string `json:"detail"`
}
