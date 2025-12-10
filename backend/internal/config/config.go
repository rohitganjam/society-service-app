package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port        string
	Environment string

	// Database (Supabase)
	DatabaseURL     string
	SupabaseURL     string
	SupabaseAnonKey string

	// JWT
	JWTSecret          string
	JWTExpiryHours     int
	RefreshExpiryHours int

	// Razorpay
	RazorpayKeyID        string
	RazorpayKeySecret    string
	RazorpayWebhookSecret string

	// Firebase Cloud Messaging
	FCMServerKey string

	// MSG91 (SMS)
	MSG91AuthKey  string
	MSG91SenderID string
	MSG91FlowID   string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		// Server
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Database
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		SupabaseURL:     getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey: getEnv("SUPABASE_ANON_KEY", ""),

		// JWT
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpiryHours:     getEnvInt("JWT_EXPIRY_HOURS", 1),
		RefreshExpiryHours: getEnvInt("REFRESH_EXPIRY_HOURS", 168), // 7 days

		// Razorpay
		RazorpayKeyID:         getEnv("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret:     getEnv("RAZORPAY_KEY_SECRET", ""),
		RazorpayWebhookSecret: getEnv("RAZORPAY_WEBHOOK_SECRET", ""),

		// FCM
		FCMServerKey: getEnv("FCM_SERVER_KEY", ""),

		// MSG91
		MSG91AuthKey:  getEnv("MSG91_AUTH_KEY", ""),
		MSG91SenderID: getEnv("MSG91_SENDER_ID", ""),
		MSG91FlowID:   getEnv("MSG91_FLOW_ID", ""),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
