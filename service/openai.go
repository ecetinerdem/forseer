// services/openai.go
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ecetinerdem/forseer/types"
)

type OpenAIService struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AnalyzeStock analyzes a single stock and provides insights
func (o *OpenAIService) AnalyzeStock(ctx context.Context, stock *types.Stock) (*types.StockAnalysis, error) {
	prompt := o.buildStockAnalysisPrompt(stock)

	analysis, err := o.getCompletion(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock analysis: %w", err)
	}

	return &types.StockAnalysis{
		StockID:     stock.ID,
		Symbol:      stock.Symbol,
		Analysis:    analysis,
		GeneratedAt: time.Now(),
	}, nil
}

// AnalyzePortfolio analyzes an entire portfolio and provides insights
func (o *OpenAIService) AnalyzePortfolio(ctx context.Context, portfolio *types.Portfolio) (*types.PortfolioAnalysis, error) {
	if len(portfolio.Stocks) == 0 {
		return nil, fmt.Errorf("portfolio has no stocks to analyze")
	}

	prompt := o.buildPortfolioAnalysisPrompt(portfolio)

	analysis, err := o.getCompletion(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio analysis: %w", err)
	}

	return &types.PortfolioAnalysis{
		PortfolioID: portfolio.ID,
		UserID:      portfolio.UserID,
		Analysis:    analysis,
		StockCount:  len(portfolio.Stocks),
		GeneratedAt: time.Now(),
	}, nil
}

// getCompletion makes the actual API call to OpenAI
func (o *OpenAIService) getCompletion(ctx context.Context, prompt string) (string, error) {
	reqBody := OpenAIRequest{
		Model: "gpt-3.5-turbo", // You can change to "gpt-4" if you have access
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a professional financial analyst with expertise in stock market analysis. Provide detailed, actionable insights based on the stock data provided.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.3, // Lower temperature for more consistent analysis
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// buildStockAnalysisPrompt creates a detailed prompt for single stock analysis
func (o *OpenAIService) buildStockAnalysisPrompt(stock *types.Stock) string {
	return fmt.Sprintf(`
Please analyze the following stock data and provide a comprehensive analysis:

Stock Symbol: %s
Month: %s
Open Price: $%.2f
High Price: $%.2f
Low Price: $%.2f
Close Price: $%.2f
Volume: %d

Please provide analysis covering:
1. Price Performance: Analyze the price movement (open vs close, high vs low)
2. Volatility Assessment: Comment on the price volatility based on the high-low range
3. Volume Analysis: Interpret the trading volume significance
4. Technical Indicators: Basic technical analysis (price trends, support/resistance if applicable)
5. Risk Assessment: Identify potential risks based on the data
6. Recommendations: Provide actionable insights or recommendations

Please format your response in clear sections and be specific about the data points you're referencing.
`, stock.Symbol, stock.Month, stock.Open, stock.High, stock.Low, stock.Close, stock.Volume)
}

// buildPortfolioAnalysisPrompt creates a detailed prompt for portfolio analysis
func (o *OpenAIService) buildPortfolioAnalysisPrompt(portfolio *types.Portfolio) string {
	var stocksData strings.Builder
	stocksData.WriteString("Portfolio Stocks:\n\n")

	totalValue := 0.0
	for i, stock := range portfolio.Stocks {
		stocksData.WriteString(fmt.Sprintf(`%d. %s (%s):
   Open: $%.2f, High: $%.2f, Low: $%.2f, Close: $%.2f
   Volume: %d
   
`, i+1, stock.Symbol, stock.Month, stock.Open, stock.High, stock.Low, stock.Close, stock.Volume))
		totalValue += stock.Close
	}

	return fmt.Sprintf(`
Please analyze the following investment portfolio and provide a comprehensive analysis:

Portfolio Name: %s
Number of Stocks: %d
Total Portfolio Close Value: $%.2f

%s

Please provide analysis covering:
1. Portfolio Diversification: Analyze the spread across different stocks
2. Overall Performance: Comment on the general performance of the portfolio
3. Risk Assessment: Identify portfolio risks and volatility
4. Sector Analysis: If you can identify sectors from the stock symbols, provide sector insights
5. Performance Leaders and Laggards: Identify best and worst performing stocks
6. Portfolio Balance: Comment on the portfolio composition
7. Recommendations: Provide specific recommendations for portfolio optimization
8. Risk Management: Suggest risk management strategies

Please format your response in clear sections with specific data references and actionable insights.
`, portfolio.Name, len(portfolio.Stocks), totalValue, stocksData.String())
}
