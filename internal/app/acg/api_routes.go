package acg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

		if err = s.store.Categories().Create(cat); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, fmt.Sprintf("Category (%s) successfully created", cat.ID.Hex()))

	}
}

func (s *Server) handleCategoryGetBySlug() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Query().Get("slug")

		if slug == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		cat, err := s.store.Categories().FindBySlug(slug)

		switch err {
		case mongo.ErrNoDocuments:
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoCategory)
			s.error(w, r, http.StatusNotFound, helpers.ErrNoCategory)
			return
		case nil:
			s.respond(w, r, http.StatusOK, cat)
			return
		default:
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *Server) handleCategoryGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cats, err := s.store.Categories().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, cats)
	}
}

func (s *Server) handleCategoryDelete() http.HandlerFunc {
	type req struct {
		ID primitive.ObjectID `json:"deletedID"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &req{}
		var err error

		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.Categories().Delete(req.ID); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Category (%s) successfully deleted", req.ID.Hex()))
	}
}

/*
 * Categories handlers END
 */

/*
 * Post handlers
 */
func (s *Server) handlePostCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		post := &models.Post{
			ID:   primitive.NewObjectID(),
			Time: time.Now(),
		}

		var err error

		if err = json.NewDecoder(r.Body).Decode(post); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		cat, err := s.store.Categories().FindByID(post.CategoryID)
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				s.logger.Logf("[ERROR] %v\n", helpers.ErrNoCategory)
				s.error(w, r, http.StatusNotFound, helpers.ErrNoCategory)
			default:
				s.logger.Logf("[ERROR] %v\n", err)
				s.error(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		post.Slug = helpers.GenerateSlug(post.Title)
		post.URL = cat.URL() + "/" + post.Slug

		// post.URL = fmt.Sprintf("/%s/%s", cat.Slug, helpers.GenerateSlug(post.Title))

		if err = s.store.Posts().Create(post); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, fmt.Sprintf("Post (%s) successfully created", post.ID.Hex()))
	}
}

func (s *Server) handlePostGetBySlug() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Query().Get("slug")

		if slug == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		post, err := s.store.Posts().FindBySlug(slug)

		switch err {
		case mongo.ErrNoDocuments:
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoPost)
			s.error(w, r, http.StatusNotFound, helpers.ErrNoPost)
			return
		case nil:
			s.respond(w, r, http.StatusOK, post)
			return
		default:
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *Server) handlePostDelete() http.HandlerFunc {
	type req struct {
		ID primitive.ObjectID `json:"deletedID"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &req{}
		var err error

		if err = json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.Posts().Delete(req.ID); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Post (%s) successfully deleted", req.ID.Hex()))
	}
}

func (s *Server) handlePostGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := s.store.Posts().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, posts)
	}
}

/*
 * Post handlers END
 */
