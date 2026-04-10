package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agenthub/pkg/config"
	"agenthub/pkg/logger"
	"agenthub/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger.Init()
	log.SetOutput(logger.NewLogrusWriter())
	log.Info("Starting AgentHub API Gateway...")

	// Load configuration
	cfg := config.Load()
	log.Infof("Configuration loaded: port=%s, env=%s", cfg.Port, cfg.Env)

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Info("Database connected successfully")

	// Initialize Redis
	redisClient := initRedis(cfg)
	defer redisClient.Close()
	log.Info("Redis connected successfully")

	// Initialize Kafka
	kafkaProducer := initKafkaProducer(cfg)
	defer kafkaProducer.Close()
	log.Info("Kafka producer initialized")

	// Setup Gin router
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit(redisClient))

	// Health check
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck(db))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handleRegister(db))
			auth.POST("/login", handleLogin(db, redisClient))
			auth.POST("/refresh", handleRefreshToken(redisClient))
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.Auth(cfg.JWTSecret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", handleGetProfile(db))
				users.PATCH("/me", handleUpdateProfile(db))
				users.GET("/me/balance", handleGetBalance(db))
			}

			// Agent routes
			agents := protected.Group("/agents")
			{
				agents.GET("", handleListAgents(db, redisClient))
				agents.GET("/categories", handleListCategories(db))
				agents.GET("/:id", handleGetAgent(db, redisClient))
			}

			// Task routes
			tasks := protected.Group("/tasks")
			{
				tasks.POST("", handleCreateTask(db, kafkaProducer))
				tasks.GET("", handleListTasks(db))
				tasks.GET("/:id", handleGetTask(db))
				tasks.POST("/:id/cancel", handleCancelTask(db))
				tasks.POST("/:id/retry", handleRetryTask(db, kafkaProducer))
			}

			// Plan routes
			plans := protected.Group("/plans")
			{
				plans.GET("", handleListPlans(db))
				plans.GET("/:id", handleGetPlan(db))
			}

			// Subscription routes
			subs := protected.Group("/subscriptions")
			{
				subs.GET("", handleListSubscriptions(db))
				subs.POST("", handleCreateSubscription(db, kafkaProducer))
				subs.POST("/:id/cancel", handleCancelSubscription(db))
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				transactions.GET("", handleListTransactions(db))
			}

			// Billing routes
			billing := protected.Group("/billing")
			{
				billing.POST("/recharge", handleRecharge(db, kafkaProducer))
				billing.GET("/history", handleGetBillingHistory(db))
			}
		}
	}

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout:  30 * time.Second,
		IdleTimeout:   60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Infof("API Gateway listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "agenthub-api",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness check endpoint
func readinessCheck(db interface{ Ping(context.Context) error }) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "database connection failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
}

// Placeholder functions for initialization
func initDatabase(cfg *config.Config) (interface{ Close() error; Ping(context.Context) error }, error) {
	// In real implementation, this would be PostgreSQL connection
	return &mockDB{}, nil
}

func initRedis(cfg *config.Config) interface{ Close() error } {
	// In real implementation, this would be Redis connection
	return &mockRedis{}
}

func initKafkaProducer(cfg *config.Config) interface{ Close() error; WriteMessages(context.Context, ...interface{}) error } {
	// In real implementation, this would be Kafka producer
	return &mockKafka{}
}

// Mock implementations for compilation
type mockDB struct{}
func (m *mockDB) Close() error { return nil }
func (m *mockDB) Ping(_ context.Context) error { return nil }

type mockRedis struct{}
func (m *mockRedis) Close() error { return nil }

type mockKafka struct{}
func (m *mockKafka) Close() error { return nil }
func (m *mockKafka) WriteMessages(_ context.Context, _ ...interface{}) error { return nil }

// Placeholder handlers
func handleRegister(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "register endpoint"})
	}
}

func handleLogin(db, redis interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "login endpoint"})
	}
}

func handleRefreshToken(redis interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "refresh endpoint"})
	}
}

func handleGetProfile(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get profile endpoint"})
	}
}

func handleUpdateProfile(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "update profile endpoint"})
	}
}

func handleGetBalance(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get balance endpoint"})
	}
}

func handleListAgents(db, redis interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list agents endpoint"})
	}
}

func handleListCategories(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list categories endpoint"})
	}
}

func handleGetAgent(db, redis interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get agent endpoint"})
	}
}

func handleCreateTask(db interface{}, kafka interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "create task endpoint"})
	}
}

func handleListTasks(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list tasks endpoint"})
	}
}

func handleGetTask(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get task endpoint"})
	}
}

func handleCancelTask(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "cancel task endpoint"})
	}
}

func handleRetryTask(db, kafka interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "retry task endpoint"})
	}
}

func handleListPlans(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list plans endpoint"})
	}
}

func handleGetPlan(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get plan endpoint"})
	}
}

func handleListSubscriptions(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list subscriptions endpoint"})
	}
}

func handleCreateSubscription(db interface{}, kafka interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "create subscription endpoint"})
	}
}

func handleCancelSubscription(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "cancel subscription endpoint"})
	}
}

func handleListTransactions(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list transactions endpoint"})
	}
}

func handleRecharge(db interface{}, kafka interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "recharge endpoint"})
	}
}

func handleGetBillingHistory(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get billing history endpoint"})
	}
}
