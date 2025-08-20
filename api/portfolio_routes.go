package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) HandleGetPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	portfolio, err := s.db.GetStocks(ctx)

	if err != nil {
		http.Error(w, "Could not get portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&portfolio); err != nil {
		http.Error(w, "Could not encode portfolio", http.StatusInternalServerError)
	}
}

func (s *Server) HandleGetStockByID(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) HandleAddStockToPortfolio(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) HandleDeleteStockByID(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) HandleGetStockBySymbol(w http.ResponseWriter, r *http.Request) {

}
