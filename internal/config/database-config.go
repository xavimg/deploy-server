package config

import (
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

// SetupDatabaseConnection is creating a new connection to our database
func SetupDatabaseConnection() *gorm.DB {

	// errEnv := godotenv.Load(".env")
	// if errEnv != nil {
	// 	log.Println("impossible get .env")
	// }

	// host := os.Getenv("HOST")
	// port := os.Getenv("PORT_DB")
	// user := os.Getenv("USER")
	// name := os.Getenv("NAME")
	// password := os.Getenv("PORT")

	connectString := "host=ec2-176-34-211-0.eu-west-1.compute.amazonaws.com port=5432 user=afmzfqrwhhfjfa dbname=df13dp2td99rig password=4090d2bdce300b93316fc2b0831bba6cd55fccd8bf6df71681532d2553fe7a3e"
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
		panic("failed to close connection from database")
	}

	dbSQL.Close()
}
