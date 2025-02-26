package auth

import (
	"github.com/kamchatkin/practicum-shortener/internal/auth"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"net/http"
)

// WithAuth Выдавать пользователю симметрично подписанную куку, содержащую уникальный идентификатор пользователя, если такой куки не существует или она не проходит проверку подлинности.
func WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger := logs.NewLogger()

		setNewCookie := false

		authCookie, errNoCookie := r.Cookie(auth.CookineName)
		if errNoCookie != nil {
			logger.Info("No cookie found")
			setNewCookie = true
		}

		if !setNewCookie {
			if auth.GetUserID(authCookie.Value) == -1 {
				logger.Info("No user_id found")
				setNewCookie = true
			}
		}

		if !setNewCookie {
			next.ServeHTTP(w, r)
			return
		}

		token, err := auth.BuildJWTString()
		if err != nil {
			logger.Info("Error building JWT token")
			next.ServeHTTP(w, r)
			return
		}

		// кука на чувашском дядя
		cuca := &http.Cookie{
			Name:  auth.CookineName,
			Value: token,
		}

		logger.Info("Setting new cookie")
		http.SetCookie(w, cuca)
		r.AddCookie(cuca)

		next.ServeHTTP(w, r)
	}
}
