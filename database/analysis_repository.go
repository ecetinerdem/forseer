// Add to your database package (database/analysis.go)
package database

import (
	"context"
	"fmt"

	"github.com/ecetinerdem/forseer/types"
)

type AnalysisRepo interface {
	// Stock analysis methods
	SaveStockAnalysis(ctx context.Context, analysis *types.StockAnalysis) (*types.StockAnalysis, error)
	GetStockAnalysis(ctx context.Context, userID, stockID string) (*types.StockAnalysis, error)
	GetUserStockAnalyses(ctx context.Context, userID string) ([]*types.StockAnalysis, error)
	DeleteStockAnalysis(ctx context.Context, userID, analysisID string) error

	// Portfolio analysis methods
	SavePortfolioAnalysis(ctx context.Context, analysis *types.PortfolioAnalysis) (*types.PortfolioAnalysis, error)
	GetPortfolioAnalysis(ctx context.Context, userID, portfolioID string) (*types.PortfolioAnalysis, error)
	GetUserPortfolioAnalyses(ctx context.Context, userID string) ([]*types.PortfolioAnalysis, error)
	DeletePortfolioAnalysis(ctx context.Context, userID, analysisID string) error
}

// SaveStockAnalysis saves a stock analysis to the database
func (db *DB) SaveStockAnalysis(ctx context.Context, analysis *types.StockAnalysis) (*types.StockAnalysis, error) {
	query := `
		INSERT INTO stock_analyses (stock_id, symbol, analysis, generated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, stock_id, symbol, analysis, generated_at, created_at, updated_at
	`

	var saved types.StockAnalysis
	err := db.QueryRowContext(ctx, query,
		analysis.StockID,
		analysis.Symbol,
		analysis.Analysis,
		analysis.GeneratedAt,
	).Scan(
		&saved.ID,
		&saved.StockID,
		&saved.Symbol,
		&saved.Analysis,
		&saved.GeneratedAt,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to save stock analysis: %w", err)
	}

	return &saved, nil
}

// GetStockAnalysis retrieves a stock analysis if the user owns the stock
func (db *DB) GetStockAnalysis(ctx context.Context, userID, stockID string) (*types.StockAnalysis, error) {
	query := `
		SELECT sa.id, sa.stock_id, sa.symbol, sa.analysis, sa.generated_at, sa.created_at, sa.updated_at
		FROM stock_analyses sa
		INNER JOIN stocks s ON sa.stock_id = s.id
		INNER JOIN portfolios p ON s.portfolio_id = p.id
		WHERE sa.stock_id = $1 AND p.user_id = $2
		ORDER BY sa.generated_at DESC
		LIMIT 1
	`

	var analysis types.StockAnalysis
	err := db.QueryRowContext(ctx, query, stockID, userID).Scan(
		&analysis.ID,
		&analysis.StockID,
		&analysis.Symbol,
		&analysis.Analysis,
		&analysis.GeneratedAt,
		&analysis.CreatedAt,
		&analysis.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stock analysis: %w", err)
	}

	return &analysis, nil
}

