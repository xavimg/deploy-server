package entity

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID        int64     `json:"id"`
	Detail    string    `json:"detail"`
	Price     int       `json:"price"`
	CreatedAT time.Time `json:"created_at,omitempty"`
}

type Carts struct {
	IDCart     int64           `json:"id_cart,omitempty"`
	Product    json.RawMessage `json:"products,omitempty"`
	CreditCard string          `json:"credit_card,omitempty"`
}

type Buy struct {
	ID        int64  `json:"id"`
	Iduser    int64  `json:"iduser"`
	Createdat string `json:"createdat"`
}
