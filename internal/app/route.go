package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kiwipanel/kiwipanel/config"
	"github.com/kiwipanel/kiwipanel/internal/modules/providers/routes"
)

func NewRoutes(appconfig *config.AppConfig) http.Handler {
	r := chi.NewRouter()
	Middlewares(r)
	routes.ProvidersRoutes(r, appconfig)
	return r
}
