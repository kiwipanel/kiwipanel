package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kiwipanel/kiwipanel/config"
	"github.com/kiwipanel/kiwipanel/internal/modules/providers/routes"
)

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>KiwiPanel Login</title>
		<style>
			body { font-family: Arial; background: #f0f0f0; margin: 0; padding: 20px; }
			.container { max-width: 400px; margin: 100px auto; background: white; padding: 30px; border-radius: 5px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
			h1 { color: #333; text-align: center; }
			input { width: 100%; padding: 10px; margin: 10px 0; box-sizing: border-box; border: 1px solid #ddd; border-radius: 4px; }
			button { width: 100%; padding: 10px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
			button:hover { background: #0056b3; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>KiwiPanel</h1>
			<form method="post" action="">
				<input type="text" name="username" placeholder="Username" required>
				<input type="password" name="password" placeholder="Password" required>
				<button type="submit">Login</button>
			</form>
		</div>
	</body>
	</html>
	`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate credentials
	if validateCredentials(username, password) {
		token := generateSessionToken()
		http.SetCookie(w, &http.Cookie{
			Name:     "kiwipanel_session",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400 * 7, // 7 days
		})

		// After login, redirect to dashboard (no passcode in URL)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Invalid credentials"))
}

// Placeholder functions - implement properly in production
func validateCredentials(username, password string) bool {
	// TODO: Use bcrypt to compare hashed passwords
	return username == "admin" && password == "kiwi12345"
}

func generateSessionToken() string {
	// TODO: Use crypto/rand for real token generation
	return "session_token_placeholder"
}

func isValidToken(token string) bool {
	// TODO: Validate against stored tokens in database/cache
	return token == "session_token_placeholder"
}

func NewRoutes(appconfig *config.AppConfig) http.Handler {
	r := chi.NewRouter()
	Middlewares(r)

	r.Route("/{passcode}", func(r chi.Router) {
		r.Use(PasscodeMiddleware)
		r.Get("/", handleLoginPage)
		r.Post("/login", handleLoginSubmit)
	})

	routes.ProvidersRoutes(r, appconfig)

	r.NotFound(http.HandlerFunc(NewNotFoundHandler))
	return r
}

func NewNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Page not found at KiwiPanel"))

}