// GetUserStockAnalyses retrieves all stock analyses for a user
func (db *DB) GetUserStockAnalyses(ctx context.Context, userID string) ([]*types.StockAnalysis, error) {
	query := `
		SELECT sa.id, sa.stock_id, sa.symbol, sa.analysis, sa.generated_at, sa.created_at, sa.updated_at
		FROM stock_analyses sa
		INNER JOIN stocks s ON sa.stock_id = s.id
		INNER JOIN portfolios p ON s.portfolio_id = p.id
		WHERE p.user_id = $1
		ORDER BY sa.generated_at DESC
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query stock analyses: %w", err)
	}
	defer rows.Close()

	var analyses []*types.StockAnalysis
	for rows.Next() {
		var analysis types.StockAnalysis
		err := rows.Scan(
			&analysis.ID,
			&analysis.StockID,
			&analysis.Symbol,
			&analysis.Analysis,
			&analysis.GeneratedAt,
			&analysis.CreatedAt,
			&analysis.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock analysis: %w", err)
		}
		analyses = append(analyses, &analysis)
	}

	return analyses, nil
}

// DeleteStockAnalysis deletes a stock analysis if the user owns it
func (db *DB) DeleteStockAnalysis(ctx context.Context, userID, analysisID string) error {
	query := `
		DELETE FROM stock_analyses 
		WHERE id = $1 
		AND stock_id IN (
			SELECT s.id 
			FROM stocks s 
			INNER JOIN portfolios p ON s.portfolio_id = p.id 
			WHERE p.user_id = $2
		)
	`

	result, err := db.ExecContext(ctx, query, analysisID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete stock analysis: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("analysis not found or you don't have permission to delete it")
	}

	return nil
}

// SavePortfolioAnalysis saves a portfolio analysis to the database
func (db *DB) SavePortfolioAnalysis(ctx context.Context, analysis *types.PortfolioAnalysis) (*types.PortfolioAnalysis, error) {
	query := `
		INSERT INTO portfolio_analyses (portfolio_id, user_id, analysis, stock_count, generated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, portfolio_id, user_id, analysis, stock_count, generated_at, created_at, updated_at
	`

	var saved types.PortfolioAnalysis
	err := db.QueryRowContext(ctx, query,
		analysis.PortfolioID,
		analysis.UserID,
		analysis.Analysis,
		analysis.StockCount,
		analysis.GeneratedAt,
	).Scan(
		&saved.ID,
		&saved.PortfolioID,
		&saved.UserID,
		&saved.Analysis,
		&saved.StockCount,
		&saved.GeneratedAt,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to save portfolio analysis: %w", err)
	}

	return &saved, nil
}

// GetPortfolioAnalysis retrieves the latest portfolio analysis for a user's portfolio
func (db *DB) GetPortfolioAnalysis(ctx context.Context, userID, portfolioID string) (*types.PortfolioAnalysis, error) {
	query := `
		SELECT pa.id, pa.portfolio_id, pa.user_id, pa.analysis, pa.stock_count, pa.generated_at, pa.created_at, pa.updated_at
		FROM portfolio_analyses pa
		INNER JOIN portfolios p ON pa.portfolio_id = p.id
		WHERE pa.portfolio_id = $1 AND p.user_id = $2
		ORDER BY pa.generated_at DESC
		LIMIT 1
	`

	var analysis types.PortfolioAnalysis
	err := db.QueryRowContext(ctx, query, portfolioID, userID).Scan(
		&analysis.ID,
		&analysis.PortfolioID,
		&analysis.UserID,
		&analysis.Analysis,
		&analysis.StockCount,
		&analysis.GeneratedAt,
		&analysis.CreatedAt,
		&analysis.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio analysis: %w", err)
	}

	return &analysis, nil
}

// GetUserPortfolioAnalyses retrieves all portfolio analyses for a user
func (db *DB) GetUserPortfolioAnalyses(ctx context.Context, userID string) ([]*types.PortfolioAnalysis, error) {
	query := `
		SELECT pa.id, pa.portfolio_id, pa.user_id, pa.analysis, pa.stock_count, pa.generated_at, pa.created_at, pa.updated_at
		FROM portfolio_analyses pa
		INNER JOIN portfolios p ON pa.portfolio_id = p.id
		WHERE p.user_id = $1
		ORDER BY pa.generated_at DESC
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query portfolio analyses: %w", err)
	}
	defer rows.Close()

	var analyses []*types.PortfolioAnalysis
	for rows.Next() {
		var analysis types.PortfolioAnalysis
		err := rows.Scan(
			&analysis.ID,
			&analysis.PortfolioID,
			&analysis.UserID,
			&analysis.Analysis,
			&analysis.StockCount,
			&analysis.GeneratedAt,
			&analysis.CreatedAt,
			&analysis.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan portfolio analysis: %w", err)
		}
		analyses = append(analyses, &analysis)
	}

	return analyses, nil
}

// DeletePortfolioAnalysis deletes a portfolio analysis if the user owns it
func (db *DB) DeletePortfolioAnalysis(ctx context.Context, userID, analysisID string) error {
	query := `
		DELETE FROM portfolio_analyses 
		WHERE id = $1 
		AND portfolio_id IN (
			SELECT id FROM portfolios WHERE user_id = $2
		)
	`

	result, err := db.ExecContext(ctx, query, analysisID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete portfolio analysis: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("analysis not found or you don't have permission to delete it")
	}

	return nil
}
