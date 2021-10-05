package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment variables from .env file
func LoadEnv() {
	fmt.Println("Loading environment Variables...")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

// IsDebugMode checks if debug mode is enabled
func IsDebugMode() bool {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Println("Couldn't find debug value, continuing with `true`")
		debug = true
	}
	return debug
}
