package repository

import (
	"log"

	"gorm.io/gorm"
)

// UserRepository is a contract what UserRepository can do to db.
type AuthRepository interface {
	VerifyCodeByEmail(email string, code int) (bool, error)
	FindEmail(email string) int
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

	codeV := db.FindEmail(email)

	if codeV != code {
		log.Println("code not valid")
		return false, nil
	}

	return true, nil
}

/*func (db *authConnection) FindEmail(email string) (em string, err error) {

	type Result struct {
		Email string
	}
	var result Result
	res1 := db.connection.Raw("SELECT name FROM users WHERE email = ?", email).Scan(&result)

	fmt.Println(&res1.U)
	fmt.Println(*res1)
	fmt.Println(res1)

	user := &entity.User{}
	if err := db.connection.Model(user).Where("email = ?", email).Take(user.Email); err != nil {
		log.Println("Error: ", err)
		return "", err.Error
	}

	return user.Email, nil
}*/

func (db *authConnection) FindEmail(emailQuery string) (cv int) {

	var code int
	db.connection.Raw("select code_verify FROM users WHERE email = ?", emailQuery).Scan(&code)

	return code
}
