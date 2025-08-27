// Add to your api package (api/analysis.go)
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ecetinerdem/forseer/middleware"
	"github.com/ecetinerdem/forseer/types"
	"github.com/go-chi/chi/v5"
)

// HandleAnalyzeStock generates AI analysis for a specific stock
func (s *Server) HandleAnalyzeStock(w http.ResponseWriter, r *http.Request) {
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

	// Get the stock and verify ownership
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

	// Generate analysis using OpenAI
	analysis, err := s.openAIService.AnalyzeStock(ctx, stock)
	if err != nil {
		http.Error(w, "Failed to generate stock analysis", http.StatusInternalServerError)
		return
	}

	// Save analysis to database
	savedAnalysis, err := s.db.SaveStockAnalysis(ctx, analysis)
	if err != nil {
		http.Error(w, "Failed to save analysis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(savedAnalysis); err != nil {
		http.Error(w, "Could not encode analysis", http.StatusInternalServerError)
		return
	}
}

// HandleAnalyzePortfolio generates AI analysis for the user's entire portfolio
func (s *Server) HandleAnalyzePortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	// Get the user's portfolio
	portfolio, err := s.db.GetUserPortfolio(ctx, user.ID)
	if err != nil {
		http.Error(w, "Could not get portfolio", http.StatusInternalServerError)
		return
	}

	if len(portfolio.Stocks) == 0 {
		http.Error(w, "Portfolio has no stocks to analyze", http.StatusBadRequest)
		return
	}

	// Generate analysis using OpenAI
	analysis, err := s.openAIService.AnalyzePortfolio(ctx, portfolio)
	if err != nil {
		http.Error(w, "Failed to generate portfolio analysis", http.StatusInternalServerError)
		return
	}

	// Save analysis to database
	savedAnalysis, err := s.db.SavePortfolioAnalysis(ctx, analysis)
	if err != nil {
		http.Error(w, "Failed to save analysis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(savedAnalysis); err != nil {
		http.Error(w, "Could not encode analysis", http.StatusInternalServerError)
		return
	}
}

// HandleGetStockAnalysis retrieves the latest analysis for a specific stock
func (s *Server) HandleGetStockAnalysis(w http.ResponseWriter, r *http.Request) {
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

	analysis, err := s.db.GetStockAnalysis(ctx, user.ID, stockID)
	if err != nil {
		http.Error(w, "Analysis not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		http.Error(w, "Could not encode analysis", http.StatusInternalServerError)
		return
	}
}

// HandleGetPortfolioAnalysis retrieves the latest analysis for the user's portfolio
func (s *Server) HandleGetPortfolioAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	// Get portfolio ID (optional - if not provided, get default portfolio)
	portfolioID := r.URL.Query().Get("portfolio_id")

	// If no portfolio ID provided, get the user's default portfolio
	if portfolioID == "" {
		portfolio, err := s.db.GetUserPortfolio(ctx, user.ID)
		if err != nil {
			http.Error(w, "Could not get portfolio", http.StatusInternalServerError)
			return
		}
		portfolioID = portfolio.ID
	}

	analysis, err := s.db.GetPortfolioAnalysis(ctx, user.ID, portfolioID)
	if err != nil {
		http.Error(w, "Portfolio analysis not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		http.Error(w, "Could not encode analysis", http.StatusInternalServerError)
		return
	}
}

// HandleGetAllStockAnalyses retrieves all stock analyses for the user
func (s *Server) HandleGetAllStockAnalyses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	analyses, err := s.db.GetUserStockAnalyses(ctx, user.ID)
	if err != nil {
		http.Error(w, "Could not retrieve analyses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(analyses); err != nil {
		http.Error(w, "Could not encode analyses", http.StatusInternalServerError)
		return
	}
}

// HandleDeleteStockAnalysis deletes a stock analysis
func (s *Server) HandleDeleteStockAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	analysisID := chi.URLParam(r, "id")
	if analysisID == "" {
		http.Error(w, "Analysis ID cannot be empty", http.StatusBadRequest)
		return
	}

	err := s.db.DeleteStockAnalysis(ctx, user.ID, analysisID)
	if err != nil {
		http.Error(w, "Could not delete analysis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message":     "Analysis deleted successfully",
		"analysis_id": analysisID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
		return
	}
}

// HandleDeletePortfolioAnalysis deletes a portfolio analysis
func (s *Server) HandleDeletePortfolioAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	analysisID := chi.URLParam(r, "id")
	if analysisID == "" {
		http.Error(w, "Analysis ID cannot be empty", http.StatusBadRequest)
		return
	}

	err := s.db.DeletePortfolioAnalysis(ctx, user.ID, analysisID)
	if err != nil {
		http.Error(w, "Could not delete analysis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message":     "Analysis deleted successfully",
		"analysis_id": analysisID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
		return
	}
}

// HandleGetAllPortfolioAnalyses retrieves all portfolio analyses for the user
func (s *Server) HandleGetAllPortfolioAnalyses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := middleware.User(ctx)
	if user == nil {
		http.Error(w, "Could not get user from context", http.StatusUnauthorized)
		return
	}

	analyses, err := s.db.GetUserPortfolioAnalyses(ctx, user.ID)
	if err != nil {
		http.Error(w, "Could not retrieve analyses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(analyses); err != nil {
		http.Error(w, "Could not encode analyses", http.StatusInternalServerError)
		return
	}
}
