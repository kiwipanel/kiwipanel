package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kiwipanel/kiwipanel/config"
	admin "github.com/kiwipanel/kiwipanel/internal/modules/app/admin/routes"
	public "github.com/kiwipanel/kiwipanel/internal/modules/app/home/routes"
)

func ProvidersRoutes(r *chi.Mux, appconfig *config.AppConfig) {
	public.PublicRoutes(r)
	admin.PublicRoutes(r, appconfig)
}
