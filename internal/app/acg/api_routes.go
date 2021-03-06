package acg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const maxFileSize = 32 << 20 // Set max upload file size to 32MB

/*
 * Response helpers
 */
// respond method manage response with json encoding and optional data
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
 * Upload handler
 */
func (s *Server) handleUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(maxFileSize)

		uFile, uHeader, err := r.FormFile("acg_upload")
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		defer uFile.Close()

		// Match MIME types
		isImage, _ := regexp.MatchString(`image\/(jpg|png|jpeg|webp|gif|svg\+xml)`, uHeader.Header["Content-Type"][0])
		subFolder := "documents/"
		if isImage {
			subFolder = "images/"
		}

		// Generate filename with specific format
		suf := time.Now().Format("02-01-2006_15-04-05")
		uploadPath := "uploads/" + subFolder + suf + "_" + strings.ReplaceAll(uHeader.Filename, " ", "_")

		f, err := os.OpenFile(uploadPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			s.logger.Logf("[ERROR] during new file creation %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer f.Close()

		// Copy file bytes into file
		bytesWritten, err := io.Copy(f, uFile)
		if err != nil {
			s.logger.Logf("[ERROR] during new file creation %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Logf("[DEBUG] File %s uploaded in %s. Size %v\n", uHeader.Filename, f.Name(), bytesWritten)

		s.respond(w, r, http.StatusOK, uploadPath)
	}
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

func (s *Server) handleCategoryUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		category := &models.Category{}

		if err = json.NewDecoder(r.Body).Decode(category); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.Categories().Update(category); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Category (%v) successfully updated", category.ID.Hex()))
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

func (s *Server) handleCategoryGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("ID")

		if ID == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		objID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrInvalidObjectID)
			s.error(w, r, http.StatusBadRequest, helpers.ErrInvalidObjectID)
			return
		}

		category, err := s.store.Categories().FindByID(objID)

		switch err {
		case mongo.ErrNoDocuments:
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoCategory)
			s.error(w, r, http.StatusNotFound, helpers.ErrNoCategory)
			return
		case nil:
			s.respond(w, r, http.StatusOK, category)
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

		if req.ID.IsZero() {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrEmptyObjectID)
			s.error(w, r, http.StatusInternalServerError, helpers.ErrEmptyObjectID)
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

		_, err = s.store.Categories().FindByID(post.CategoryID)
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

func (s *Server) handlePostGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("ID")

		if ID == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		objID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrInvalidObjectID)
			s.error(w, r, http.StatusBadRequest, helpers.ErrInvalidObjectID)
			return
		}

		post, err := s.store.Posts().FindByID(objID)

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

// handlePostUpdate update post by it id
// * NOTE: must receive whole post struct
func (s *Server) handlePostUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		post := &models.Post{}

		if err = json.NewDecoder(r.Body).Decode(post); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.Posts().Update(post); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Post (%v) successfully updated", post.ID.Hex()))
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

		if req.ID.IsZero() {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrEmptyObjectID)
			s.error(w, r, http.StatusInternalServerError, helpers.ErrEmptyObjectID)
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
		var (
			err    error
			val    int64
			strVal string
		)

		findOpts := options.Find()
		findOpts.SetSort(bson.D{{Key: "time", Value: -1}})

		if r.URL.Query().Has("limit") {
			strVal = r.URL.Query().Get("limit")
			val, err = strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				s.logger.Logf("[DEBUG] during parse limit: %v\n", err)
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			findOpts.SetLimit(val)
		}

		if r.URL.Query().Has("skip") {
			strVal = r.URL.Query().Get("skip")
			val, err = strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				s.logger.Logf("[DEBUG] during parse skip: %v\n", err)
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			findOpts.SetSkip(val)
		}

		posts, err := s.store.Posts().Find(bson.M{"deleted": false}, findOpts)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, posts)
	}
}

func (s *Server) handlePostCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		numOfPosts, err := s.store.Posts().Count(bson.D{{Key: "deleted", Value: false}})
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]int64{"count": numOfPosts})
	}
}

/*
 * Post handlers END
 */

/*
 * Service handlers
 */
func (s *Server) handleServiceCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := &models.Service{
			ID: primitive.NewObjectID(),
		}

		var err error

		if err = json.NewDecoder(r.Body).Decode(service); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		service.Slug = helpers.GenerateSlug(service.Title)

		if err = s.store.Services().Create(service); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, fmt.Sprintf("Service (%s) successfully created", service.ID.Hex()))
	}
}

