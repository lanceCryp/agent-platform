package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// TaskHandler handles task-related requests
type TaskHandler struct {
	db    *pgxpool.Pool
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(db *pgxpool.Pool) *TaskHandler {
	return &TaskHandler{
		db: db,
	}
}

// CreateTaskRequest represents task creation input
type CreateTaskRequest struct {
	AgentID    string                 `json:"agent_id" binding:"required"`
	Prompt     string                 `json:"prompt" binding:"required,min=1,max=10000"`
	Priority   int                    `json:"priority"`
	MaxRetries int                    `json:"max_retries"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// CreateTaskResponse represents task creation response
type CreateTaskResponse struct {
	TaskID           string  `json:"task_id"`
	Status           string  `json:"status"`
	EstimatedCost    float64 `json:"estimated_cost"`
	EstimatedSeconds int     `json:"estimated_duration"`
}

// ListTasks returns paginated list of tasks
func (h *TaskHandler) ListTasks(c *gin.Context) {
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT t.id, t.agent_id, a.name as agent_name, t.prompt, t.result,
		       t.status, t.priority, t.cost, t.tokens_used, t.duration_seconds,
		       t.error_message, t.retry_count, t.created_at, t.started_at, t.completed_at
		FROM tasks t
		JOIN agents a ON t.agent_id = a.id
		WHERE t.user_id = $1
	`
	countQuery := "SELECT COUNT(*) FROM tasks WHERE user_id = $1"
	args := []interface{}{userID}
	argIndex := 2

	if status != "" {
		query += " AND t.status = $" + strconv.Itoa(argIndex)
		countQuery += " AND status = $" + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	query += " ORDER BY t.created_at DESC LIMIT $" + strconv.Itoa(argIndex)
	args = append(args, limit)
	argIndex++
	query += " OFFSET $" + strconv.Itoa(argIndex)
	args = append(args, offset)

	// Get total count
	var total int
	err := h.db.QueryRow(ctx, countQuery, args[:1]...).Scan(&total)
	if err != nil {
		logrus.Errorf("Failed to count tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal error"})
		return
	}

	// Get tasks
	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("Failed to query tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal error"})
		return
	}
	defer rows.Close()

	tasks := []map[string]interface{}{}
	for rows.Next() {
		var id, agentID, agentName, prompt *string
		var result, errorMessage *string
		var status string
		var priority, tokensUsed, durationSeconds, retryCount int
		var cost float64
		var createdAt, startedAt, completedAt *time.Time

		err := rows.Scan(
			&id, &agentID, &agentName, &prompt, &result,
			&status, &priority, &cost, &tokensUsed, &durationSeconds,
			&errorMessage, &retryCount, &createdAt, &startedAt, &completedAt,
		)
		if err != nil {
			logrus.Errorf("Failed to scan task: %v", err)
			continue
		}

		task := map[string]interface{}{
			"id":                id,
			"agent_id":          agentID,
			"agent_name":        agentName,
			"status":            status,
			"priority":          priority,
			"cost":              cost,
			"tokens_used":       tokensUsed,
			"duration_seconds":  durationSeconds,
			"retry_count":       retryCount,
			"created_at":        createdAt,
		}

		if prompt != nil {
			task["prompt"] = *prompt
		}
		if result != nil {
			task["result"] = *result
		}
		if errorMessage != nil {
			task["error_message"] = *errorMessage
		}
		if startedAt != nil {
			task["started_at"] = *startedAt
		}
		if completedAt != nil {
			task["completed_at"] = *completedAt
		}

		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
			"tasks": tasks,
		},
	})
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request: " + err.Error()})
		return
	}

	// Set defaults
	if req.Priority == 0 {
		req.Priority = 5
	}
	if req.MaxRetries == 0 {
		req.MaxRetries = 3
	}

	// Get agent info
	var agentUUID uuid.UUID
	var agentName, runtimeType string
	var pricePerRequest float64
	var avgDuration int

	err := h.db.QueryRow(ctx, `
		SELECT id, name, runtime_type, price_per_request, avg_duration_seconds
		FROM agents WHERE agent_id = $1 AND is_active = true
	`, req.AgentID).Scan(&agentUUID, &agentName, &runtimeType, &pricePerRequest, &avgDuration)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Agent not found"})
		return
	}

	// Check user balance
	var balance float64
	err = h.db.QueryRow(ctx, "SELECT COALESCE(balance, 0) FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to check balance"})
		return
	}

	if balance < pricePerRequest {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"success": false,
			"error":   "Insufficient balance. Please recharge.",
			"required": pricePerRequest,
			"current":  balance,
		})
		return
	}

	// Serialize context
	var contextJSON []byte
	if req.Context != nil {
		contextJSON, _ = json.Marshal(req.Context)
	}

	// Create task
	var taskID uuid.UUID
	err = h.db.QueryRow(ctx, `
		INSERT INTO tasks (user_id, agent_id, prompt, status, priority, cost, max_retries, context)
		VALUES ($1, $2, $3, 'pending', $4, $5, $6, $7)
		RETURNING id
	`, userID, agentUUID, req.Prompt, req.Priority, pricePerRequest, req.MaxRetries, contextJSON).Scan(&taskID)

	if err != nil {
		logrus.Errorf("Failed to create task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to create task"})
		return
	}

	// Publish to message queue (Kafka/Redis)
	// In production, this would publish to Kafka
	h.publishTaskCreated(ctx, taskID.String(), agentUUID.String(), runtimeType, req.Prompt)

	logrus.Infof("Task created: %s for user: %s, agent: %s", taskID, userID, req.AgentID)

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data": CreateTaskResponse{
			TaskID:           taskID.String(),
			Status:           "pending",
			EstimatedCost:    pricePerRequest,
			EstimatedSeconds: avgDuration,
		},
	})
}

