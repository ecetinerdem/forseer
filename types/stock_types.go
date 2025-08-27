package types

import "time"

type Portfolio struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Stocks    []Stock   `json:"stocks,omitempty"` // Optional for when you want to include stocks
}

type Stock struct {
	ID          string    `json:"id"`
	PortfolioID string    `json:"portfolio_id"`
	Symbol      string    `json:"symbol"`
	Month       string    `json:"month"` // Format: YYYY-MM
	Open        float64   `json:"open"`  // Changed to float64 for decimal values
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	Volume      int64     `json:"volume"` // Changed to int64 for bigint
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
