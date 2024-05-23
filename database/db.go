package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/davidalvarez305/budgeting/models"
)

type connection struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

var DB gorm.DB

func Connect() (*gorm.DB, error) {
	conn := connection{
		host:     os.Getenv("POSTGRES_HOST"),
		port:     os.Getenv("POSTGRES_PORT"),
		user:     os.Getenv("PGUSER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		dbName:   os.Getenv("POSTGRES_DB"),
	}

	db, err := gorm.Open(postgres.Open(connToString(conn)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		FullSaveAssociations: true,
	})

	if err != nil {
		return db, err
	}

	db.AutoMigrate(
		&models.Lead{},
		&models.LeadMarketing{},
		&models.VendingType{},
		&models.VendingLocation{},
	)

	DB = *db

	return db, nil
}

func connToString(info connection) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		info.host, info.port, info.user, info.password, info.dbName)
}
