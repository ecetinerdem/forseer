package api

import (
	"github.com/ecetinerdem/forseer/database"
	"github.com/ecetinerdem/forseer/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	db     *database.DB
	Router *chi.Mux
}

func NewServer(database *database.DB) *Server {
	s := &Server{
		db:     database,
		Router: chi.NewRouter(),
	}
	s.setUpRoutes()
	return s
}

func (s *Server) setUpRoutes() *chi.Mux {
	// Public routes
	s.Router.Get("/", s.HandleGreeting)
	s.Router.Post("/register", s.HandleCreateUser)
	s.Router.Post("/login", s.HandleLoginUser)

	// API v1 routes
	s.Router.Route("/api/v1", func(r chi.Router) {
		// User routes
		r.Route("/users", func(userRouter chi.Router) {
			userRouter.Use(middleware.UserAuthentication)
			userRouter.Get("/", s.HandleGetUsers)
			userRouter.Get("/search", s.HandleGetUserByEmail) // Changed to use query param
			userRouter.Get("/{id}", s.HandleGetUserById)
			userRouter.Put("/{id}", s.HandleUpdateUser)
			userRouter.Delete("/{id}", s.HandleDeleteUserById)
		})

		// Portfolio routes
		r.Route("/portfolio", func(portfolioRouter chi.Router) {
			portfolioRouter.Use(middleware.UserAuthentication)

			// Portfolio operations
			portfolioRouter.Get("/", s.HandleGetPortfolio)
			portfolioRouter.Post("/", s.HandleCreatePortfolio) // For creating new portfolios

			// Stock operations
			portfolioRouter.Route("/stocks", func(stockRouter chi.Router) {
				stockRouter.Get("/", s.HandleGetUserStocks)                // Get all stocks
				stockRouter.Post("/{symbol}", s.HandleAddStockToPortfolio) // Add stock by symbol
				stockRouter.Get("/search", s.HandleGetStockBySymbol)       // Search stocks by symbol (query param)
				stockRouter.Get("/{id}", s.HandleGetStockByID)             // Get specific stock
				stockRouter.Delete("/{id}", s.HandleDeleteStockByID)       // Delete stock
			})
		})
	})

	return s.Router
}
