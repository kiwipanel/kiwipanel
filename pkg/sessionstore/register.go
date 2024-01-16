package sessionstore

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Register(r *echo.Echo) {
	//https://github.com/CurtisVermeeren/gorilla-sessions-tutorial/blob/master/CookieStoreSession/main.go

	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	Store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	ss := session.Middleware(Store)
	r.Use(ss)
}
