package repository

import (
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"gorm.io/gorm"
)

type AdminRepository interface {
	InsertAdmin(admin entity.User) (entity.User, error)
	ListAllUsers() ([]entity.User, error)
	ListAllUsersByActive() ([]entity.User, error)
	ListAllUsersByTypeAdmin() ([]entity.User, error)
	ListAllUsersByTypeUser() ([]entity.User, error)
	BanUser(userID string) error
	UnbanUser(userID string) error
	NewFeature(feature entity.Feature) (entity.Feature, error)
}

type adminConnection struct {
	connection *gorm.DB
}

func NewAdminRepository(dbConn *gorm.DB) AdminRepository {
	return &adminConnection{
		connection: dbConn}
}

func (db *adminConnection) InsertAdmin(user entity.User) (entity.User, error) {
	user.Password = hashAndSalt([]byte(user.Password))
	user.TypeUser = "admin"

	if err := db.connection.Create(&user); err != nil {
		return entity.User{}, err.Error
	}

	return user, nil
}

func (db *adminConnection) ListAllUsers() ([]entity.User, error) {
	var users []entity.User
	if err := db.connection.Model(users).Find(&users); err != nil {
		return users, err.Error
	}

	return users, nil
}

func (db *adminConnection) ListAllUsersByActive() ([]entity.User, error) {
	var users []entity.User
	active := false
	if err := db.connection.Model(users).Where("active = ?", active).Find(&users); err != nil {
		return users, err.Error
	}

	return users, nil
}

func (db *adminConnection) ListAllUsersByTypeAdmin() ([]entity.User, error) {
	var users []entity.User
	typeUser := "admin"
	if err := db.connection.Model(users).Where("type_user = ?", typeUser).Find(&users); err != nil {
		return users, err.Error
	}

	return users, nil
}

func (db *adminConnection) ListAllUsersByTypeUser() ([]entity.User, error) {
	var users []entity.User
	typeUser := "user"
	if err := db.connection.Model(users).Where("type_user = ?", typeUser).Find(&users); err != nil {
		return users, err.Error
	}

	return users, nil
}

func (db *adminConnection) BanUser(userID string) error {
	var user entity.User
	if err := db.connection.Model(user).Where("id = ?", userID).Update("active", false); err != nil {
		return err.Error
	}

	return nil
}

func (db *adminConnection) UnbanUser(userID string) error {
	var user entity.User
	if err := db.connection.Model(user).Where("id = ?", userID).Update("active", true); err != nil {
		return err.Error
	}

	return nil
}

func (db *adminConnection) NewFeature(feature entity.Feature) (entity.Feature, error) {
	if err := db.connection.Save(&feature); err != nil {
		return entity.Feature{}, err.Error
	}
	return feature, nil
}
