package config

import (
	"fmt"
	"insurance-ng/src/server/models"
	"log"
	"net/url"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectToDb() {
	dsn := url.URL{
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASS")),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		Path:     os.Getenv("DB_NAME"),
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	var err error
	Database, err = gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	migrateDb()
}

func migrateDb() {
	Database.AutoMigrate(&models.Users{})
	Database.AutoMigrate(&models.UserConsents{})
}
