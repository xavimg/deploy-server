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
	UpdateUser(user entity.User, userID string, newInfo dto.UserUpdateDTO) (entity.User, error)
	VerifyCredential(email, password string) (interface{}, error)
	VerifyUserExist(userID interface{}) (interface{}, error)
	VerifyUserActive(email interface{}) (entity.User, error)
	IsDuplicateEmail(email string) (ctx *gorm.DB, err error)
	FindByEmail(username string) (entity.User, error)
	ProfileUser(userID string) (entity.User, error)
	SaveToken(user entity.User, token string) error
	DeleteToken(user entity.User, token string) error
	DeleteAccount(userID uint64) error
	GetToken(userID interface{}) (entity.User, error)
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

func (db *userConnection) UpdateUser(user entity.User, userID string, newInfo dto.UserUpdateDTO) (entity.User, error) {
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

	return user, nil
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(hash)
}

func (db *userConnection) VerifyCredential(email string, password string) (interface{}, error) {
	var user entity.User
	res := db.connection.Where("email = ?", email).Take(&user)
	if res == nil {
		return nil, res.Error
	}

	return user, nil
}

func (db *userConnection) VerifyUserExist(userID interface{}) (interface{}, error) {
	var user entity.User
	res := db.connection.Where("id = ?", userID).Take(&user)
	if res == nil {
		return nil, res.Error
	}

	return user, nil
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB, err error) {
	var user entity.User
	return db.connection.Where("email = ?", email).Take(&user), nil
}

func (db *userConnection) ProfileUser(userID string) (entity.User, error) {
	var user entity.User
	db.connection.Find(&user, userID)

	return user, nil
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

func (db *userConnection) SaveToken(user entity.User, token string) error {
	user.Token = token
	db.connection.Save(&user)
	return nil
}

func (db *userConnection) DeleteToken(user entity.User, s string) error {
	user.Token = s
	db.connection.Save(&user)
	return nil
}

func (db *userConnection) GetToken(UserID interface{}) (entity.User, error) {
	var user entity.User
	db.connection.Find(&user, UserID)
	return user, nil
}

func (db *userConnection) VerifyUserActive(email interface{}) (entity.User, error) {
	var user entity.User
	db.connection.Find(&user, email)
	return user, nil
}

func (db *userConnection) CheckRole(id interface{}) (typeUser entity.TypeUser, err error) {
	if err := db.connection.Raw("select type_user FROM users WHERE id = ?", id).Scan(&typeUser); err != nil {
		return typeUser, err.Error
	}

	return typeUser, nil

}
