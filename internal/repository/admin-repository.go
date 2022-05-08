package repository

import (
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"gorm.io/gorm"
)

type AdminRepository interface {
	InsertAdmin(admin entity.User) entity.User
	ListAllUsers() []entity.User
	ListAllUsersByActive() []entity.User
	ListAllUsersByTypeAdmin() []entity.User
	ListAllUsersByTypeUser() []entity.User
	BanUser(userID string)
	UnbanUser(userID string)
	NewFeature(feature entity.Feature) entity.Feature
}

type adminConnection struct {
	connection *gorm.DB
}

func NewAdminRepository(dbConn *gorm.DB) AdminRepository {
	return &adminConnection{
		connection: dbConn}
}

func (db *adminConnection) InsertAdmin(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))
	user.TypeUser = "admin"

	db.connection.Create(&user)

	return user
}

func (db *adminConnection) ListAllUsers() []entity.User {
	var users []entity.User
	db.connection.Model(users).Find(&users)
	return users
}

func (db *adminConnection) ListAllUsersByActive() []entity.User {
	var users []entity.User
	active := false
	db.connection.Model(users).Where("active = ?", active).Find(&users)
	return users
}

func (db *adminConnection) ListAllUsersByTypeAdmin() []entity.User {
	var users []entity.User
	typeUser := "admin"
	db.connection.Model(users).Where("type_user = ?", typeUser).Find(&users)
	return users
}

func (db *adminConnection) ListAllUsersByTypeUser() []entity.User {
	var users []entity.User
	typeUser := "user"
	db.connection.Model(users).Where("type_user = ?", typeUser).Find(&users)
	return users
}

func (db *adminConnection) BanUser(userID string) {
	var user entity.User
	db.connection.Model(user).Where("id = ?", userID).Update("active", false)
}

func (db *adminConnection) UnbanUser(userID string) {
	var user entity.User
	db.connection.Model(user).Where("id = ?", userID).Update("active", true)
}

func (db *adminConnection) NewFeature(feature entity.Feature) entity.Feature {

	db.connection.Save(&feature)
	return feature
}
