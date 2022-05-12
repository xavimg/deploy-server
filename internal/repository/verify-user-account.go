package repository

import (
	"errors"

	"gorm.io/gorm"
)

// UserRepository is a contract what UserRepository can do to db.
type AuthRepository interface {
	VerifyCodeByEmail(email string, code int) (bool, error)
	FindEmail(email string) (int, error)
}

type authConnection struct {
	connection *gorm.DB
}

// NewUserRepository is creates a new instance of UserRepository
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authConnection{
		connection: db,
	}
}

func (db *authConnection) VerifyCodeByEmail(email string, code int) (bool, error) {
	codeV, err := db.FindEmail(email)
	if err != nil {
		return false, err
	}

	if codeV != code {
		return false, errors.New("Invalid code")
	}

	return true, nil
}

func (db *authConnection) FindEmail(emailQuery string) (cv int, err error) {
	var code int
	db.connection.Raw("select code_verify FROM users WHERE email = ?", emailQuery).Scan(&code)

	return code, nil
}
