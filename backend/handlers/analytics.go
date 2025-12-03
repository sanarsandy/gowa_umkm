package handlers

import (
	"log"
	"net/http"
	"time"

	"gowa-backend/db"

	"github.com/labstack/echo/v4"
)

// DailyStats represents daily analytics data
type DailyStats struct {
	Date             string  `db:"date" json:"date"`
	MessagesReceived int     `db:"messages_received" json:"messages_received"`
	MessagesSent     int     `db:"messages_sent" json:"messages_sent"`
	UniqueCustomers  int     `db:"unique_customers" json:"unique_customers"`
	NewCustomers     int     `db:"new_customers" json:"new_customers"`
	AIResponses      int     `db:"ai_responses" json:"ai_responses"`
	AIEscalations    int     `db:"ai_escalations" json:"ai_escalations"`
	AIAvgConfidence  float64 `db:"ai_avg_confidence" json:"ai_avg_confidence"`
	AITokensUsed     int     `db:"ai_tokens_used" json:"ai_tokens_used"`
	AICostUSD        float64 `db:"ai_cost_usd" json:"ai_cost_usd"`
	BroadcastsSent   int     `db:"broadcasts_sent" json:"broadcasts_sent"`
}

// GetAnalyticsOverview returns overview analytics
// GET /api/analytics/overview
func GetAnalyticsOverview(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var overview struct {
		TotalCustomers     int     `db:"total_customers" json:"total_customers"`
		TotalMessages      int     `db:"total_messages" json:"total_messages"`
		TodayMessages      int     `db:"today_messages" json:"today_messages"`
		WeekMessages       int     `db:"week_messages" json:"week_messages"`
		MonthMessages      int     `db:"month_messages" json:"month_messages"`
		NewCustomersToday  int     `db:"new_customers_today" json:"new_customers_today"`
		NewCustomersWeek   int     `db:"new_customers_week" json:"new_customers_week"`
		NewCustomersMonth  int     `db:"new_customers_month" json:"new_customers_month"`
		AIResponsesTotal   int     `db:"ai_responses_total" json:"ai_responses_total"`
		AIEscalationsTotal int     `db:"ai_escalations_total" json:"ai_escalations_total"`
		AITotalCost        float64 `db:"ai_total_cost" json:"ai_total_cost"`
		AITotalTokens      int64   `db:"ai_total_tokens" json:"ai_total_tokens"`
		ActiveCustomers    int     `db:"active_customers" json:"active_customers"`
	}

	// Total customers
	db.DB.Get(&overview.TotalCustomers, `SELECT COUNT(*) FROM customer_insights WHERE tenant_id = $1`, tenantID)

	// Total messages
	db.DB.Get(&overview.TotalMessages, `SELECT COUNT(*) FROM whatsapp_messages WHERE tenant_id = $1`, tenantID)

	// Today's messages
	db.DB.Get(&overview.TodayMessages, `
		SELECT COUNT(*) FROM whatsapp_messages 
		WHERE tenant_id = $1 AND DATE(to_timestamp(timestamp)) = CURRENT_DATE
	`, tenantID)

	// This week's messages
	db.DB.Get(&overview.WeekMessages, `
		SELECT COUNT(*) FROM whatsapp_messages 
		WHERE tenant_id = $1 AND to_timestamp(timestamp) >= NOW() - INTERVAL '7 days'
	`, tenantID)

	// This month's messages
	db.DB.Get(&overview.MonthMessages, `
		SELECT COUNT(*) FROM whatsapp_messages 
		WHERE tenant_id = $1 AND to_timestamp(timestamp) >= NOW() - INTERVAL '30 days'
	`, tenantID)

	// New customers today
	db.DB.Get(&overview.NewCustomersToday, `
		SELECT COUNT(*) FROM customer_insights 
		WHERE tenant_id = $1 AND DATE(created_at) = CURRENT_DATE
	`, tenantID)

	// New customers this week
	db.DB.Get(&overview.NewCustomersWeek, `
		SELECT COUNT(*) FROM customer_insights 
		WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '7 days'
	`, tenantID)

	// New customers this month
	db.DB.Get(&overview.NewCustomersMonth, `
		SELECT COUNT(*) FROM customer_insights 
		WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
	`, tenantID)

	// AI stats from ai_configs
	db.DB.QueryRow(`
		SELECT 
			COALESCE(total_requests, 0),
			COALESCE(total_cost_usd, 0),
			COALESCE(total_tokens_used, 0)
		FROM ai_configs WHERE tenant_id = $1
	`, tenantID).Scan(&overview.AIResponsesTotal, &overview.AITotalCost, &overview.AITotalTokens)

	// AI escalations from logs
	db.DB.Get(&overview.AIEscalationsTotal, `
		SELECT COUNT(*) FROM ai_conversation_logs 
		WHERE tenant_id = $1 AND action_taken = 'escalated'
	`, tenantID)

	// Active customers (messaged in last 7 days)
	db.DB.Get(&overview.ActiveCustomers, `
		SELECT COUNT(DISTINCT chat_jid) FROM whatsapp_messages 
		WHERE tenant_id = $1 AND to_timestamp(timestamp) >= NOW() - INTERVAL '7 days' AND NOT is_from_me
	`, tenantID)

	return c.JSON(http.StatusOK, overview)
}

