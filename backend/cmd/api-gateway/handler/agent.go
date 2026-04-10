package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"agenthub/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// AgentHandler handles agent-related requests
type AgentHandler struct {
	db    *pgxpool.Pool
	cache *cache.RedisCache
}

// NewAgentHandler creates a new agent handler
func NewAgentHandler(db *pgxpool.Pool, cache *cache.RedisCache) *AgentHandler {
	return &AgentHandler{
		db:    db,
		cache: cache,
	}
}

// ListAgents returns paginated list of agents
func (h *AgentHandler) ListAgents(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	category := c.Query("category")
	tierStr := c.Query("tier")
	tier := 0
	if tierStr != "" {
		tier, _ = strconv.Atoi(tierStr)
	}
	search := c.Query("search")
	runtime := c.Query("runtime")
	sortBy := c.DefaultQuery("sort", "popular")

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
		SELECT id, agent_id, name, name_en, description, category, tags, tier,
		       runtime_type, price_per_request, avg_duration_seconds,
		       success_rate, rating, total_tasks, is_active
		FROM agents
		WHERE is_active = true
	`
	countQuery := "SELECT COUNT(*) FROM agents WHERE is_active = true"
	args := []interface{}{}
	argIndex := 1

	if category != "" {
		query += " AND category = $" + strconv.Itoa(argIndex)
		countQuery += " AND category = $" + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	if tier > 0 {
		query += " AND tier <= $" + strconv.Itoa(argIndex)
		countQuery += " AND tier <= $" + strconv.Itoa(argIndex)
		args = append(args, tier)
		argIndex++
	}

	if search != "" {
		query += " AND (name ILIKE $" + strconv.Itoa(argIndex) + " OR description ILIKE $" + strconv.Itoa(argIndex) + ")"
		countQuery += " AND (name ILIKE $" + strconv.Itoa(argIndex) + " OR description ILIKE $" + strconv.Itoa(argIndex) + ")"
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if runtime != "" {
		query += " AND runtime_type = $" + strconv.Itoa(argIndex)
		countQuery += " AND runtime_type = $" + strconv.Itoa(argIndex)
		args = append(args, runtime)
		argIndex++
	}

	// Add sorting
	switch sortBy {
	case "rating":
		query += " ORDER BY rating DESC"
	case "price-low":
		query += " ORDER BY price_per_request ASC"
	case "price-high":
		query += " ORDER BY price_per_request DESC"
	default: // popular
		query += " ORDER BY total_tasks DESC"
	}

	// Add pagination
	query += " LIMIT $" + strconv.Itoa(argIndex)
	args = append(args, limit)
	argIndex++
	query += " OFFSET $" + strconv.Itoa(argIndex)
	args = append(args, offset)

	// Get total count
	var total int
	err := h.db.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		logrus.Errorf("Failed to count agents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal error"})
		return
	}

	// Get agents
	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("Failed to query agents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal error"})
		return
	}
	defer rows.Close()

	agents := []map[string]interface{}{}
	for rows.Next() {
		var id string
		var agentID, name, nameEn, description, cat, runtimeType string
		var tags []string
		var tier, avgDuration, totalTasks int
		var pricePerRequest, successRate, rating float64
		var isActive bool

		err := rows.Scan(
			&id, &agentID, &name, &nameEn, &description, &cat, &tags, &tier,
			&runtimeType, &pricePerRequest, &avgDuration, &successRate, &rating,
			&totalTasks, &isActive,
		)
		if err != nil {
			logrus.Errorf("Failed to scan agent: %v", err)
			continue
		}

		agents = append(agents, map[string]interface{}{
			"id":                   id,
			"agent_id":             agentID,
			"name":                 name,
			"name_en":              nameEn,
			"description":          description,
			"category":             cat,
			"tags":                 tags,
			"tier":                 tier,
			"runtime_type":          runtimeType,
			"price_per_request":     pricePerRequest,
			"avg_duration_seconds": avgDuration,
			"success_rate":          successRate,
			"rating":               rating,
			"total_tasks":          totalTasks,
			"is_active":             isActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"total":    total,
			"page":     page,
			"limit":    limit,
			"agents":   agents,
		},
	})
}

// GetAgent returns agent details
func (h *AgentHandler) GetAgent(c *gin.Context) {
	agentID := c.Param("id")
	ctx := c.Request.Context()

	// Try cache first
	cacheKey := "agent:" + agentID
	if cached, err := h.cache.Get(ctx, cacheKey); err == nil {
		c.Header("X-Cache", "HIT")
		c.Data(http.StatusOK, "application/json", []byte(cached))
		return
	}

	// Query database
	var id, agentId, name, nameEn, description, cat, runtimeType string
	var tags []string
	var tier, avgDuration, totalTasks, totalSuccess, totalFailed int
	var pricePerRequest, successRate, rating float64
	var isActive, isFeatured bool
	var config, metadata []byte
	var inputExample, outputExample *string

	err := h.db.QueryRow(ctx, `
		SELECT id, agent_id, name, name_en, description, category, tags, tier,
		       runtime_type, price_per_request, avg_duration_seconds,
		       success_rate, rating, total_tasks, total_success, total_failed,
		       is_active, is_featured, config, metadata, input_example, output_example
		FROM agents
		WHERE agent_id = $1 AND is_active = true
	`, agentID).Scan(
		&id, &agentId, &name, &nameEn, &description, &cat, &tags, &tier,
		&runtimeType, &pricePerRequest, &avgDuration, &successRate, &rating,
		&totalTasks, &totalSuccess, &totalFailed, &isActive, &isFeatured,
		&config, &metadata, &inputExample, &outputExample,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Agent not found",
		})
		return
	}

	// Get category info
	var categoryName, categoryIcon string
	h.db.QueryRow(ctx, "SELECT name, icon FROM categories WHERE slug = $1", cat).Scan(&categoryName, &categoryIcon)

	result := map[string]interface{}{
		"id":                   id,
		"agent_id":             agentId,
		"name":                 name,
		"name_en":              nameEn,
		"description":          description,
		"category":             cat,
		"category_name":        categoryName,
		"category_icon":        categoryIcon,
		"tags":                 tags,
		"tier":                 tier,
		"runtime_type":          runtimeType,
		"price_per_request":     pricePerRequest,
		"avg_duration_seconds": avgDuration,
		"success_rate":          successRate,
		"rating":               rating,
		"total_tasks":          totalTasks,
		"total_success":         totalSuccess,
		"total_failed":          totalFailed,
		"is_active":             isActive,
		"is_featured":           isFeatured,
		"input_example":         inputExample,
		"output_example":       outputExample,
	}

	// Cache for 5 minutes
	// h.cache.Set(ctx, cacheKey, result, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ListCategories returns all categories
func (h *AgentHandler) ListCategories(c *gin.Context) {
	ctx := c.Request.Context()

	// Query categories with agent count
	rows, err := h.db.Query(ctx, `
		SELECT c.id, c.name, c.slug, c.icon, c.description, c.sort_order,
		       COUNT(a.id) as agent_count
		FROM categories c
		LEFT JOIN agents a ON a.category = c.slug AND a.is_active = true
		WHERE c.is_active = true
		GROUP BY c.id, c.name, c.slug, c.icon, c.description, c.sort_order
		ORDER BY c.sort_order ASC
	`)
	if err != nil {
		logrus.Errorf("Failed to query categories: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal error"})
		return
	}
	defer rows.Close()

	categories := []map[string]interface{}{}
	for rows.Next() {
		var id, name, slug, icon, description string
		var sortOrder, agentCount int

		err := rows.Scan(&id, &name, &slug, &icon, &description, &sortOrder, &agentCount)
		if err != nil {
			logrus.Errorf("Failed to scan category: %v", err)
			continue
		}

		categories = append(categories, map[string]interface{}{
			"id":           id,
			"name":         name,
			"slug":         slug,
			"icon":         icon,
			"description":  description,
			"sort_order":   sortOrder,
			"agent_count":  agentCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       categories,
	})
}

// GetFeaturedAgents returns featured agents
func (h *AgentHandler) GetFeaturedAgents(c *gin.Context) {
	ctx := c.Request.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, agent_id, name, name_en, description, category, tags, tier,
		       runtime_type, price_per_request, rating, total_tasks
		FROM agents
		WHERE is_active = true AND is_featured = true
		ORDER BY rating DESC, total_tasks DESC
		LIMIT 8
	`)
	if err != nil {
		logrus.Errorf("Failed to query featured agents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal error"})
		return
	}
	defer rows.Close()

	agents := []map[string]interface{}{}
	for rows.Next() {
		var id, agentID, name, nameEn, description, cat, runtimeType string
		var tags []string
		var tier, totalTasks int
		var pricePerRequest, rating float64

		err := rows.Scan(
			&id, &agentID, &name, &nameEn, &description, &cat, &tags, &tier,
			&runtimeType, &pricePerRequest, &rating, &totalTasks,
		)
		if err != nil {
			continue
		}

		agents = append(agents, map[string]interface{}{
			"id":               id,
			"agent_id":         agentID,
			"name":             name,
			"name_en":          nameEn,
			"description":      description,
			"category":         cat,
			"tags":             tags,
			"tier":             tier,
			"runtime_type":     runtimeType,
			"price_per_request": pricePerRequest,
			"rating":           rating,
			"total_tasks":      totalTasks,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    agents,
	})
}
