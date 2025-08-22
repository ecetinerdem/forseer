package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ecetinerdem/forseer/types"
)

type PortfolioRepo interface {
	GetStocks(context.Context, string) (*types.Portfolio, error)
	GetStockByID(context.Context, string) (*types.Stock, error)
	GetStockBySymbol(context.Context, string) (*types.Stock, error)
	AddStockToPortfolio(context.Context, *types.Stock) (*types.Stock, error)
	DeleteStockByID(context.Context, string) error
}

func (db *DB) GetStocks(ctx context.Context, userID string) (*types.Portfolio, error) {
	query := `
		SELECT id, user_id, stocks FROM portfolios
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var portfolio types.Portfolio
	var stocksJSON []byte

	err := db.QueryRowContext(ctx, query, userID).Scan(
		&portfolio.ID,
		&portfolio.UserID,
		&stocksJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find portfolio %w", err)
	}

	err = json.Unmarshal(stocksJSON, &portfolio.Stocks)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stocks %w", err)
	}

	return &portfolio, nil
}

func (db *DB) GetStockByID(ctx context.Context, stockID string) (*types.Stock, error) {
	query := `
		SELECT id, portfolio_id, symbol, month, open, high, low, close, volume
		FROM stocks
		WHERE id = $1
	`

	var stock types.Stock

	err := db.QueryRowContext(ctx, query, stockID).Scan(
		&stock.ID,
		&stock.PortfolioID,
		&stock.Symbol,
		&stock.Month,
		&stock.High,
		&stock.Low,
		stock.Close,
		&stock.Volume,
	)

	if err != nil {
		return nil, fmt.Errorf("stock with given id does not exist %w", err)
	}

	return &stock, nil
}

func (db *DB) GetStockBySymbol(ctx context.Context, stockSymbol string) (*types.Stock, error) {
	query := `
		SELECT id, portfolio_id, symbol, month, open, high, low, close, volume
		FROM stocks
		WHERE symbol = $1
	`

	var stock types.Stock

	err := db.QueryRowContext(ctx, query, stockSymbol).Scan(
		&stock.ID,
		&stock.PortfolioID,
		&stock.Symbol,
		&stock.Month,
		&stock.High,
		&stock.Low,
		stock.Close,
		&stock.Volume,
	)

	if err != nil {
		return nil, fmt.Errorf("stock with given id does not exist %w", err)
	}

	return &stock, nil
}

func (db *DB) AddStockToPortfolio(ctx context.Context, stock *types.Stock) (*types.Stock, error) {
	query := `
		INSERT INTO stocks(portfolio_id, symbol month, open, high, low, close, volume),
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, portfolio_id, symbol, month, open, high, low, close, volume
	`

	var stockID string
	var portfolioID string
	var symbol string
	var month string
	var open string
	var high string
	var low string
	var close string
	var volume string

	err := db.QueryRowContext(ctx, query, stock).Scan(&stockID, &portfolioID, &symbol, &month, &open, &high, &low, &close, &volume)

	if err != nil {
		return nil, fmt.Errorf("could not save the stock %w", err)
	}

	return &types.Stock{
		ID:          stockID,
		PortfolioID: portfolioID,
		Symbol:      symbol,
		Month:       month,
		Open:        open,
		High:        high,
		Low:         low,
		Close:       close,
		Volume:      volume,
	}, nil
}

func (db *DB) DeleteStockByID(ctx context.Context, userID string) error {
	return nil
}
