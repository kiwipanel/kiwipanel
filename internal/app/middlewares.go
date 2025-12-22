package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kiwipanel/kiwipanel/pkg/helpers"
)

func MaintenanceMiddleware(metaPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := os.Stat(metaPath); err == nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(
					"KiwiPanel is temporarily unavailable (maintenance mode).",
				))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Verify passcode only for login page access
func PasscodeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passcode := chi.URLParam(r, "passcode")
		kiwipanelPasscode, _ := helpers.LoadGatePasscode()
		if passcode != kiwipanelPasscode {
			w.WriteHeader(http.StatusNotFound)

			notFound := `<!DOCTYPE html>
				<html style="height:100%">
				<head>
				<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
				<title> 404 Not Found
				</title><style>@media (prefers-color-scheme:dark){body{background-color:#000!important}}</style></head>
				<body style="color: #444; margin:0;font: normal 14px/20px Arial, Helvetica, sans-serif; height:100%; background-color: #fff;">
				<div style="height:auto; min-height:100%; ">     <div style="text-align: center; width:800px; margin-left: -400px; position:absolute; top: 30%; left:50%;">
						<h1 style="margin:0; font-size:150px; line-height:150px; font-weight:bold;">404</h1>
				<h2 style="margin-top:20px;font-size: 30px;">Not Found
				</h2>
				<p>The resource requested could not be found on this server!</p>
				</div></div></body></html>`

			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(notFound))

			//		w.Write([]byte("Not found"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Middlewares(r chi.Router) {
	metaPath := helpers.MaintenanceFlagPath()
	r.Use(MaintenanceMiddleware(metaPath))
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

}
