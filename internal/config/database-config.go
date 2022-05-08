package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

// SetupDatabaseConnection is creating a new connection to our database
func SetupDatabaseConnection() *gorm.DB {

	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Println("impossible get .env")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT_DB")
	user := os.Getenv("USER")
	name := os.Getenv("NAME")
	password := os.Getenv("PORT")

	connectString := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s",
		host, port, user, name, password)

	db, err := gorm.Open(postgres.Open(connectString))

	if err != nil {
		panic("Failed to create a connection to database")
	}

	db.AutoMigrate(
		&entity.User{},
		&entity.Feature{},
	)

	return db
}

// CLoseDatabaseConnection method is closing a connection between your app and your database
func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()

	if err != nil {
		panic("Failed to close connection from database")
	}

	dbSQL.Close()
}
