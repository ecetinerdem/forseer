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
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY&symbol=%s&apikey=%s", stockSymbol, apiKey)

	r, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("error getting stock data %w", err)
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&alphaVentageStockResponse)

	if err != nil {
		return nil, fmt.Errorf("error decoding stock data %w", err)
	}

	var latestDate string
	for date := range alphaVentageStockResponse.TimeSeries {
		if latestDate == "" || date > latestDate {
			latestDate = date
		}
	}

	returnStock.Symbol = alphaVentageStockResponse.MetaData.Symbol
	returnStock.Month = time.Now().Month().String()
	returnStock.Open = alphaVentageStockResponse.TimeSeries[latestDate].Open
	returnStock.High = alphaVentageStockResponse.TimeSeries[latestDate].High
	returnStock.Close = alphaVentageStockResponse.TimeSeries[latestDate].Close
	returnStock.Volume = alphaVentageStockResponse.TimeSeries[latestDate].Volume
	returnStock.PortfolioID = user.Portfolio.ID

	return &returnStock, nil

}
