package dto

type ProductToCart struct {
	ID       int    `json:"id"`
	Name     string `json:"detail"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type ConfirmPayment struct {
	CreditCard string `json:"credit_card"`
}
