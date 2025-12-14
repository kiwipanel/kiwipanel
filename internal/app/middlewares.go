package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kiwipanel/kiwipanel/pkg/helpers"
)

const hiddenFlagProduction = "/opt/kiwipanel/meta/hidden"
const hiddenFlagDev = "kiwipanel/meta/hidden"

var metaPath string

func MaintenanceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if helpers.IsInstalled() {
			metaPath = hiddenFlagProduction
		} else {
			metaPath = hiddenFlagDev
		}
		if _, err := os.Stat(metaPath); err == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("KiwiPanel is temporarily unavailable (maintenance mode)."))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Middlewares(r chi.Router) {
	r.Use(MaintenanceMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

}