func (s *Server) handleServiceGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("ID")

		if ID == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		objID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrInvalidObjectID)
			s.error(w, r, http.StatusBadRequest, helpers.ErrInvalidObjectID)
			return
		}

		service, err := s.store.Services().FindByID(objID)

		switch err {
		case mongo.ErrNoDocuments:
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoService)
			s.error(w, r, http.StatusNotFound, helpers.ErrNoService)
			return
		case nil:
			s.respond(w, r, http.StatusOK, service)
			return
		default:
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *Server) handleServiceUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		service := &models.Service{}

		if err = json.NewDecoder(r.Body).Decode(service); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.Services().Update(service); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Service (%v) successfully updated", service.ID.Hex()))
	}
}

func (s *Server) handleServiceDelete() http.HandlerFunc {
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

		if req.ID.IsZero() {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrEmptyObjectID)
			s.error(w, r, http.StatusInternalServerError, helpers.ErrEmptyObjectID)
			return
		}

		if err = s.store.Services().Delete(req.ID); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Service (%s) successfully deleted", req.ID.Hex()))
	}
}

func (s *Server) handleServiceGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := s.store.Services().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, services)
	}
}

/*
 * Service handlers END
 */

/*
 * MatCategory handlers
 */

func (s *Server) handleMatCategoryCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matcat := &models.MatCategory{
			ID: primitive.NewObjectID(),
		}

		var err error

		if err = json.NewDecoder(r.Body).Decode(matcat); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		matcat.Slug = helpers.GenerateSlug(matcat.Title)

		if err = s.store.MatCategories().Create(matcat); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, fmt.Sprintf("Material category (%s) successfully created", matcat.ID.Hex()))

	}
}

func (s *Server) handleMatCategoryGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("ID")

		if ID == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		objID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrInvalidObjectID)
			s.error(w, r, http.StatusBadRequest, helpers.ErrInvalidObjectID)
			return
		}

		matcategory, err := s.store.MatCategories().FindByID(objID)

		switch err {
		case mongo.ErrNoDocuments:
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoMatCategory)
			s.error(w, r, http.StatusNotFound, helpers.ErrNoMatCategory)
			return
		case nil:
			s.respond(w, r, http.StatusOK, matcategory)
			return
		default:
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *Server) handleMatCategoryGetBySlug() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Query().Get("slug")

		if slug == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		matcat, err := s.store.MatCategories().FindBySlug(slug)

		switch err {
		case mongo.ErrNoDocuments:
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoMatCategory)
			s.error(w, r, http.StatusNotFound, helpers.ErrNoMatCategory)
			return
		case nil:
			s.respond(w, r, http.StatusOK, matcat)
			return
		default:
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *Server) handleMatCategoryUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		matcategory := &models.MatCategory{}

		if err = json.NewDecoder(r.Body).Decode(matcategory); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.MatCategories().Update(matcategory); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("MatCategory (%v) successfully updated", matcategory.ID.Hex()))
	}
}
func (s *Server) handleMatCategoryDelete() http.HandlerFunc {
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

		if req.ID.IsZero() {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrEmptyObjectID)
			s.error(w, r, http.StatusInternalServerError, helpers.ErrEmptyObjectID)
			return
		}

		if err = s.store.MatCategories().Delete(req.ID); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Material category (%s) successfully deleted", req.ID.Hex()))
	}
}

func (s *Server) handleMatCategoryGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := s.store.MatCategories().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, services)
	}
}

/*
 * MatCategory handlers END
 */

/*
 * Material handlers
 */

func (s *Server) handleMaterialCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		material := &models.Material{
			ID:   primitive.NewObjectID(),
			Time: time.Now(),
		}

		var err error

		if err = json.NewDecoder(r.Body).Decode(material); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		_, err = s.store.MatCategories().FindByID(material.MatCategoryID)
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

		material.Slug = helpers.GenerateSlug(material.Title)

		if err = s.store.Materials().Create(material); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, material.ID.Hex())
	}
}

func (s *Server) handleMaterialGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("ID")

		if ID == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		objID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrInvalidObjectID)
			s.error(w, r, http.StatusBadRequest, helpers.ErrInvalidObjectID)
			return
		}

		material, err := s.store.Materials().FindByID(objID)

		switch err {
		case mongo.ErrNoDocuments:
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoMaterial)
			s.error(w, r, http.StatusNotFound, helpers.ErrNoMaterial)
			return
		case nil:
			s.respond(w, r, http.StatusOK, material)
			return
		default:
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

// handleMaterialUpdate update post by it id
// * NOTE: must receive whole material struct
func (s *Server) handleMaterialUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		material := &models.Material{}

		if err = json.NewDecoder(r.Body).Decode(material); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.Materials().Update(material); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Post (%v) successfully updated", material.ID.Hex()))
	}
}
func (s *Server) handleMaterialDelete() http.HandlerFunc {
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

		if req.ID.IsZero() {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrEmptyObjectID)
			s.error(w, r, http.StatusInternalServerError, helpers.ErrEmptyObjectID)
			return
		}

		if err = s.store.Materials().Delete(req.ID); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Material (%s) successfully deleted", req.ID.Hex()))
	}
}

