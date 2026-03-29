package server

import (
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/config"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/handler"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/service"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Router *chi.Mux
	pool   *pgxpool.Pool
}

func New(pool *pgxpool.Pool, cfg *config.Config) *Server {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(chimw.RequestID)

	// CORS for Next.js frontend
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Queries
	userQ := queries.NewUserQueries(pool)
	sessionQ := queries.NewSessionQueries(pool)
	categoryQ := queries.NewCategoryQueries(pool)
	itemQ := queries.NewItemQueries(pool)
	formQ := queries.NewFormQueries(pool)
	workflowQ := queries.NewWorkflowQueries(pool)

	// Services
	authSvc := service.NewAuthService(userQ, sessionQ)
	formSvc := service.NewFormService()
	workflowSvc := service.NewWorkflowService(workflowQ)

	// Handlers
	healthH := handler.NewHealthHandler()
	authH := handler.NewAuthHandler(authSvc)
	dashboardH := handler.NewDashboardHandler(itemQ, categoryQ, formQ, workflowQ)
	itemsH := handler.NewItemsHandler(itemQ, categoryQ)
	categoriesH := handler.NewCategoriesHandler(categoryQ)
	formsH := handler.NewFormsHandler(formQ, formSvc)
	workflowsH := handler.NewWorkflowsHandler(workflowQ, workflowSvc)

	s := &Server{Router: r, pool: pool}
	s.registerRoutes(r, authSvc, healthH, authH, dashboardH, itemsH, categoriesH, formsH, workflowsH)

	return s
}
