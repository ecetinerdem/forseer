package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ecetinerdem/forseer/api/utils"
	"github.com/ecetinerdem/forseer/middleware"
	"github.com/ecetinerdem/forseer/types"
	"github.com/go-chi/chi/v5"
)

// HandleGetPortfolio returns the authenticated user's portfolio with all stocks
func (s *Server) HandleGetPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	if s.db == nil {
		http.Error(w, "Database connection not available", http.StatusInternalServerError)
		return
	}

	portfolio, err := s.db.GetUserPortfolio(ctx, user.ID)
	if err != nil {
		http.Error(w, "Could not get portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(portfolio); err != nil {
		http.Error(w, "Could not encode portfolio", http.StatusInternalServerError)
		return
	}
}

// HandleGetStockByID returns a specific stock if it belongs to the authenticated user
func (s *Server) HandleGetStockByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	stockID := chi.URLParam(r, "id")
	if stockID == "" {
		http.Error(w, "Stock ID cannot be empty", http.StatusBadRequest)
		return
	}

	stock, err := s.db.GetUserStockByID(ctx, user.ID, stockID)
	if err != nil {
		var ownershipErr *types.StockOwnershipError
		if errors.As(err, &ownershipErr) {
			http.Error(w, "Stock not found or you don't have access to it", http.StatusNotFound)
			return
		}
		http.Error(w, "Could not retrieve stock", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stock); err != nil {
		http.Error(w, "Could not encode stock", http.StatusInternalServerError)
		return
	}
}

// HandleAddStockToPortfolio adds a stock to the authenticated user's portfolio
func (s *Server) HandleAddStockToPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	stockSymbol := chi.URLParam(r, "symbol")
	if stockSymbol == "" {
		http.Error(w, "Stock symbol cannot be empty", http.StatusBadRequest)
		return
	}

	// Check if user already has this stock
	existingStock, err := s.db.GetUserStockBySymbol(ctx, user.ID, stockSymbol)
	if err == nil && existingStock != nil {
		http.Error(w, "Stock already exists in your portfolio", http.StatusConflict)
		return
	}

	// Fetch stock data from external API
	stock, err := utils.GetAlphaVentageStock(user, stockSymbol)
	if err != nil {
		http.Error(w, "Error while fetching stock data", http.StatusInternalServerError)
		return
	}

	// Add stock to user's portfolio
	addedStock, err := s.db.AddStockToUserPortfolio(ctx, user.ID, stock)
	if err != nil {
		http.Error(w, "Could not add stock to portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(addedStock); err != nil {
		http.Error(w, "Could not encode stock", http.StatusInternalServerError)
		return
	}
}

// HandleDeleteStockByID deletes a stock from the authenticated user's portfolio
func (s *Server) HandleDeleteStockByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	stockID := chi.URLParam(r, "id")
	if stockID == "" {
		http.Error(w, "Stock ID cannot be empty", http.StatusBadRequest)
		return
	}

	err := s.db.DeleteUserStockByID(ctx, user.ID, stockID)
	if err != nil {
		var ownershipErr *types.StockOwnershipError
		if errors.As(err, &ownershipErr) {
			http.Error(w, "Stock not found or you don't have access to it", http.StatusNotFound)
			return
		}
		http.Error(w, "Could not delete stock", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Stock deleted successfully",
		"stockId": stockID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
		return
	}
}

// HandleGetStockBySymbol returns a stock by symbol if it belongs to the authenticated user
func (s *Server) HandleGetStockBySymbol(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	symbol := chi.URLParam(r, "symbol")
	if symbol == "" {
		http.Error(w, "Stock symbol cannot be empty", http.StatusBadRequest)
		return
	}

	stock, err := s.db.GetUserStockBySymbol(ctx, user.ID, symbol)
	if err != nil {
		http.Error(w, "Stock with given symbol not found in your portfolio", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stock); err != nil {
		http.Error(w, "Could not encode stock", http.StatusInternalServerError)
		return
	}
}

// HandleGetUserStocks returns all stocks in the authenticated user's portfolio
func (s *Server) HandleGetUserStocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	stocks, err := s.db.GetUserStocks(ctx, user.ID)
	if err != nil {
		http.Error(w, "Could not retrieve stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stocks); err != nil {
		http.Error(w, "Could not encode stocks", http.StatusInternalServerError)
		return
	}
}

// HandleCreatePortfolio creates a new portfolio for the authenticated user
func (s *Server) HandleCreatePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		req.Name = "My Portfolio"
	}

	portfolio, err := s.db.CreateUserPortfolio(ctx, user.ID, req.Name)
	if err != nil {
		http.Error(w, "Could not create portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(portfolio); err != nil {
		http.Error(w, "Could not encode portfolio", http.StatusInternalServerError)
		return
	}
}
