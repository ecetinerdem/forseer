package api

import (
	"github.com/ecetinerdem/forseer/database"
	"github.com/ecetinerdem/forseer/middleware"
	services "github.com/ecetinerdem/forseer/service"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	db            *database.DB
	Router        *chi.Mux
	openAIService *services.OpenAIService
}

func NewServer(database *database.DB, openAIAPIKey string) *Server {
	s := &Server{
		db:            database,
		Router:        chi.NewRouter(),
		openAIService: services.NewOpenAIService(openAIAPIKey),
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

		// AI Analysis routes
		r.Route("/analysis", func(analysisRouter chi.Router) {
			analysisRouter.Use(middleware.UserAuthentication)

			// Stock analysis endpoints
			analysisRouter.Route("/stocks", func(stockAnalysisRouter chi.Router) {
				stockAnalysisRouter.Post("/{id}/analyze", s.HandleAnalyzeStock)  // Generate analysis for specific stock
				stockAnalysisRouter.Get("/{id}", s.HandleGetStockAnalysis)       // Get latest analysis for stock
				stockAnalysisRouter.Get("/", s.HandleGetAllStockAnalyses)        // Get all stock analyses for user
				stockAnalysisRouter.Delete("/{id}", s.HandleDeleteStockAnalysis) // Delete stock analysis
			})

			// Portfolio analysis endpoints
			analysisRouter.Route("/portfolio", func(portfolioAnalysisRouter chi.Router) {
				portfolioAnalysisRouter.Post("/analyze", s.HandleAnalyzePortfolio)       // Generate portfolio analysis
				portfolioAnalysisRouter.Get("/", s.HandleGetPortfolioAnalysis)           // Get latest portfolio analysis
				portfolioAnalysisRouter.Get("/all", s.HandleGetAllPortfolioAnalyses)     // Get all portfolio analyses for user
				portfolioAnalysisRouter.Delete("/{id}", s.HandleDeletePortfolioAnalysis) // Delete portfolio analysis
			})
		})
	})

	return s.Router
}
