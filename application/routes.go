package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/malav4all/golang-api/handler"
	"github.com/malav4all/golang-api/repository/alert"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/alert", a.loadAlertRoutes)

	a.router = router
}

func (a *App) loadAlertRoutes(router chi.Router) {
	alertHandler := &handler.Alert{
		Repo: &alert.RedisRepo{
			Client: a.rdb,
		},
	}

	router.Post("/", alertHandler.Create)
	router.Get("/", alertHandler.ListAlert)
	router.Get("/{id}", alertHandler.GetAlertId)
	router.Put("/{id}", alertHandler.UpdateById)
	router.Delete("/{id}", alertHandler.DeleteById)
}
