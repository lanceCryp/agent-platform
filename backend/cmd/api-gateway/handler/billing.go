package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// BillingHandler handles billing and subscription requests
type BillingHandler struct {
	db *pgxpool.Pool
}

// NewBillingHandler creates a new billing handler
func NewBillingHandler(db *pgxpool.Pool) *BillingHandler {
	return &BillingHandler{db: db}
}

// ListPlans returns all available plans
func (h *BillingHandler) ListPlans(c *gin.Context) {
	ctx := c.Request.Context()

	rows, err := h.db.Query(ctx, `
		SELECT id, name, type, description, price_monthly, price_yearly,
		       discount_percentage, features, task_limit, agent_tier_limit,
		       api_access, priority_support, custom_agents, max_concurrent_tasks,
		       is_featured, sort_order
		FROM plans
		WHERE is_active = true
		ORDER BY sort_order ASC
	`)
	if err != nil {
		logrus.Errorf("Failed to query plans: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal error"})
		return
	}
	defer rows.Close()

	plans := []map[string]interface{}{}
	for rows.Next() {
		var id, name, planType, description string
		var priceMonthly, priceYearly, discountPercentage float64
		var features []byte
		var taskLimit, agentTierLimit, maxConcurrentTasks *int
		var apiAccess, prioritySupport, customAgents, isFeatured bool

		err := rows.Scan(
			&id, &name, &planType, &description, &priceMonthly, &priceYearly,
			&discountPercentage, &features, &taskLimit, &agentTierLimit,
			&apiAccess, &prioritySupport, &customAgents, &maxConcurrentTasks,
			&isFeatured, nil,
		)
		if err != nil {
			logrus.Errorf("Failed to scan plan: %v", err)
			continue
		}

		plan := map[string]interface{}{
			"id":                    id,
			"name":                  name,
			"type":                  planType,
			"description":           description,
			"price_monthly":         priceMonthly,
			"price_yearly":          priceYearly,
			"discount_percentage":   discountPercentage,
			"api_access":            apiAccess,
			"priority_support":      prioritySupport,
			"custom_agents":         customAgents,
			"is_featured":            isFeatured,
		}

		// Parse features JSON
		// var featuresMap map[string]interface{}
		// json.Unmarshal(features, &featuresMap)
		plan["features"] = string(features)
		if taskLimit != nil {
			plan["task_limit"] = *taskLimit
		} else {
			plan["task_limit"] = nil // unlimited
		}
		if agentTierLimit != nil {
			plan["agent_tier_limit"] = *agentTierLimit
		}
		if maxConcurrentTasks != nil {
			plan["max_concurrent_tasks"] = *maxConcurrentTasks
		}

		plans = append(plans, plan)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plans,
	})
}

// GetPlan returns plan details
func (h *BillingHandler) GetPlan(c *gin.Context) {
	planID := c.Param("id")
	ctx := c.Request.Context()

	var id, name, planType, description string
	var priceMonthly, priceYearly float64
	var features []byte

	err := h.db.QueryRow(ctx, `
		SELECT id, name, type, description, price_monthly, price_yearly, features
		FROM plans WHERE id = $1 AND is_active = true
	`, planID).Scan(&id, &name, &planType, &description, &priceMonthly, &priceYearly, &features)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"id":              id,
			"name":            name,
			"type":            planType,
			"description":     description,
			"price_monthly":   priceMonthly,
			"price_yearly":    priceYearly,
			"features":        string(features),
		},
	})
}

// CreateSubscriptionRequest represents subscription creation input
type CreateSubscriptionRequest struct {
	PlanID       string `json:"plan_id" binding:"required"`
	BillingCycle string `json:"billing_cycle" binding:"required,oneof=monthly yearly"`
}

