package config

import (
	"fmt"
	"insurance-ng/src/server/models"
	"log"
	"net/url"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func ConnectToDb() {
	fmt.Println("Connecting to database...")
	dsn := url.URL{
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASS")),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		Path:     os.Getenv("DB_NAME"),
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	var err error
	dbLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             2,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	Database, err = gorm.Open(postgres.Open(dsn.String()), &gorm.Config{Logger: dbLogger})

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	migrateDb()
}

func migrateDb() {
	fmt.Println("Migrating models...")
	Database.AutoMigrate(&models.Users{})
	Database.AutoMigrate(&models.UserConsents{})
	Database.AutoMigrate(&models.UserPlanScores{})
	Database.AutoMigrate(&models.UserExistingInsurance{})
	Database.AutoMigrate(&models.Insurance{})
	InitInsuranceSeed()
}
