package auth

import (
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	// Test hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should not equal original password")
	}

	// Test that same password produces different hashes (due to salt)
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password second time: %v", err)
	}

	if hash == hash2 {
		t.Error("Same password should produce different hashes due to salt")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	// Hash the password
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	if !CheckPassword(password, hash) {
		t.Error("Correct password should match hash")
	}

	// Test wrong password
	if CheckPassword(wrongPassword, hash) {
		t.Error("Wrong password should not match hash")
	}

	// Test empty password
	if CheckPassword("", hash) {
		t.Error("Empty password should not match hash")
	}
}

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-123", 1*time.Hour, 7*24*time.Hour)

	userID := "user-123"
	email := "test@example.com"
	username := "testuser"
	role := "user"

	token, err := manager.GenerateToken(userID, email, username, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}

	// Token should be a valid JWT format (3 parts separated by dots)
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("Token should have 3 parts, got %d dots", parts)
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-123", 1*time.Hour, 7*24*time.Hour)

	userID := "user-123"
	email := "test@example.com"
	username := "testuser"
	role := "user"

	token, err := manager.GenerateToken(userID, email, username, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID mismatch: got %s, want %s", claims.UserID, userID)
	}

	if claims.Email != email {
		t.Errorf("Email mismatch: got %s, want %s", claims.Email, email)
	}

	if claims.Username != username {
		t.Errorf("Username mismatch: got %s, want %s", claims.Username, username)
	}

	if claims.Role != role {
		t.Errorf("Role mismatch: got %s, want %s", claims.Role, role)
	}
}

func TestJWTManager_ValidateInvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-123", 1*time.Hour, 7*24*time.Hour)

	// Test invalid token
	_, err := manager.ValidateToken("invalid.token.here")
	if err == nil {
		t.Error("Invalid token should return error")
	}

	// Test empty token
	_, err = manager.ValidateToken("")
	if err == nil {
		t.Error("Empty token should return error")
	}
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-123", 1*time.Hour, 7*24*time.Hour)

	userID := "user-123"

	refreshToken, err := manager.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	if refreshToken == "" {
		t.Error("Refresh token should not be empty")
	}

	// Validate refresh token
	claims, err := manager.ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("Failed to validate refresh token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID mismatch: got %s, want %s", claims.UserID, userID)
	}
}

func TestJWTManager_RefreshTokens(t *testing.T) {
	manager := NewJWTManager("test-secret-key-123", 1*time.Hour, 7*24*time.Hour)

	userID := "user-123"

	// Generate initial tokens
	originalToken, refreshToken, err := manager.RefreshTokens("")
	if err != nil {
		t.Fatalf("Failed to generate initial tokens: %v", err)
	}

	// Validate original token
	_, err = manager.ValidateToken(originalToken)
	if err != nil {
		t.Error("Original token should be valid")
	}

	// Generate new tokens using refresh token
	newToken, newRefreshToken, err := manager.RefreshTokens(refreshToken)
	if err != nil {
		t.Fatalf("Failed to refresh tokens: %v", err)
	}

	if newToken == "" || newRefreshToken == "" {
		t.Error("New tokens should not be empty")
	}

	// New token should be different from original
	if newToken == originalToken {
		t.Error("New token should be different from original")
	}
}

func TestJWTManager_WrongSecret(t *testing.T) {
	manager1 := NewJWTManager("secret-key-1", 1*time.Hour, 7*24*time.Hour)
	manager2 := NewJWTManager("secret-key-2", 1*time.Hour, 7*24*time.Hour)

	// Generate token with manager1
	token, err := manager1.GenerateToken("user-123", "test@example.com", "testuser", "user")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with manager2 (different secret)
	_, err = manager2.ValidateToken(token)
	if err == nil {
		t.Error("Token signed with different secret should not validate")
	}
}

func TestJWTManager_Expiration(t *testing.T) {
	// Create manager with very short expiration
	manager := NewJWTManager("test-secret", 1*time.Millisecond, 1*time.Millisecond)

	token, err := manager.GenerateToken("user-123", "test@example.com", "testuser", "user")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Token should be expired
	_, err = manager.ValidateToken(token)
	if err != ErrExpiredToken {
		t.Errorf("Expected ErrExpiredToken, got: %v", err)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "benchmarkpassword123"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(password, hash)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	manager := NewJWTManager("test-secret", 1*time.Hour, 7*24*time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.GenerateToken("user-123", "test@example.com", "testuser", "user")
	}
}
