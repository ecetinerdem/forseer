package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ecetinerdem/forseer/types"
)

func GetAlphaVentageStock(stockSymbol string) (*types.Stock, error) {
	var alphaVentageStockResponse types.AlphaVentageStockResponse
	var returnStock types.Stock

	apiKey := os.Getenv("ALPHAVENTAGE_API_KEY")
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY&%s=IBM&%s=demo", stockSymbol, apiKey)

	response, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Error getting stock data %w", err)
	}

	err = json.NewDecoder(response.Body).Decode(&alphaVentageStockResponse)

	if err != nil {
		return nil, fmt.Errorf("Error decoding stock data %w", err)
	}

	returnStock.Symbol = alphaVentageStockResponse.MetaData.Symbol
	returnStock.Month = string(time.Now().Month())
	returnStock.Open = alphaVentageStockResponse.TimeSeries["2025-08-20"].Open
	returnStock.High = alphaVentageStockResponse.TimeSeries["2025-08-20"].Open
	returnStock.Close = alphaVentageStockResponse.TimeSeries["2025-08-20"].Open
	returnStock.Volume = alphaVentageStockResponse.TimeSeries["2025-08-20"].Open

}
