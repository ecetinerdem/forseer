package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ecetinerdem/forseer/types"
)

type PortfolioRepo interface {
	// Portfolio operations
	GetUserPortfolio(ctx context.Context, userID string) (*types.Portfolio, error)
	CreateUserPortfolio(ctx context.Context, userID, portfolioName string) (*types.Portfolio, error)

	// Stock operations - all user-scoped
	GetUserStocks(ctx context.Context, userID string) ([]types.Stock, error)
	GetUserStockByID(ctx context.Context, userID, stockID string) (*types.Stock, error)
	GetUserStockBySymbol(ctx context.Context, userID, stockSymbol string) (*types.Stock, error)
	AddStockToUserPortfolio(ctx context.Context, userID string, stock *types.Stock) (*types.Stock, error)
	DeleteUserStockByID(ctx context.Context, userID, stockID string) error

	// Ownership validation
	UserOwnsStock(ctx context.Context, userID, stockID string) (bool, error)
	UserOwnsPortfolio(ctx context.Context, userID, portfolioID string) (bool, error)
}

// GetUserPortfolio retrieves the user's portfolio with all stocks
func (db *DB) GetUserPortfolio(ctx context.Context, userID string) (*types.Portfolio, error) {
	// First get or create portfolio
	portfolio, err := db.getOrCreatePortfolio(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	// Get all stocks for this portfolio
	stocks, err := db.GetUserStocks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stocks: %w", err)
	}

	portfolio.Stocks = stocks
	return portfolio, nil
}

// CreateUserPortfolio creates a new portfolio for the user
func (db *DB) CreateUserPortfolio(ctx context.Context, userID, portfolioName string) (*types.Portfolio, error) {
	query := `
		INSERT INTO portfolios (user_id, name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, user_id, name, created_at, updated_at
	`

	var portfolio types.Portfolio
	err := db.QueryRowContext(ctx, query, userID, portfolioName).Scan(
		&portfolio.ID,
		&portfolio.UserID,
		&portfolio.Name,
		&portfolio.CreatedAt,
		&portfolio.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create portfolio: %w", err)
	}

	return &portfolio, nil
}

