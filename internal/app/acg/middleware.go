package acg

import (
	"net/http"

	"github.com/the-NZA/acg-nikolaev/internal/app/auth"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
)

// authMiddleware check and varify cookie with token
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("TKN")

		if err != nil {
			switch err {
			case http.ErrNoCookie:
				s.logger.Logf("[ERROR] During cookie parse: %v\n", helpers.ErrUnauthorized)
				s.error(w, r, http.StatusUnauthorized, helpers.ErrUnauthorized)
				return
			default:
				s.logger.Logf("[ERROR] During cookie pares: %v\n", err)
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		// try to verify token and maybe go to 'next'
		isTokUpdate, err := auth.CheckToken(token.Value, s.config.SecretKey)
		if err != nil {
			s.logger.Logf("[ERROR] During token check: %v\n", err)
			s.error(w, r, http.StatusUnauthorized, helpers.ErrUnauthorized)
			return
		}

		// s.logger.Logf("[INFO] isTokUpdate: %v\n", isTokUpdate)
		if isTokUpdate {
			newToken, newExpTime, err := auth.UpdateToken(token.Value, s.config.SecretKey)
			if err != nil {
				s.logger.Logf("[ERROR] During token updated: %v\n", helpers.ErrUnauthorized)
				s.error(w, r, http.StatusUnauthorized, helpers.ErrUnauthorized)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "TKN",
				Value:    newToken,
				Expires:  newExpTime,
				HttpOnly: true,
				Path:     "/",
				Domain:   s.config.AppDomain,
			})
		}

		next.ServeHTTP(w, r)
	})
}
