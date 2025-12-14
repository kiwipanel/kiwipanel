package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kiwipanel/kiwipanel/internal/modules/app/home/controllers"
)

func PublicRoutes(r *chi.Mux) {
	controller := controllers.New()
	r.Get("/", controller.Homepage)
	r.Get("/hello", controller.Homepage)

}
