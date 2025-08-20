package database

import (
	"context"

	"github.com/ecetinerdem/forseer/types"
)

type PortfolioRepo interface {
	GetStocks(context.Context, string) (*types.Portfolio, error)
	GetStockByID(context.Context, string) (*types.Stock, error)
	GetStockBySymbol(context.Context, string) (*types.Stock, error)
	AddStockToPortfolio(context.Context, *types.Stock) (*types.Stock, error)
	DeleteStockByID(context.Context, string) error
}
