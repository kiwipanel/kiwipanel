package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kiwipanel/kiwipanel/config"
	"github.com/kiwipanel/kiwipanel/internal/modules/app/admin/controllers"
)

func PublicRoutes(r *chi.Mux, appconfig *config.AppConfig) {
	controller := controllers.New(appconfig)
	r.Get("/admin", controller.AdminHompage)

}