// CreateSubscription creates a new subscription
func (h *BillingHandler) CreateSubscription(c *gin.Context) {
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}

	// Get plan details
	var planType string
	var priceMonthly, priceYearly float64
	var agentTierLimit int
	err := h.db.QueryRow(ctx, `
		SELECT type, price_monthly, price_yearly, agent_tier_limit
		FROM plans WHERE id = $1 AND is_active = true
	`, req.PlanID).Scan(&planType, &priceMonthly, &priceYearly, &agentTierLimit)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Plan not found"})
		return
	}

	// Calculate price
	price := priceMonthly
	if req.BillingCycle == "yearly" {
		price = priceYearly
	}

	// Check if user already has active subscription of same type
	var existingSubID *uuid.UUID
	h.db.QueryRow(ctx, `
		SELECT id FROM subscriptions
		WHERE user_id = $1 AND plan_type = $2 AND status = 'active'
	`, userID, planType).Scan(&existingSubID)

	if existingSubID != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "You already have an active subscription of this type",
		})
		return
	}

	// Calculate dates
	now := time.Now()
	var endDate time.Time
	if req.BillingCycle == "monthly" {
		endDate = now.AddDate(0, 1, 0)
	} else {
		endDate = now.AddDate(1, 0, 0)
	}

	// Create subscription
	var subscriptionID uuid.UUID
	err = h.db.QueryRow(ctx, `
		INSERT INTO subscriptions (user_id, plan_id, plan_type, status, start_date, end_date, price, billing_cycle)
		VALUES ($1, $2, $3, 'pending', $4, $5, $6, $7)
		RETURNING id
	`, userID, req.PlanID, planType, now, endDate, price, req.BillingCycle).Scan(&subscriptionID)

	if err != nil {
		logrus.Errorf("Failed to create subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to create subscription"})
		return
	}

	// In production, redirect to payment gateway
	// For now, we simulate successful payment
	paymentURL := "/api/v1/subscriptions/" + subscriptionID.String() + "/payment?gateway=mock"

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"subscription_id": subscriptionID.String(),
			"payment_url":      paymentURL,
			"amount":           price,
			"billing_cycle":    req.BillingCycle,
		},
	})
}

// CancelSubscription cancels an active subscription
func (h *BillingHandler) CancelSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	// Check ownership and status
	var currentStatus string
	err := h.db.QueryRow(ctx, `
		SELECT status FROM subscriptions WHERE id = $1 AND user_id = $2
	`, subscriptionID, userID).Scan(&currentStatus)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Subscription not found"})
		return
	}

	if currentStatus != "active" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Subscription is not active",
		})
		return
	}

	// Update subscription
	_, err = h.db.Exec(ctx, `
		UPDATE subscriptions SET status = 'cancelled', cancelled_at = NOW()
		WHERE id = $1
	`, subscriptionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to cancel subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    map[string]string{"message": "Subscription cancelled successfully"},
	})
}

// ListSubscriptions returns user's subscriptions
func (h *BillingHandler) ListSubscriptions(c *gin.Context) {
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	rows, err := h.db.Query(ctx, `
		SELECT s.id, s.plan_id, p.name, s.plan_type, s.status,
		       s.start_date, s.end_date, s.price, s.billing_cycle, s.auto_renew
		FROM subscriptions s
		JOIN plans p ON s.plan_id = p.id
		WHERE s.user_id = $1
		ORDER BY s.created_at DESC
	`, userID)

	if err != nil {
		logrus.Errorf("Failed to query subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal error"})
		return
	}
	defer rows.Close()

	subscriptions := []map[string]interface{}{}
	for rows.Next() {
		var id, planID, planName, planType, status, billingCycle string
		var startDate, endDate time.Time
		var price float64
		var autoRenew bool

		err := rows.Scan(&id, &planID, &planName, &planType, &status, &startDate, &endDate, &price, &billingCycle, &autoRenew)
		if err != nil {
			continue
		}

		subscriptions = append(subscriptions, map[string]interface{}{
			"id":            id,
			"plan_id":       planID,
			"plan_name":     planName,
			"plan_type":     planType,
			"status":        status,
			"start_date":    startDate,
			"end_date":      endDate,
			"price":         price,
			"billing_cycle": billingCycle,
			"auto_renew":    autoRenew,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"data":         subscriptions,
	})
}

