package acg

import (
	"net/http"

	"github.com/the-NZA/acg-nikolaev/internal/app/auth"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("TKN")

		if err != nil {
			switch err {
			case http.ErrNoCookie:
				s.logger.Logf("[ERROR] %v\n", helpers.ErrUnauthorized)
				s.error(w, r, http.StatusUnauthorized, helpers.ErrUnauthorized)
				return
			default:
				s.logger.Logf("[ERROR] %v\n", err)
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		// s.logger.Logf("[DEBUG] %v\n", token.Value)

		// try to verify token and maybe go to 'next'
		if err = auth.CheckToken(token.Value, s.config.SecretKey); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusUnauthorized, helpers.ErrUnauthorized)
			return
		}

		// s.logger.Logf("[DEBUG] Auth'ed\n")
		next.ServeHTTP(w, r)
	})
}
