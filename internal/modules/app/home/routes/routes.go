package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kiwipanel/kiwipanel/internal/modules/app/home/controllers"
)

func PublicRoutes(r *chi.Mux) {
	controller := controllers.New()
	r.Get("/", controller.Homepage)
	r.Get("/hello", controller.Homepage)
	r.Get("/health", HealthHandler)

}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