func (s *Server) handleMaterialGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err    error
			val    int64
			strVal string
		)

		findOpts := options.Find()
		findOpts.SetSort(bson.D{{Key: "time", Value: -1}})

		if r.URL.Query().Has("limit") {
			strVal = r.URL.Query().Get("limit")
			val, err = strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				s.logger.Logf("[DEBUG] during parse limit: %v\n", err)
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			findOpts.SetLimit(val)
		}

		if r.URL.Query().Has("skip") {
			strVal = r.URL.Query().Get("skip")
			val, err = strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				s.logger.Logf("[DEBUG] during parse skip: %v\n", err)
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			findOpts.SetSkip(val)
		}

		materials, err := s.store.Materials().Find(bson.M{"deleted": false}, findOpts)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, materials)
	}
}

func (s *Server) handleMaterialGetAllBySlug() http.HandlerFunc {
	type matCat struct {
		*models.MatCategory
		Materials []*models.Material `json:"materials,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		matcat, err := s.store.MatCategories().FindBySlug(chi.URLParam(r, "matCatSlug"))
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		materials, err := s.store.Materials().FindAll(bson.M{"matcategory_id": matcat.ID, "deleted": false})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, &matCat{
			MatCategory: matcat,
			Materials:   materials,
		})
	}
}

func (s *Server) handleMaterialCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		numOfMaterials, err := s.store.Materials().Count(bson.D{{Key: "deleted", Value: false}})
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]int64{"count": numOfMaterials})
	}
}

/*
 * Material handlers END
 */

/*
 * Page handlers
 */
func (s *Server) handlePageCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := &models.Page{
			ID: primitive.NewObjectID(),
		}

		var err error

		if err = json.NewDecoder(r.Body).Decode(page); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		page.URL = "/" + helpers.GenerateSlug(page.Title)

		if err = s.store.Pages().Create(page); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, fmt.Sprintf("Page (%s) successfully created", page.ID.Hex()))
	}
}

func (s *Server) handlePageGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("ID")

		if ID == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		objID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrInvalidObjectID)
			s.error(w, r, http.StatusBadRequest, helpers.ErrInvalidObjectID)
			return
		}

		page, err := s.store.Pages().FindByID(objID)

		switch err {
		case mongo.ErrNoDocuments:
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoService)
			s.error(w, r, http.StatusNotFound, helpers.ErrNoService)
			return
		case nil:
			s.respond(w, r, http.StatusOK, page)
			return
		default:
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}
func (s *Server) handlePageGetByURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")

		if url == "" {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoRequestParams)
			s.error(w, r, http.StatusBadRequest, helpers.ErrNoRequestParams)
			return
		}

		page, err := s.store.Pages().FindByURL(url)

		switch err {
		case mongo.ErrNoDocuments:
			s.logger.Logf("[ERROR] %v\n", helpers.ErrNoPage)
			s.error(w, r, http.StatusNotFound, helpers.ErrNoPage)
			return
		case nil:
			s.respond(w, r, http.StatusOK, page)
			return
		default:
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

// handlePageUpdate update page by it id
// * NOTE: must receive whole page struct
func (s *Server) handlePageUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		page := &models.Page{}

		if err = json.NewDecoder(r.Body).Decode(page); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.Pages().Update(page); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Page (%v) successfully updated", page.ID.Hex()))
	}
}

func (s *Server) handlePageDelete() http.HandlerFunc {
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

		if req.ID.IsZero() {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrEmptyObjectID)
			s.error(w, r, http.StatusInternalServerError, helpers.ErrEmptyObjectID)
			return
		}

		if err = s.store.Pages().Delete(req.ID); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("Page (%s) successfully deleted", req.ID.Hex()))
	}
}

func (s *Server) handlePageGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages, err := s.store.Pages().FindAll(bson.M{})
		if err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, pages)
	}
}

/*
 * Page handlers END
 */

/*
 * User handlers
 */
func (s *Server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := &models.User{
			ID: primitive.NewObjectID(),
		}

		var err error

		if err = json.NewDecoder(r.Body).Decode(usr); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.Users().Create(usr); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, fmt.Sprintf("User '%s' successfully created", usr.Username))
	}
}

func (s *Server) handleUserDelete() http.HandlerFunc {
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

		if req.ID.IsZero() {
			s.logger.Logf("[ERROR] %v\n", helpers.ErrEmptyObjectID)
			s.error(w, r, http.StatusInternalServerError, helpers.ErrEmptyObjectID)
			return
		}

		if err = s.store.Users().Delete(req.ID); err != nil {
			s.logger.Logf("[ERROR] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, fmt.Sprintf("User (%s) successfully deleted", req.ID.Hex()))
	}
}

/*
 * User handlers END
 */
