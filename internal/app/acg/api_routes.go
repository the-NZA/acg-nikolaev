package acg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
 * Response helpers
 */

// respond method manage response with json encoding end optional data
func (s Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// error method manage response with error with wrapping it
func (s Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

/*
 * Categories handlers
 */

func (s *Server) handleCategoryCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cat := &models.Category{
			ID: primitive.NewObjectID(),
		}

		var err error

		if err = json.NewDecoder(r.Body).Decode(cat); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		cat.Slug = helpers.GenerateSlug(cat.Title)

		s.logger.Logf("[DEBUG] before\n")

		if err = s.store.Categories().Create(cat); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Logf("[DEBUG] After create\n")

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Category (%s) successfully created", cat.ID))

	}
}

func (s *Server) handleCategoryGetBySlug() http.HandlerFunc {
	type req struct {
		Slug string `json:"slug"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
			req *req
		)

		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		cat, err := s.store.Categories().FindBySlug(req.Slug)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, cat)

	}
}
