package entity

import "encoding/json"

type TypeCard string

type TypeUser string

const (
	TypeVisa       TypeCard = "Visa"
	TypeMastercard TypeCard = "Mastercard"
	TypeAmex       TypeCard = "Amex"

	TypeUserNormal TypeUser = "user"
	TypeUserAdmin  TypeUser = "admin"
)

type User struct {
	ID         uint64   `gorm:"primary_key:auto_increment" json:"id"`
	Name       string   `gorm:"type:varchar(255)" json:"name"`
	Email      string   `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password   string   `gorm:"->;<-;not null" json:"-"`
	Token      string   `gorm:"type:varchar(255)" json:"-"`
	Active     bool     `gorm:"default:true" json:"-"`
	Isvalid    bool     `gorm:"default:false" json:"-"`
	CodeVerify int      `gorm:"default:null" json:"-"`
	TypeUser   TypeUser `gorm:"default:user" json:"-"`
	//CreditCard CreditCard `gorm:"type:varchar(255)" json:"credit_card"`
}

type CreditCard struct {
	ISBN string   `gorm:"type:varchar(255)" json:"isbn"`
	Type TypeCard `gorm:"type:varchar(255)" json:"type_card"`
}

type Friend struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
}

type User_Friends struct {
	ID         uint64          `json:"id"`
	IDUser     uint64          `json:"id_user"`
	IDFriend   uint64          `json:"id_friend"`
	Friendlist json.RawMessage `json:"friendlist"`
}

type User_Messages struct {
	ID       uint64
	Receiver float64
	Sender   float64
	Tittle   string
	Detail   string
}
