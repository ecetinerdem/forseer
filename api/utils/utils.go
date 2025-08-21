package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ecetinerdem/forseer/types"
)

func GetAlphaVentageStock(user *types.User, stockSymbol string) (*types.Stock, error) {
	var alphaVentageStockResponse types.AlphaVentageStockResponse
	var returnStock types.Stock

	apiKey := os.Getenv("ALPHAVENTAGE_API_KEY")
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY&symbol=%s&apikey=%w", stockSymbol, apiKey)

	r, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Error getting stock data %w", err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&alphaVentageStockResponse)

	if err != nil {
		return nil, fmt.Errorf("Error decoding stock data %w", err)
	}

	var latestDate string
	for date := range alphaVentageStockResponse.TimeSeries {
		if latestDate == "" || date > latestDate {
			latestDate = date
		}
	}

	returnStock.Symbol = alphaVentageStockResponse.MetaData.Symbol
	returnStock.Month = string(time.Now().Month())
	returnStock.Open = alphaVentageStockResponse.TimeSeries[latestDate].Open
	returnStock.High = alphaVentageStockResponse.TimeSeries[latestDate].Open
	returnStock.Close = alphaVentageStockResponse.TimeSeries[latestDate].Open
	returnStock.Volume = alphaVentageStockResponse.TimeSeries[latestDate].Open
	returnStock.PortfolioID = user.Portfolio.ID

}
