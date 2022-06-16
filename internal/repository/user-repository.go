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
	UpdateUser(user entity.User, userID interface{}, newInfo dto.UserUpdateDTO) (entity.User, error)
	VerifyCredential(email, password string) (interface{}, error)
	VerifyUserExist(userID interface{}) (interface{}, error)
	VerifyUserActive(email interface{}) (entity.User, error)
	IsDuplicateEmail(email string) (ctx *gorm.DB, err error)
	FindByEmail(username string) (entity.User, error)
	ProfileUser(userID interface{}) (*entity.User, error)
	SaveToken(user entity.User, token string) error
	DeleteToken(user entity.User, token string) error
	DeleteAccount(userID float64) error
	GetToken(userID interface{}) (entity.User, error)
	CheckRole(id interface{}) (entity.TypeUser, error)

	AddFriend(id interface{}, user *dto.Friend) error
	ShowFriendlist(id interface{}) ([]*entity.User_Friends, error)
	RemoveFriend(id uint64) error
	IsFriend(id interface{}) (bool, error)
	SendMessage(message entity.User_Messages) error
	ListMessages(id interface{}) ([]*entity.User_Messages, error)
	MessageDetail(id int) (*entity.User_Messages, error)
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

func (db *userConnection) UpdateUser(user entity.User, userID interface{}, newInfo dto.UserUpdateDTO) (entity.User, error) {
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

func (db *userConnection) ProfileUser(userID interface{}) (*entity.User, error) {
	var user *entity.User
	db.connection.Find(&user, userID)

	return user, nil
}

func (db *userConnection) DeleteAccount(userID float64) error {
	var user entity.User
	if err := db.connection.Model(&user).Delete(&user.ID, userID); err != nil {
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

func (db *userConnection) AddFriend(id interface{}, friend *dto.Friend) error {
	if err := db.connection.Raw("INSERT INTO user_friends(id_user, id_friend, friendlist) VALUES(?,?,?);", id, friend.ID, friend).Scan(friend).Error; err != nil {
		return err
	}

	return nil
}

func (db *userConnection) ShowFriendlist(id interface{}) ([]*entity.User_Friends, error) {
	friends := []*entity.User_Friends{}
	result := db.connection.Where("id_user = ?", id).Find(&friends)
	if result.Error != nil {
		return nil, result.Error
	}

	return friends, nil
}

func (db *userConnection) RemoveFriend(id uint64) error {
	friends := []*entity.User_Friends{}
	if err := db.connection.Where("id_friend = ?", id).Delete(&friends).Error; err != nil {
		return err
	}

	return nil
}

func (db *userConnection) IsFriend(id interface{}) (bool, error) {
	friends := entity.User_Friends{}
	var err error
	if err = db.connection.Select("id_user").Where("id_friend = ?", id).Find(&friends).Error; err != nil {
		return false, err
	}
	if friends.IDUser == 0 {
		return false, err
	}
	return true, nil
}

func (db *userConnection) SendMessage(m entity.User_Messages) error {
	if err := db.connection.Raw("INSERT INTO user_messages(receiver, sender, tittle, detail) VALUES(?,?,?,?);", m.Receiver, m.Sender, m.Tittle, m.Detail).Scan(m).Error; err != nil {
		return err
	}

	return nil
}

func (db *userConnection) ListMessages(id interface{}) ([]*entity.User_Messages, error) {
	notifications := []*entity.User_Messages{}
	if err := db.connection.Select("sender, tittle").Where("receiver = ?", id).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func (db *userConnection) MessageDetail(id int) (*entity.User_Messages, error) {
	message := entity.User_Messages{}
	if err := db.connection.Select("detail, sender").Where("id = ?", id).Find(&message).Error; err != nil {
		return nil, err
	}

	return &message, nil
}