// GetTask returns task details
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	var id, agentID string
	var agentName string
	var prompt, result, errorMessage *string
	var status string
	var priority, tokensUsed, durationSeconds, retryCount int
	var cost float64
	var createdAt, startedAt, completedAt *time.Time

	err := h.db.QueryRow(ctx, `
		SELECT t.id, t.agent_id, a.name, t.prompt, t.result,
		       t.status, t.priority, t.cost, t.tokens_used, t.duration_seconds,
		       t.error_message, t.retry_count, t.created_at, t.started_at, t.completed_at
		FROM tasks t
		JOIN agents a ON t.agent_id = a.id
		WHERE t.id = $1 AND t.user_id = $2
	`, taskID, userID).Scan(
		&id, &agentID, &agentName, &prompt, &result,
		&status, &priority, &cost, &tokensUsed, &durationSeconds,
		&errorMessage, &retryCount, &createdAt, &startedAt, &completedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Task not found"})
		return
	}

	task := map[string]interface{}{
		"id":               id,
		"agent_id":         agentID,
		"agent_name":       agentName,
		"status":           status,
		"priority":         priority,
		"cost":             cost,
		"tokens_used":      tokensUsed,
		"duration_seconds": durationSeconds,
		"retry_count":      retryCount,
		"created_at":       createdAt,
	}

	if prompt != nil {
		task["prompt"] = *prompt
	}
	if result != nil {
		task["result"] = *result
	}
	if errorMessage != nil {
		task["error_message"] = *errorMessage
	}
	if startedAt != nil {
		task["started_at"] = *startedAt
	}
	if completedAt != nil {
		task["completed_at"] = *completedAt
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    task,
	})
}

// CancelTask cancels a pending task
func (h *TaskHandler) CancelTask(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	// Check task ownership and status
	var currentStatus string
	err := h.db.QueryRow(ctx, `
		SELECT status FROM tasks WHERE id = $1 AND user_id = $2
	`, taskID, userID).Scan(&currentStatus)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Task not found"})
		return
	}

	if currentStatus != "pending" && currentStatus != "processing" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Task cannot be cancelled in current status: " + currentStatus,
		})
		return
	}

	// Update status
	_, err = h.db.Exec(ctx, `
		UPDATE tasks SET status = 'cancelled', completed_at = NOW()
		WHERE id = $1
	`, taskID)

	if err != nil {
		logrus.Errorf("Failed to cancel task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to cancel task"})
		return
	}

	// Refund balance (if already deducted)
	var cost float64
	h.db.QueryRow(ctx, "SELECT cost FROM tasks WHERE id = $1", taskID).Scan(&cost)
	if cost > 0 {
		h.db.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", cost, userID)
		
		// Record refund transaction
		h.db.Exec(ctx, `
			INSERT INTO transactions (user_id, type, amount, balance_before, balance_after, description, reference_id, reference_type)
			SELECT $1, 'refund', $2, balance - $2, balance, 'Task cancelled', $3, 'task'
			FROM users WHERE id = $1
		`, userID, cost, taskID)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    map[string]string{"message": "Task cancelled successfully"},
	})
}

// RetryTask retries a failed task
func (h *TaskHandler) RetryTask(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	// Get original task
	var status, agentID, prompt string
	var maxRetries, retryCount int
	err := h.db.QueryRow(ctx, `
		SELECT status, agent_id, prompt, max_retries, retry_count
		FROM tasks WHERE id = $1 AND user_id = $2
	`, taskID, userID).Scan(&status, &agentID, &prompt, &maxRetries, &retryCount)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Task not found"})
		return
	}

	if status != "failed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Only failed tasks can be retried",
		})
		return
	}

	if retryCount >= maxRetries {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Maximum retry attempts reached",
		})
		return
	}

	// Create new task with reference to parent
	var newTaskID uuid.UUID
	err = h.db.QueryRow(ctx, `
		INSERT INTO tasks (user_id, agent_id, prompt, status, priority, cost, max_retries, parent_task_id)
		SELECT user_id, agent_id, prompt, 'pending', priority, cost, max_retries, $1
		FROM tasks WHERE id = $2
		RETURNING id
	`, taskID, taskID).Scan(&newTaskID)

	if err != nil {
		logrus.Errorf("Failed to retry task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to retry task"})
		return
	}

	// Update original task retry count
	h.db.Exec(ctx, "UPDATE tasks SET retry_count = retry_count + 1 WHERE id = $1", taskID)

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"task_id":  newTaskID.String(),
			"status":   "pending",
			"message":  "Task queued for retry",
		},
	})
}

// Helper function to publish task creation event
func (h *TaskHandler) publishTaskCreated(ctx context.Context, taskID, agentID, runtimeType, prompt string) {
	// In production, this would publish to Kafka
	// For now, we just log it
	logrus.Infof("Publishing task created event: taskID=%s, agentID=%s, runtime=%s", taskID, agentID, runtimeType)
	
	// Example Kafka message:
	// message := map[string]interface{}{
	// 	"task_id":      taskID,
	// 	"agent_id":     agentID,
	// 	"runtime_type": runtimeType,
	// 	"prompt":       prompt,
	// 	"timestamp":    time.Now(),
	// }
}
