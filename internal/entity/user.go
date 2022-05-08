package entity

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
