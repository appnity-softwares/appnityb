package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	JWTExpiryHours     int
	RefreshExpiryHours int
	CORSOrigins        []string
	UploadDir          string
	MaxUploadSize      int64
	BaseURL            string
	AdminEmail         string
	AdminPassword      string
	Environment        string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	jwtExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if jwtExpiry == 0 {
		jwtExpiry = 24
	}

	refreshExpiry, _ := strconv.Atoi(os.Getenv("REFRESH_EXPIRY_HOURS"))
	if refreshExpiry == 0 {
		refreshExpiry = 168
	}

	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:5174,http://localhost:3000"
	}

	maxUploadSize, _ := strconv.ParseInt(os.Getenv("MAX_UPLOAD_SIZE"), 10, 64)
	if maxUploadSize == 0 {
		maxUploadSize = 10 << 20
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return Config{
		Port:               port,
		DatabaseURL:        databaseURL,
		JWTSecret:          jwtSecret,
		JWTExpiryHours:     jwtExpiry,
		RefreshExpiryHours: refreshExpiry,
		CORSOrigins:        strings.Split(corsOrigins, ","),
		UploadDir:          getEnvOrDefault("UPLOAD_DIR", "./uploads"),
		MaxUploadSize:      maxUploadSize,
		BaseURL:            baseURL,
		AdminEmail:         os.Getenv("ADMIN_EMAIL"),
		AdminPassword:      os.Getenv("ADMIN_PASSWORD"),
		Environment:        env,
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
