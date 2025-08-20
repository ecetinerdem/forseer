package types

type Stock struct {
	ID          string `json:"id"`
	PortfolioID string `json:"portfolio_id"`
	Symbol      string `json:"symbol"`
	Month       string `json:"month"`
	Open        string `json:"open"`
	High        string `json:"high"`
	Low         string `json:"low"`
	Close       string `json:"close"`
	Volume      string `json:"volume"`
}

type Portfolio struct {
	ID     string  `json:"id"`
	UserID string  `json:"user_id"`
	Stocks []Stock `json:"stocks"`
}
