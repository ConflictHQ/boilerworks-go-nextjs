package server

import (
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/handler"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/service"
	"github.com/go-chi/chi/v5"
)

func (s *Server) registerRoutes(
	r *chi.Mux,
	authSvc *service.AuthService,
	healthH *handler.HealthHandler,
	authH *handler.AuthHandler,
	dashboardH *handler.DashboardHandler,
	itemsH *handler.ItemsHandler,
	categoriesH *handler.CategoriesHandler,
	formsH *handler.FormsHandler,
	workflowsH *handler.WorkflowsHandler,
) {
	// Health check (no auth)
	r.Get("/health", healthH.Health)

	// Public auth routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", authH.Login)
		r.Post("/register", authH.Register)
		r.Post("/logout", authH.Logout)
	})

	// Authenticated API routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authSvc))

		r.Get("/api/auth/me", authH.Me)
		r.Get("/api/dashboard", dashboardH.Dashboard)

		// Items
		r.Route("/api/items", func(r chi.Router) {
			r.With(middleware.RequirePermission("items.view")).Get("/", itemsH.List)
			r.With(middleware.RequirePermission("items.view")).Get("/{uuid}", itemsH.Get)
			r.With(middleware.RequirePermission("items.create")).Post("/", itemsH.Create)
			r.With(middleware.RequirePermission("items.edit")).Put("/{uuid}", itemsH.Update)
			r.With(middleware.RequirePermission("items.delete")).Delete("/{uuid}", itemsH.Delete)
		})

		// Categories
		r.Route("/api/categories", func(r chi.Router) {
			r.With(middleware.RequirePermission("categories.view")).Get("/", categoriesH.List)
			r.With(middleware.RequirePermission("categories.view")).Get("/{uuid}", categoriesH.Get)
			r.With(middleware.RequirePermission("categories.create")).Post("/", categoriesH.Create)
			r.With(middleware.RequirePermission("categories.edit")).Put("/{uuid}", categoriesH.Update)
			r.With(middleware.RequirePermission("categories.delete")).Delete("/{uuid}", categoriesH.Delete)
		})

		// Forms
		r.Route("/api/forms", func(r chi.Router) {
			r.With(middleware.RequirePermission("forms.view")).Get("/", formsH.ListDefinitions)
			r.With(middleware.RequirePermission("forms.view")).Get("/{uuid}", formsH.GetDefinition)
			r.With(middleware.RequirePermission("forms.create")).Post("/", formsH.CreateDefinition)
			r.With(middleware.RequirePermission("forms.edit")).Put("/{uuid}", formsH.UpdateDefinition)
			r.With(middleware.RequirePermission("forms.delete")).Delete("/{uuid}", formsH.DeleteDefinition)
			r.With(middleware.RequirePermission("forms.view")).Get("/{uuid}/submissions", formsH.ListSubmissions)
			r.With(middleware.RequirePermission("forms.create")).Post("/{uuid}/submissions", formsH.CreateSubmission)
		})

		// Workflows
		r.Route("/api/workflows", func(r chi.Router) {
			r.With(middleware.RequirePermission("workflows.view")).Get("/", workflowsH.ListDefinitions)
			r.With(middleware.RequirePermission("workflows.view")).Get("/{uuid}", workflowsH.GetDefinition)
			r.With(middleware.RequirePermission("workflows.create")).Post("/", workflowsH.CreateDefinition)
			r.With(middleware.RequirePermission("workflows.edit")).Put("/{uuid}", workflowsH.UpdateDefinition)
			r.With(middleware.RequirePermission("workflows.delete")).Delete("/{uuid}", workflowsH.DeleteDefinition)
			r.With(middleware.RequirePermission("workflows.view")).Get("/{uuid}/instances", workflowsH.ListInstances)
			r.With(middleware.RequirePermission("workflows.create")).Post("/{uuid}/instances", workflowsH.CreateInstance)
			r.With(middleware.RequirePermission("workflows.view")).Get("/instances/{uuid}", workflowsH.GetInstance)
			r.With(middleware.RequirePermission("workflows.edit")).Post("/instances/{uuid}/transition", workflowsH.TransitionInstance)
		})
	})
}
