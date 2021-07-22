package acg

import "net/http"

func (s *Server) handleAuthRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{
			"1st RULE": "You do not talk about FIGHT CLUB.",
			"2nd RULE": "You DO NOT talk about FIGHT CLUB.",
		})
	}
}

func (s *Server) handleAuthLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{
			"login": "this is placeholder",
		})
	}
}

func (s *Server) handleAuthLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{
			"logout": "this is placeholder",
		})
	}
}
