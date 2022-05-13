package config

import (
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func SetupDatabaseConnection() *gorm.DB {
	connectString := "host=ec2-34-246-227-219.eu-west-1.compute.amazonaws.com port=5432 user=dwiuwcchyjajhh dbname=d637cf012r303n password=4ee5daba915b8323b2750b319c0df7c5b167cc1d4e4abfbd1406a0b36354e3c4"
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

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}

	dbSQL.Close()
}
