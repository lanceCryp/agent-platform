package config

import (
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Server
	Port string
	Env  string
	
	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	
	// Redis
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
	
	// Kafka
	KafkaBrokers []string
	
	// JWT
	JWTSecret        string
	JWTExpiresIn     time.Duration
	RefreshTokenTTL  time.Duration
	
	// Rate Limiting
	RateLimitReqPerMin int
	RateLimitBurst     int
	
	// Agent Config
	AgentTimeout       time.Duration
	MaxConcurrentTasks int
}

func Load() *Config {
	viper.SetEnvPrefix("AGENTHUB")
	viper.AutomaticEnv()
	
	// Defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("JWT_EXPIRES_IN", "24h")
	viper.SetDefault("REFRESH_TOKEN_TTL", "168h") // 7 days
	viper.SetDefault("RATE_LIMIT_REQ_PER_MIN", 100)
	viper.SetDefault("RATE_LIMIT_BURST", 20)
	viper.SetDefault("AGENT_TIMEOUT", "5m")
	viper.SetDefault("MAX_CONCURRENT_TASKS", 100)
	
	cfg := &Config{
		Port:               viper.GetString("PORT"),
		Env:                viper.GetString("ENV"),
		DBHost:             viper.GetString("DB_HOST"),
		DBPort:             viper.GetInt("DB_PORT"),
		DBUser:             viper.GetString("DB_USER"),
		DBPassword:         viper.GetString("DB_PASSWORD"),
		DBName:             viper.GetString("DB_NAME"),
		DBSSLMode:          viper.GetString("DB_SSL_MODE"),
		RedisHost:          viper.GetString("REDIS_HOST"),
		RedisPort:          viper.GetInt("REDIS_PORT"),
		RedisPassword:      viper.GetString("REDIS_PASSWORD"),
		RedisDB:            viper.GetInt("REDIS_DB"),
		KafkaBrokers:       viper.GetStringSlice("KAFKA_BROKERS"),
		JWTSecret:          viper.GetString("JWT_SECRET"),
		JWTExpiresIn:       viper.GetDuration("JWT_EXPIRES_IN"),
		RefreshTokenTTL:    viper.GetDuration("REFRESH_TOKEN_TTL"),
		RateLimitReqPerMin: viper.GetInt("RATE_LIMIT_REQ_PER_MIN"),
		RateLimitBurst:     viper.GetInt("RATE_LIMIT_BURST"),
		AgentTimeout:       viper.GetDuration("AGENT_TIMEOUT"),
		MaxConcurrentTasks: viper.GetInt("MAX_CONCURRENT_TASKS"),
	}
	
	// Override with environment variables if set
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}
	if env := os.Getenv("ENV"); env != "" {
		cfg.Env = env
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.DBHost = dbHost
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.JWTSecret = jwtSecret
	}
	
	return cfg
}

// LoadFromEnv allows overriding config with environment variables at runtime
func LoadFromEnv() *Config {
	cfg := Load()
	
	// Override with environment variables
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.DBHost = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.DBPort = port
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.DBUser = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.DBPassword = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.DBName = v
	}
	if v := os.Getenv("REDIS_HOST"); v != "" {
		cfg.RedisHost = v
	}
	if v := os.Getenv("KAFKA_BROKERS"); v != "" {
		cfg.KafkaBrokers = []string{v}
	}
	
	return cfg
}