// GetAnalyticsMessages returns message analytics over time
// GET /api/analytics/messages?period=7d|30d|90d
func GetAnalyticsMessages(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	period := c.QueryParam("period")
	if period == "" {
		period = "30d"
	}

	var interval string
	switch period {
	case "7d":
		interval = "7 days"
	case "30d":
		interval = "30 days"
	case "90d":
		interval = "90 days"
	default:
		interval = "30 days"
	}

	type MessageStat struct {
		Date     string `db:"date" json:"date"`
		Received int    `db:"received" json:"received"`
		Sent     int    `db:"sent" json:"sent"`
	}

	var stats []MessageStat

	query := `
		SELECT 
			to_char(DATE(to_timestamp(timestamp)), 'YYYY-MM-DD') as date,
			COUNT(*) FILTER (WHERE NOT is_from_me) as received,
			COUNT(*) FILTER (WHERE is_from_me) as sent
		FROM whatsapp_messages
		WHERE tenant_id = $1 AND to_timestamp(timestamp) >= NOW() - INTERVAL '` + interval + `'
		GROUP BY DATE(to_timestamp(timestamp))
		ORDER BY DATE(to_timestamp(timestamp)) ASC
	`

	err := db.DB.Select(&stats, query, tenantID)
	if err != nil {
		log.Println("[Analytics Messages] Error:", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get message analytics")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"period": period,
		"data":   stats,
	})
}

// GetAnalyticsCustomers returns customer analytics
// GET /api/analytics/customers?period=7d|30d|90d
func GetAnalyticsCustomers(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	period := c.QueryParam("period")
	if period == "" {
		period = "30d"
	}

	var interval string
	switch period {
	case "7d":
		interval = "7 days"
	case "30d":
		interval = "30 days"
	case "90d":
		interval = "90 days"
	default:
		interval = "30 days"
	}

	type CustomerStat struct {
		Date         string `db:"date"`
		NewCustomers int    `db:"new_customers"`
	}

	var stats []CustomerStat

	query := `
		SELECT 
			to_char(DATE(created_at), 'YYYY-MM-DD') as date,
			COUNT(*)::int as new_customers
		FROM customer_insights
		WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '` + interval + `'
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at) ASC
	`

	err := db.DB.Select(&stats, query, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get customer analytics: "+err.Error())
	}

	// Get cumulative count
	var totalBefore int
	db.DB.Get(&totalBefore, `
		SELECT COUNT(*) FROM customer_insights 
		WHERE tenant_id = $1 AND created_at < NOW() - INTERVAL '`+interval+`'
	`, tenantID)

	// Convert to result with cumulative
	type ResultStat struct {
		Date         string `json:"date"`
		NewCustomers int    `json:"new_customers"`
		TotalActive  int    `json:"total_active"`
	}

	result := make([]ResultStat, len(stats))
	cumulative := totalBefore
	for i, s := range stats {
		cumulative += s.NewCustomers
		result[i] = ResultStat{
			Date:         s.Date,
			NewCustomers: s.NewCustomers,
			TotalActive:  cumulative,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"period": period,
		"data":   result,
	})
}

