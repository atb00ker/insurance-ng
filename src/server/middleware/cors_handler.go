package middleware

import (
	"insurance-ng/src/server/config"
	"os"
	"strings"

	"github.com/rs/cors"
)

// CorsMiddleware creates middleware to add cors header to response
func CorsMiddleware() *cors.Cors {
	allowedCors := strings.Split(os.Getenv("APP_CORS"), ",")
	return cors.New(cors.Options{
		AllowedOrigins:   allowedCors,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		Debug:            config.IsDebugMode(),
	})
}
