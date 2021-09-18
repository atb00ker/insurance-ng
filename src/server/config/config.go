package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	fmt.Println("Loading environment Variables...")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func IsDebugMode() bool {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Println("Couldn't find debug value, continuing with `true`")
		debug = true
	}
	return debug
}
