package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
	Mail   MailConfig
	Cookie CookieConfig
}

type CookieConfig struct {
	Secure   bool
	HttpOnly bool
	Domain   string
	SameSite string
}

type DBConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

type ServerConfig struct {
	Port           string
	AllowedOrigins []string
	FrontendURL    string
}

type JWTConfig struct {
	Secret              string
	AccessExpiryMinutes int
	RefreshExpiryDays   int
}

type MailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

var Cfg *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	Cfg = &Config{
		DB: DBConfig{
			User: getEnv("DB_USER", "root"),
			Pass: getEnv("DB_PASS", ""),
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "3306"),
			Name: getEnv("DB_NAME", "educlypro_db"),
		},
		Server: ServerConfig{
			Port:           getEnv("PORT", "8080"),
			AllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ","),
			FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:5173"),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", ""),
			AccessExpiryMinutes: getEnvAsInt("JWT_ACCESS_EXPIRY_MINUTES", 15),
			RefreshExpiryDays:   getEnvAsInt("JWT_REFRESH_EXPIRY_DAYS", 7),
		},
		Mail: MailConfig{
			Host:        getEnv("MAIL_HOST", "smtp.sendgrid.net"),
			Port:        getEnvAsInt("MAIL_PORT", 587),
			Username:    getEnv("MAIL_USERNAME", "apikey"),
			Password:    getEnv("MAIL_PASSWORD", ""),
			Encryption:  getEnv("MAIL_ENCRYPTION", "tls"),
			FromAddress: getEnv("MAIL_FROM_ADDRESS", "mohamedbaya.dev@gmail.com"),
			FromName:    getEnv("MAIL_FROM_NAME", "EduclyPro"),
		},
		Cookie: CookieConfig{
			Secure:   getEnvAsBool("COOKIE_SECURE", false),
			HttpOnly: getEnvAsBool("COOKIE_HTTP_ONLY", true),
			Domain:   getEnv("COOKIE_DOMAIN", ""),
			SameSite: getEnv("COOKIE_SAMESITE", "Lax"),
		},
	}

	if Cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is not set in environment")
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valueStr := strings.ToLower(getEnv(key, ""))
	if valueStr == "true" || valueStr == "1" {
		return true
	}
	if valueStr == "false" || valueStr == "0" {
		return false
	}
	return fallback
}
