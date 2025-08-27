// Add these to your types package (types/analysis.go)
package types

import (
	"time"
)

// StockAnalysis represents AI-generated analysis for a single stock
type StockAnalysis struct {
	ID          string    `json:"id" db:"id"`
	StockID     string    `json:"stock_id" db:"stock_id"`
	Symbol      string    `json:"symbol" db:"symbol"`
	Analysis    string    `json:"analysis" db:"analysis"`
	GeneratedAt time.Time `json:"generated_at" db:"generated_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PortfolioAnalysis represents AI-generated analysis for an entire portfolio
type PortfolioAnalysis struct {
	ID          string    `json:"id" db:"id"`
	PortfolioID string    `json:"portfolio_id" db:"portfolio_id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Analysis    string    `json:"analysis" db:"analysis"`
	StockCount  int       `json:"stock_count" db:"stock_count"`
	GeneratedAt time.Time `json:"generated_at" db:"generated_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AnalysisRequest represents the request payload for analysis
type AnalysisRequest struct {
	Type string `json:"type"` // "stock" or "portfolio"
	ID   string `json:"id"`   // stock_id or portfolio_id
}

// Custom error for analysis
type AnalysisError struct {
	Type    string
	Message string
}

func (e *AnalysisError) Error() string {
	return e.Message
}
