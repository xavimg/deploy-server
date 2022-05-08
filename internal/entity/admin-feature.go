package entity

type Feature struct {
	ID    uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Title string `gorm:"type:varchar(255)" json:"title"`
	Body  string `gorm:"type:varchar(255)" json:"body"`
}
