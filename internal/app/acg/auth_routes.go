package acg

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
)

// handleAuthRoot just a placeholder
func (s *Server) handleAuthRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{
			"1st RULE": "You do not talk about FIGHT CLUB.",
			"2nd RULE": "You DO NOT talk about FIGHT CLUB.",
		})
	}
}

// handleAuthLogin read username and password from Body and try to authorize user
func (s *Server) handleAuthLogin() http.HandlerFunc {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		cred := &req{}
		var err error

		if err = json.NewDecoder(r.Body).Decode(cred); err != nil {
			s.logger.Logf("[ERROR] During decode body: %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if cred.Username == "" || cred.Password == "" {
			s.logger.Logf("[ERROR] Empty credentials in body: %v\n", helpers.ErrNoBodyParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoBodyParams)
			return
		}

		token, expTime, err := s.store.Users().Login(cred.Username, cred.Password, s.config.SecretKey)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "TKN",
			Value:    token,
			Expires:  expTime,
			HttpOnly: true,
			Path:     "/",
			Domain:   s.config.AppDomain,
		})

		s.respond(w, r, http.StatusOK, map[string]string{
			"login": "successful",
			// "user":  cred.Username,
			"token": token,
		})
	}
}

// handleAuthLogout remove cookie with token
func (s *Server) handleAuthLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "TKN",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Path:     "/",
			Domain:   s.config.AppDomain,
		})

		s.respond(w, r, http.StatusOK, map[string]string{
			"logout": "successful",
		})
	}
}
