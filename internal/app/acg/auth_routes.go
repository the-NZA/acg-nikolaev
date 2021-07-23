package acg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
)

func (s *Server) handleAuthRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{
			"1st RULE": "You do not talk about FIGHT CLUB.",
			"2nd RULE": "You DO NOT talk about FIGHT CLUB.",
		})
	}
}

func (s *Server) handleAuthLogin() http.HandlerFunc {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		cred := &req{}
		var err error

		if err = json.NewDecoder(r.Body).Decode(cred); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Logf("[DEBUG] %v with %v\n", cred.Username, cred.Password)

		if cred.Username == "" || cred.Password == "" {
			// s.logger.Logf("[ERROR] %v\n", http.StatusBadRequest)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoBodyParams)
			return
		}

		if err = s.store.Users().Login(cred.Username, cred.Password); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("User (%v) successfully auth'ed", cred.Username))
	}
}

func (s *Server) handleAuthLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{
			"logout": "this is placeholder",
		})
	}
}
