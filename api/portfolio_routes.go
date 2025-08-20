package api

import (
	"encoding/json"
	"net/http"

	"github.com/ecetinerdem/forseer/middleware"
	"github.com/go-chi/chi/v5"
)

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

	portfolio, err := s.db.GetStocks(ctx, user.ID)

	if err != nil {
		http.Error(w, "Could not get portfolio", http.StatusInternalServerError)
		return
	}

	if len(portfolio.Stocks) == 0 {
		http.Error(w, "There is no stocks in the portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&portfolio); err != nil {
		http.Error(w, "Could not encode portfolio", http.StatusInternalServerError)
	}
}

func (s *Server) HandleGetStockByID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	stockId := chi.URLParam(r, "id")

	if stockId == "" {
		http.Error(w, "Stock ID can not be empty", http.StatusBadRequest)
		return
	}

	stock, err := s.db.GetStockByID(ctx, stockId)

	if err != nil {
		http.Error(w, "Stock with given id does not exist", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&stock); err != nil {
		http.Error(w, "Cold not encode stock", http.StatusInternalServerError)
		return
	}

}

func (s *Server) HandleAddStockToPortfolio(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) HandleDeleteStockByID(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) HandleGetStockBySymbol(w http.ResponseWriter, r *http.Request) {

}
