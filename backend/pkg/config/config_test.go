package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("ENV", "testing")
	os.Setenv("JWT_SECRET", "test-jwt-secret")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ENV")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Port)
	}

	if cfg.Env != "testing" {
		t.Errorf("Expected env testing, got %s", cfg.Env)
	}

	if cfg.JWTSecret != "test-jwt-secret" {
		t.Errorf("Expected JWT secret test-jwt-secret, got %s", cfg.JWTSecret)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Clear all environment variables
	os.Unsetenv("PORT")
	os.Unsetenv("ENV")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("KAFKA_BROKERS")

	cfg := Load()

	// Check defaults
	if cfg.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Port)
	}

	if cfg.Env != "development" {
		t.Errorf("Expected default env development, got %s", cfg.Env)
	}

	if cfg.DBPort != 5432 {
		t.Errorf("Expected default DB port 5432, got %d", cfg.DBPort)
	}

	if cfg.RedisPort != 6379 {
		t.Errorf("Expected default Redis port 6379, got %d", cfg.RedisPort)
	}

	if cfg.RateLimitReqPerMin != 100 {
		t.Errorf("Expected default rate limit 100, got %d", cfg.RateLimitReqPerMin)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("REDIS_HOST", "redis.example.com")
	os.Setenv("KAFKA_BROKERS", "kafka.example.com:9092")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("KAFKA_BROKERS")
	}()

	cfg := LoadFromEnv()

	if cfg.DBHost != "db.example.com" {
		t.Errorf("Expected DB host db.example.com, got %s", cfg.DBHost)
	}

	if cfg.DBPort != 5433 {
		t.Errorf("Expected DB port 5433, got %d", cfg.DBPort)
	}

	if cfg.DBUser != "testuser" {
		t.Errorf("Expected DB user testuser, got %s", cfg.DBUser)
	}

	if cfg.DBPassword != "testpass" {
		t.Errorf("Expected DB password testpass, got %s", cfg.DBPassword)
	}

	if cfg.DBName != "testdb" {
		t.Errorf("Expected DB name testdb, got %s", cfg.DBName)
	}

	if cfg.RedisHost != "redis.example.com" {
		t.Errorf("Expected Redis host redis.example.com, got %s", cfg.RedisHost)
	}

	if len(cfg.KafkaBrokers) != 1 || cfg.KafkaBrokers[0] != "kafka.example.com:9092" {
		t.Errorf("Expected Kafka broker kafka.example.com:9092, got %v", cfg.KafkaBrokers)
	}
}

func TestConfigTimeouts(t *testing.T) {
	// Clear all env vars
	os.Unsetenv("JWT_EXPIRES_IN")
	os.Unsetenv("REFRESH_TOKEN_TTL")
	os.Unsetenv("AGENT_TIMEOUT")

	cfg := Load()

	// Check default durations
	if cfg.JWTExpiresIn != 24*time.Hour {
		t.Errorf("Expected JWT expiration 24h, got %v", cfg.JWTExpiresIn)
	}

	if cfg.RefreshTokenTTL != 168*time.Hour {
		t.Errorf("Expected refresh token TTL 168h, got %v", cfg.RefreshTokenTTL)
	}

	if cfg.AgentTimeout != 5*time.Minute {
		t.Errorf("Expected agent timeout 5m, got %v", cfg.AgentTimeout)
	}
}