// ListTransactions returns user's transaction history
func (h *BillingHandler) ListTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	txType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT id, type, amount, balance_before, balance_after,
		       description, reference_id, reference_type, status, created_at
		FROM transactions
		WHERE user_id = $1
	`
	countQuery := "SELECT COUNT(*) FROM transactions WHERE user_id = $1"
	args := []interface{}{userID}
	argIndex := 2

	if txType != "" {
		query += " AND type = $" + strconv.Itoa(argIndex)
		countQuery += " AND type = $" + strconv.Itoa(argIndex)
		args = append(args, txType)
		argIndex++
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(argIndex)
	args = append(args, limit)
	argIndex++
	query += " OFFSET $" + strconv.Itoa(argIndex)
	args = append(args, offset)

	// Get total count
	var total int
	h.db.QueryRow(ctx, countQuery, args[:1]...).Scan(&total)

	// Get transactions
	rows, err := h.db.Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("Failed to query transactions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Internal error"})
		return
	}
	defer rows.Close()

	transactions := []map[string]interface{}{}
	for rows.Next() {
		var id string
		var txType string
		var amount, balanceBefore, balanceAfter float64
		var description *string
		var referenceID *string
		var referenceType *string
		var status string
		var createdAt time.Time

		err := rows.Scan(&id, &txType, &amount, &balanceBefore, &balanceAfter,
			&description, &referenceID, &referenceType, &status, &createdAt)
		if err != nil {
			continue
		}

		tx := map[string]interface{}{
			"id":              id,
			"type":            txType,
			"amount":          amount,
			"balance_before":  balanceBefore,
			"balance_after":   balanceAfter,
			"status":          status,
			"created_at":      createdAt,
		}
		if description != nil {
			tx["description"] = *description
		}
		if referenceID != nil {
			tx["reference_id"] = *referenceID
		}
		if referenceType != nil {
			tx["reference_type"] = *referenceType
		}

		transactions = append(transactions, tx)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"total":         total,
			"page":          page,
			"limit":         limit,
			"transactions":  transactions,
		},
	})
}

// RechargeRequest represents balance recharge input
type RechargeRequest struct {
	Amount  float64 `json:"amount" binding:"required,min=1,max=100000"`
	Method  string  `json:"method" binding:"required,oneof=alipay wechatpay stripe banktransfer"`
}

// Recharge handles balance recharge
func (h *BillingHandler) Recharge(c *gin.Context) {
	userID, _ := c.Get("userID")
	ctx := c.Request.Context()

	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}

	// Generate payment URL based on method
	var paymentURL string
	switch req.Method {
	case "alipay":
		paymentURL = "https://payment.agenthub.com/alipay?amount=" + strconv.FormatFloat(req.Amount, 'f', 2, 64)
	case "wechatpay":
		paymentURL = "https://payment.agenthub.com/wechatpay?amount=" + strconv.FormatFloat(req.Amount, 'f', 2, 64)
	case "stripe":
		paymentURL = "https://payment.agenthub.com/stripe?amount=" + strconv.FormatFloat(req.Amount, 'f', 2, 64)
	case "banktransfer":
		paymentURL = "/api/v1/billing/bank-transfer?amount=" + strconv.FormatFloat(req.Amount, 'f', 2, 64)
	}

	// Create pending transaction
	var txID uuid.UUID
	h.db.QueryRow(ctx, `
		INSERT INTO transactions (user_id, type, amount, balance_before, balance_after, description, status)
		SELECT $1, 'recharge', $2, balance, balance, $3, 'pending'
		FROM users WHERE id = $1
		RETURNING id
	`, userID, req.Amount, "Balance recharge via "+req.Method).Scan(&txID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"transaction_id": txID.String(),
			"amount":         req.Amount,
			"payment_url":    paymentURL,
			"expires_at":     time.Now().Add(30 * time.Minute),
		},
	})
}