// GetUserStocks returns all stocks belonging to the user
func (db *DB) GetUserStocks(ctx context.Context, userID string) ([]types.Stock, error) {
	query := `
		SELECT s.id, s.portfolio_id, s.symbol, s.month, s.open, s.high, s.low, s.close, s.volume, s.created_at, s.updated_at
		FROM stocks s
		INNER JOIN portfolios p ON s.portfolio_id = p.id
		WHERE p.user_id = $1
		ORDER BY s.created_at DESC
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user stocks: %w", err)
	}
	defer rows.Close()

	var stocks []types.Stock
	for rows.Next() {
		var stock types.Stock
		err := rows.Scan(
			&stock.ID,
			&stock.PortfolioID,
			&stock.Symbol,
			&stock.Month,
			&stock.Open,
			&stock.High,
			&stock.Low,
			&stock.Close,
			&stock.Volume,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// GetUserStockByID retrieves a specific stock if it belongs to the user
func (db *DB) GetUserStockByID(ctx context.Context, userID, stockID string) (*types.Stock, error) {
	query := `
		SELECT s.id, s.portfolio_id, s.symbol, s.month, s.open, s.high, s.low, s.close, s.volume, s.created_at, s.updated_at
		FROM stocks s
		INNER JOIN portfolios p ON s.portfolio_id = p.id
		WHERE s.id = $1 AND p.user_id = $2
	`

	var stock types.Stock
	err := db.QueryRowContext(ctx, query, stockID, userID).Scan(
		&stock.ID,
		&stock.PortfolioID,
		&stock.Symbol,
		&stock.Month,
		&stock.Open,
		&stock.High,
		&stock.Low,
		&stock.Close,
		&stock.Volume,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &types.StockOwnershipError{UserID: userID, StockID: stockID}
		}
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}

	return &stock, nil
}

// GetUserStockBySymbol retrieves a stock by symbol if it belongs to the user
func (db *DB) GetUserStockBySymbol(ctx context.Context, userID, stockSymbol string) (*types.Stock, error) {
	query := `
		SELECT s.id, s.portfolio_id, s.symbol, s.month, s.open, s.high, s.low, s.close, s.volume, s.created_at, s.updated_at
		FROM stocks s
		INNER JOIN portfolios p ON s.portfolio_id = p.id
		WHERE s.symbol = $1 AND p.user_id = $2
		ORDER BY s.created_at DESC
		LIMIT 1
	`

	var stock types.Stock
	err := db.QueryRowContext(ctx, query, stockSymbol, userID).Scan(
		&stock.ID,
		&stock.PortfolioID,
		&stock.Symbol,
		&stock.Month,
		&stock.Open,
		&stock.High,
		&stock.Low,
		&stock.Close,
		&stock.Volume,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no stock with symbol %s found for user %s", stockSymbol, userID)
		}
		return nil, fmt.Errorf("failed to get stock by symbol: %w", err)
	}

	return &stock, nil
}

// AddStockToUserPortfolio adds a stock to the user's portfolio
func (db *DB) AddStockToUserPortfolio(ctx context.Context, userID string, stock *types.Stock) (*types.Stock, error) {
	// First ensure user has a portfolio
	portfolio, err := db.getOrCreatePortfolio(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolio: %w", err)
	}

	// Insert the stock
	query := `
		INSERT INTO stocks (portfolio_id, symbol, month, open, high, low, close, volume, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, portfolio_id, symbol, month, open, high, low, close, volume, created_at, updated_at
	`

	var newStock types.Stock
	err = db.QueryRowContext(ctx, query,
		portfolio.ID,
		stock.Symbol,
		stock.Month,
		stock.Open,
		stock.High,
		stock.Low,
		stock.Close,
		stock.Volume,
	).Scan(
		&newStock.ID,
		&newStock.PortfolioID,
		&newStock.Symbol,
		&newStock.Month,
		&newStock.Open,
		&newStock.High,
		&newStock.Low,
		&newStock.Close,
		&newStock.Volume,
		&newStock.CreatedAt,
		&newStock.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("could not save the stock: %w", err)
	}

	return &newStock, nil
}

// DeleteUserStockByID deletes a stock if it belongs to the user
func (db *DB) DeleteUserStockByID(ctx context.Context, userID, stockID string) error {
	// Verify ownership and delete in one query
	query := `
		DELETE FROM stocks 
		WHERE id = $1 
		AND portfolio_id IN (
			SELECT id FROM portfolios WHERE user_id = $2
		)
	`

	result, err := db.ExecContext(ctx, query, stockID, userID)
	if err != nil {
		return fmt.Errorf("could not delete stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return &types.StockOwnershipError{UserID: userID, StockID: stockID}
	}

	return nil
}

// UserOwnsStock verifies if a user owns a specific stock
func (db *DB) UserOwnsStock(ctx context.Context, userID, stockID string) (bool, error) {
	query := `
		SELECT 1
		FROM stocks s
		INNER JOIN portfolios p ON s.portfolio_id = p.id
		WHERE s.id = $1 AND p.user_id = $2
	`

	var exists int
	err := db.QueryRowContext(ctx, query, stockID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check stock ownership: %w", err)
	}

	return true, nil
}

// UserOwnsPortfolio verifies if a user owns a specific portfolio
func (db *DB) UserOwnsPortfolio(ctx context.Context, userID, portfolioID string) (bool, error) {
	query := `
		SELECT 1
		FROM portfolios
		WHERE id = $1 AND user_id = $2
	`

	var exists int
	err := db.QueryRowContext(ctx, query, portfolioID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check portfolio ownership: %w", err)
	}

	return true, nil
}

// getOrCreatePortfolio is a helper that ensures a user has a portfolio
func (db *DB) getOrCreatePortfolio(ctx context.Context, userID string) (*types.Portfolio, error) {
	// Try to get existing portfolio
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM portfolios
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var portfolio types.Portfolio
	err := db.QueryRowContext(ctx, query, userID).Scan(
		&portfolio.ID,
		&portfolio.UserID,
		&portfolio.Name,
		&portfolio.CreatedAt,
		&portfolio.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Create default portfolio
			return db.CreateUserPortfolio(ctx, userID, "My Portfolio")
		}
		return nil, fmt.Errorf("failed to query portfolio: %w", err)
	}

	return &portfolio, nil
}
