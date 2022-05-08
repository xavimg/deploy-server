package repository

import (
	"log"

	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository is a contract what UserRepository can do to db.
type UserRepository interface {
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User, userID string, newInfo dto.UserUpdateDTO) entity.User
	VerifyCredential(email, password string) interface{}
	VerifyUserExist(userID string) interface{}
	VerifyUserActive(email string) entity.User
	IsDuplicateEmail(email string) (ctx *gorm.DB)
	FindByEmail(username string) (entity.User, error)
	ProfileUser(userID string) entity.User
	SaveToken(user entity.User, token string)
	DeleteToken(user entity.User, token string)
	DeleteAccount(userID uint64) error
	GetToken(userID string) entity.User
	CheckRole(id interface{}) (entity.TypeUser, error)
}

type userConnection struct {
	connection *gorm.DB
}

// NewUserRepository is creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))

	db.connection.Create(&user)

	return user
}

func (db *userConnection) UpdateUser(user entity.User, userID string, newInfo dto.UserUpdateDTO) entity.User {

	if newInfo.Name != "" {

		db.connection.Model(user).Where("id = ?", userID).Update("name", newInfo.Name)
	}

	if newInfo.Email != "" {

		db.connection.Model(user).Where("id = ?", userID).Update("email", newInfo.Email)
	}

	if newInfo.Password != "" {

		user.Password = hashAndSalt([]byte(newInfo.Password))

		db.connection.Model(user).Where("id = ?", userID).Update("password", user.Password)
	}

	db.connection.Preload("Characters").Preload("Characters.User").Find(&user)

	return user
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}

	return string(hash)
}

func (db *userConnection) VerifyCredential(email string, password string) interface{} {
	var user entity.User

	res := db.connection.Where("email = ?", email).Take(&user)
	// resactive := db.connection.Model(user).Where("active = ?", true).Take(&user)

	if res == nil {
		return res.Error
	}
	return user
}

func (db *userConnection) VerifyUserExist(userID string) interface{} {
	var user entity.User

	res := db.connection.Where("id = ?", userID).Take(&user)

	if res == nil {
		return res.Error
	}
	return user
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User

	return db.connection.Where("email = ?", email).Take(&user)
}

func (db *userConnection) ProfileUser(userID string) entity.User {
	var user entity.User

	db.connection.Find(&user, userID)

	return user
}

func (db *userConnection) DeleteAccount(userID uint64) error {
	var user entity.User

	if err := db.connection.Delete(&user.ID, userID); err != nil {
		log.Println("Error: ", err)
		return err.Error
	}

	return nil
}

func (db *userConnection) FindByEmail(email string) (user entity.User, err error) {
	var userToFind entity.User

	db.connection.Where("email = ? ", email).Take(&userToFind)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	return user, nil
}

func (db *userConnection) SaveToken(user entity.User, token string) {

	user.Token = token

	db.connection.Save(&user)
}

func (db *userConnection) DeleteToken(user entity.User, s string) {

	user.Token = s

	db.connection.Save(&user)

}

func (db *userConnection) GetToken(UserID string) entity.User {
	var user entity.User

	db.connection.Find(&user, UserID)

	return user
}

func (db *userConnection) VerifyUserActive(email string) entity.User {

	var user entity.User

	db.connection.Find(&user, email)

	return user
}

func (db *userConnection) CheckRole(id interface{}) (typeUser entity.TypeUser, err error) {

	if err := db.connection.Raw("select type_user FROM users WHERE id = ?", id).Scan(&typeUser); err != nil {
		return typeUser, err.Error
	}

	return typeUser, nil

}
