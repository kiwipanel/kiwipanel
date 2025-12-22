package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kiwipanel/kiwipanel/config"
	admin "github.com/kiwipanel/kiwipanel/internal/modules/app/admin/routes"
)

func ProvidersRoutes(r *chi.Mux, appconfig *config.AppConfig) {

	admin.PublicRoutes(r, appconfig)
}
