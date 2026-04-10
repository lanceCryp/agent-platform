package handler

import (
	"context"
	"net/http"
	"time"

	"agenthub/pkg/auth"
	"agenthub/pkg/cache"
	"agenthub/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	db         *pgxpool.Pool
	cache      *cache.RedisCache
	jwtManager *auth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *pgxpool.Pool, cache *cache.RedisCache, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		db:         db,
		cache:      cache,
		jwtManager: jwtManager,
	}
}

// RegisterRequest represents registration input
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

// LoginRequest represents login input
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents auth response
type AuthResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// TokenResponse represents token response
type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *UserResponse `json:"user"`
}

// UserResponse represents user response
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Avatar   string    `json:"avatar_url,omitempty"`
	Balance  float64   `json:"balance"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Check if email already exists
	var exists bool
	err := h.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		logrus.Errorf("Database error checking email: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, AuthResponse{
			Success: false,
			Error:   "Email already registered",
		})
		return
	}

	// Check if username already exists
	err = h.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
	if err != nil {
		logrus.Errorf("Database error checking username: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, AuthResponse{
			Success: false,
			Error:   "Username already taken",
		})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logrus.Errorf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	// Create user
	var userID uuid.UUID
	var user UserResponse
	err = h.db.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash, role, balance)
		VALUES ($1, $2, $3, 'user', 0.00)
		RETURNING id, email, username, role, balance
	`, req.Email, req.Username, hashedPassword).Scan(
		&user.ID, &user.Email, &user.Username, &user.Role, &user.Balance,
	)

	if err != nil {
		logrus.Errorf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Failed to create user",
		})
		return
	}

	// Generate tokens
	token, err := h.jwtManager.GenerateToken(user.ID.String(), user.Email, user.Username, user.Role)
	if err != nil {
		logrus.Errorf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		logrus.Errorf("Failed to generate refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	// Store session in cache
	sessionKey := "session:" + user.ID.String()
	h.cache.Set(ctx, sessionKey, refreshToken, 7*24*time.Hour)

	logrus.Infof("User registered successfully: %s", user.Email)

	c.JSON(http.StatusCreated, AuthResponse{
		Success: true,
		Data: TokenResponse{
			Token:        token,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(h.jwtManager.GetTokenExpiration().Seconds()),
			User:         &user,
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Get user by email
	var user UserResponse
	var passwordHash string
	err := h.db.QueryRow(ctx, `
		SELECT id, email, username, password_hash, role, COALESCE(balance, 0)
		FROM users WHERE email = $1 AND is_active = true
	`, req.Email).Scan(&user.ID, &user.Email, &user.Username, &passwordHash, &user.Role, &user.Balance)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, AuthResponse{
				Success: false,
				Error:   "Invalid email or password",
			})
			return
		}
		logrus.Errorf("Database error during login: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	// Verify password
	if !auth.CheckPassword(req.Password, passwordHash) {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	// Update last login
	h.db.Exec(ctx, `
		UPDATE users SET last_login_at = $1, last_login_ip = $2
		WHERE id = $3
	`, time.Now(), c.ClientIP(), user.ID)

	// Generate tokens
	token, err := h.jwtManager.GenerateToken(user.ID.String(), user.Email, user.Username, user.Role)
	if err != nil {
		logrus.Errorf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		logrus.Errorf("Failed to generate refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Internal server error",
		})
		return
	}

	// Store session
	sessionKey := "session:" + user.ID.String()
	h.cache.Set(ctx, sessionKey, refreshToken, 7*24*time.Hour)

	logrus.Infof("User logged in: %s", user.Email)

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Data: TokenResponse{
			Token:        token,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(h.jwtManager.GetTokenExpiration().Seconds()),
			User:         &user,
		},
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Error:   "Invalid request",
		})
		return
	}

	// Validate refresh token
	claims, err := h.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Error:   "Invalid or expired refresh token",
		})
		return
	}

	ctx := c.Request.Context()

	// Get user
	var user UserResponse
	err = h.db.QueryRow(ctx, `
		SELECT id, email, username, role, COALESCE(balance, 0)
		FROM users WHERE id = $1 AND is_active = true
	`, claims.UserID).Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.Balance)

	if err != nil {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Generate new tokens
	token, err := h.jwtManager.GenerateToken(user.ID.String(), user.Email, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	newRefreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Failed to generate refresh token",
		})
		return
	}

	// Update session
	sessionKey := "session:" + user.ID.String()
	h.cache.Set(ctx, sessionKey, newRefreshToken, 7*24*time.Hour)

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Data: TokenResponse{
			Token:        token,
			RefreshToken: newRefreshToken,
			ExpiresIn:    int64(h.jwtManager.GetTokenExpiration().Seconds()),
			User:         &user,
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()

	// Remove session from cache
	sessionKey := "session:" + userID.(string)
	h.cache.Delete(ctx, sessionKey)

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Data:    map[string]string{"message": "Logged out successfully"},
	})
}

