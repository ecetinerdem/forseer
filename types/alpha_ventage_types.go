package types

type AlphaVentageStockResponse struct {
	MetaData   MetaData               `json:"Meta Data"`
	TimeSeries map[string]MonthlyData `json:"Monthly Time Series"`
}

type MetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	TimeZone      string `json:"4. Time Zone"`
}

type MonthlyData struct {
	Open   float64 `json:"1. open"`
	High   float64 `json:"2. high"`
	Low    float64 `json:"3. low"`
	Close  float64 `json:"4. close"`
	Volume int64   `json:"5. volume"`
}
