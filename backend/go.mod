module agenthub

go 1.22

// Core web framework
github.com/gin-gonic/gin v1.9.1

// Database
github.com/jackc/pgx/v5 v5.5.0
github.com/golang-migrate/migrate/v4 v4.17.0

// Redis cache
github.com/redis/go-redis/v9 v9.4.0

// Kafka message queue
github.com/segmentio/kafka-go v0.4.47

// Authentication
github.com/golang-jwt/jwt/v5 v5.2.0
golang.org/x/crypto v0.18.0

// Configuration
github.com/spf13/viper v1.18.0

// Logging
github.com/sirupsen/logrus v1.9.3
github.com/natefinch/lumberjack v2.1.0

// Observability
go.opentelemetry.io/otel v1.21.0
go.opentelemetry.io/otel/trace v1.21.0
go.opentelemetry.io/otel/sdk v1.21.0
github.com/prometheus/client_golang v1.17.0

// Error handling
github.com/pkg/errors v0.9.1

// Utilities
github.com/google/uuid v1.5.0
github.com/gosimple/slug v1.13.1
github.com/robfig/cron/v3 v3.0.1

// Validation
github.com/go-playground/validator/v10 v10.16.0

// Testing
github.com/stretchr/testify v1.8.4
github.com/kinbiko/jsonassert v1.5.1
github.com/DATA-DOG/go-sqlmock v1.5.2
