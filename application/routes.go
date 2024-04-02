package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/malav4all/golang-api/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Route("/alert", loadAlertRoutes)

	return router
}

func loadAlertRoutes(router chi.Router) {
	alertHandler := &handler.Alert{}

	router.Post("/", alertHandler.Create)
	router.Get("/", alertHandler.ListAlert)
	// router.Get("/{id}",alertHandler.UpdateById)
	router.Put("/", alertHandler.UpdateById)
	router.Delete("/{id}", alertHandler.DeleteById)
}
