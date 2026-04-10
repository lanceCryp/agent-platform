package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, connString string) (*PostgresDB, error) {
	logrus.Info("Connecting to PostgreSQL...")

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Connection pool settings
	config.MaxConns = 50
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("PostgreSQL connection established successfully")
	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		logrus.Info("PostgreSQL connection closed")
	}
}

func (db *PostgresDB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// HealthCheck returns database health status
func (db *PostgresDB) HealthCheck(ctx context.Context) map[string]interface{} {
	status := "healthy"
	latency := time.Duration(0)

	start := time.Now()
	if err := db.Pool.Ping(ctx); err != nil {
		status = "unhealthy"
	}
	latency = time.Since(start)

	return map[string]interface{}{
		"status":        status,
		"latency_ms":    latency.Milliseconds(),
		"max_conns":     db.Pool.Config().MaxConns,
		"total_conns":   db.Pool.Stat().TotalConns(),
		"idle_conns":    db.Pool.Stat().IdleConns(),
		"working_conns": db.Pool.Stat().AcquiredConns(),
	}
}