// GetProfile returns current user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	ctx := c.Request.Context()

	var user UserResponse
	err := h.db.QueryRow(ctx, `
		SELECT id, email, username, role, COALESCE(avatar_url, ''), COALESCE(balance, 0)
		FROM users WHERE id = $1
	`, userID).Scan(&user.ID, &user.Email, &user.Username, &user.Role, &user.Avatar, &user.Balance)

	if err != nil {
		c.JSON(http.StatusNotFound, AuthResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Data:    user,
	})
}

// UpdateProfileRequest represents profile update input
type UpdateProfileRequest struct {
	Username  *string `json:"username,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UpdateProfile handles profile update
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Error:   "Invalid request",
		})
		return
	}

	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	// Build update query dynamically
	query := "UPDATE users SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if req.Username != nil {
		// Check if username is taken
		var exists bool
		h.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id != $2)", *req.Username, userID).Scan(&exists)
		if exists {
			c.JSON(http.StatusConflict, AuthResponse{
				Success: false,
				Error:   "Username already taken",
			})
			return
		}
		query += ", username = $" + string(rune('0'+argIndex))
		args = append(args, *req.Username)
		argIndex++
	}

	if req.Phone != nil {
		query += ", phone = $" + string(rune('0'+argIndex))
		args = append(args, *req.Phone)
		argIndex++
	}

	if req.AvatarURL != nil {
		query += ", avatar_url = $" + string(rune('0'+argIndex))
		args = append(args, *req.AvatarURL)
		argIndex++
	}

	query += " WHERE id = $" + string(rune('0'+argIndex))
	args = append(args, userID)

	_, err := h.db.Exec(ctx, query, args...)
	if err != nil {
		logrus.Errorf("Failed to update profile: %v", err)
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Data:    map[string]string{"message": "Profile updated successfully"},
	})
}

// GetBalance returns user balance
func (h *AuthHandler) GetBalance(c *gin.Context) {
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	var balance float64
	var subscription interface{}

	err := h.db.QueryRow(ctx, `
		SELECT COALESCE(balance, 0) FROM users WHERE id = $1
	`, userID).Scan(&balance)

	if err != nil {
		c.JSON(http.StatusNotFound, AuthResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Get active subscription
	var sub struct {
		PlanType string    `json:"plan_type"`
		EndDate  time.Time `json:"end_date"`
	}
	err = h.db.QueryRow(ctx, `
		SELECT plan_type, end_date FROM subscriptions
		WHERE user_id = $1 AND status = 'active' AND end_date > NOW()
		ORDER BY created_at DESC LIMIT 1
	`, userID).Scan(&sub.PlanType, &sub.EndDate)

	if err == nil {
		subscription = sub
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Data: map[string]interface{}{
			"balance":      balance,
			"subscription": subscription,
		},
	})
}

// ChangePasswordRequest represents password change input
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Error:   "Invalid request",
		})
		return
	}

	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	// Get current password hash
	var passwordHash string
	err := h.db.QueryRow(ctx, "SELECT password_hash FROM users WHERE id = $1", userID).Scan(&passwordHash)
	if err != nil {
		c.JSON(http.StatusNotFound, AuthResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	// Verify old password
	if !auth.CheckPassword(req.OldPassword, passwordHash) {
		c.JSON(http.StatusUnauthorized, AuthResponse{
			Success: false,
			Error:   "Current password is incorrect",
		})
		return
	}

	// Hash new password
	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Failed to process password",
		})
		return
	}

	// Update password
	_, err = h.db.Exec(ctx, "UPDATE users SET password_hash = $1 WHERE id = $2", newHash, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AuthResponse{
			Success: false,
			Error:   "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Data:    map[string]string{"message": "Password changed successfully"},
	})
}