// GetAnalyticsAI returns AI performance analytics
// GET /api/analytics/ai?period=7d|30d|90d
func GetAnalyticsAI(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)
	period := c.QueryParam("period")
	if period == "" {
		period = "30d"
	}

	var interval string
	switch period {
	case "7d":
		interval = "7 days"
	case "30d":
		interval = "30 days"
	case "90d":
		interval = "90 days"
	default:
		interval = "30 days"
	}

	type AIStat struct {
		Date          string  `db:"date"`
		Responses     int     `db:"responses"`
		Escalations   int     `db:"escalations"`
		AvgConfidence float64 `db:"avg_confidence"`
		TokensUsed    int     `db:"tokens_used"`
		CostUSD       float64 `db:"cost_usd"`
	}

	var stats []AIStat

	query := `
		SELECT 
			to_char(DATE(created_at), 'YYYY-MM-DD') as date,
			COUNT(*)::int as responses,
			COUNT(*) FILTER (WHERE action_taken = 'escalated')::int as escalations,
			COALESCE(AVG(confidence_score), 0) as avg_confidence,
			COALESCE(SUM(tokens_used), 0)::int as tokens_used,
			COALESCE(SUM(cost_usd), 0) as cost_usd
		FROM ai_conversation_logs
		WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '` + interval + `'
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at) ASC
	`

	err := db.DB.Select(&stats, query, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get AI analytics: "+err.Error())
	}

	// Calculate totals
	var totals struct {
		TotalResponses   int     `json:"total_responses"`
		TotalEscalations int     `json:"total_escalations"`
		AvgConfidence    float64 `json:"avg_confidence"`
		TotalTokens      int     `json:"total_tokens"`
		TotalCost        float64 `json:"total_cost"`
		AutoReplyRate    float64 `json:"auto_reply_rate"`
	}

	for _, s := range stats {
		totals.TotalResponses += s.Responses
		totals.TotalEscalations += s.Escalations
		totals.TotalTokens += s.TokensUsed
		totals.TotalCost += s.CostUSD
	}

	if totals.TotalResponses > 0 {
		totals.AutoReplyRate = float64(totals.TotalResponses-totals.TotalEscalations) / float64(totals.TotalResponses) * 100
	}

	db.DB.Get(&totals.AvgConfidence, `
		SELECT COALESCE(AVG(confidence_score), 0) FROM ai_conversation_logs 
		WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '`+interval+`'
	`, tenantID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"period": period,
		"data":   stats,
		"totals": totals,
	})
}

// GetAnalyticsTopCustomers returns top customers by message count
// GET /api/analytics/top-customers?limit=10
func GetAnalyticsTopCustomers(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var customers []struct {
		ID           string    `db:"id" json:"id"`
		PhoneNumber  string    `db:"phone_number" json:"phone_number"`
		Name         string    `db:"name" json:"name"`
		MessageCount int       `db:"message_count" json:"message_count"`
		LastMessage  time.Time `db:"last_message" json:"last_message"`
		LeadScore    int       `db:"lead_score" json:"lead_score"`
	}

	query := `
		SELECT 
			c.id,
			COALESCE(c.customer_phone, c.customer_jid) as phone_number,
			COALESCE(c.customer_name, c.customer_phone, c.customer_jid) as name,
			COALESCE(c.message_count, 0) as message_count,
			c.last_message_at as last_message,
			COALESCE(c.lead_score, 0) as lead_score
		FROM customer_insights c
		WHERE c.tenant_id = $1
		ORDER BY message_count DESC
		LIMIT 10
	`

	err := db.DB.Select(&customers, query, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get top customers")
	}

	return c.JSON(http.StatusOK, customers)
}

// GetAnalyticsHourly returns message distribution by hour
// GET /api/analytics/hourly
func GetAnalyticsHourly(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	type HourlyStat struct {
		Hour  int `db:"hour"`
		Count int `db:"count"`
	}

	var stats []HourlyStat

	query := `
		SELECT 
			EXTRACT(HOUR FROM to_timestamp(timestamp))::int as hour,
			COUNT(*)::int as count
		FROM whatsapp_messages
		WHERE tenant_id = $1 AND to_timestamp(timestamp) >= NOW() - INTERVAL '30 days'
		GROUP BY EXTRACT(HOUR FROM to_timestamp(timestamp))::int
		ORDER BY hour
	`

	err := db.DB.Select(&stats, query, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get hourly analytics: "+err.Error())
	}

	// Fill missing hours with 0
	hourMap := make(map[int]int)
	for _, s := range stats {
		hourMap[s.Hour] = s.Count
	}

	type ResultStat struct {
		Hour  int `json:"hour"`
		Count int `json:"count"`
	}

	result := make([]ResultStat, 24)
	for i := 0; i < 24; i++ {
		result[i].Hour = i
		result[i].Count = hourMap[i]
	}

	return c.JSON(http.StatusOK, result)
}

// GetAnalyticsIntents returns detected intent distribution
// GET /api/analytics/intents
func GetAnalyticsIntents(c echo.Context) error {
	tenantID := getTenantIDFromContext(c)

	var stats []struct {
		Intent string `db:"intent" json:"intent"`
		Count  int    `db:"count" json:"count"`
	}

	query := `
		SELECT 
			detected_intent as intent,
			COUNT(*) as count
		FROM ai_conversation_logs
		WHERE tenant_id = $1 AND detected_intent IS NOT NULL AND detected_intent != ''
		GROUP BY detected_intent
		ORDER BY count DESC
	`

	err := db.DB.Select(&stats, query, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get intent analytics")
	}

	return c.JSON(http.StatusOK, stats)
}

