package acg

import (
	"html/template"
	"net/http"

	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	tmpl *template.Template
	err  error
)

const postPerPage = 15

func init() {
	tmpl = template.Must(template.ParseGlob("internal/*/views/*.gohtml"))
}

func (s *Server) handleHomePage() http.HandlerFunc {
	type homepage struct {
		Page     *models.Page
		Services []*models.Service
		Posts    []*models.Post
	}

	return func(w http.ResponseWriter, r *http.Request) {
		page, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		services, err := s.store.Services().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[DEBUG] services: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		findOptions := options.Find()
		findOptions.SetLimit(3)
		findOptions.SetSort(bson.M{"time": -1})

		posts, err := s.store.Posts().Find(bson.M{"deleted": false}, findOptions)
		if err != nil {
			s.logger.Logf("[DEBUG] posts: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
			return
		}

		tmpl.ExecuteTemplate(w, "index.gohtml", &homepage{
			Page:     page,
			Services: services,
			Posts:    posts,
		})
	}
}

func (s *Server) handleAboutPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aboutpage, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		tmpl.ExecuteTemplate(w, "singlepage.gohtml", aboutpage)
	}
}

func (s *Server) handleMaterialsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, "This is materials page")
	}
}

func (s *Server) handleServicesPage() http.HandlerFunc {
	type servicepage struct {
		Page     *models.Page
		Services []*models.Service
	}

	return func(w http.ResponseWriter, r *http.Request) {
		page, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		services, err := s.store.Services().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[DEBUG] services: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		tmpl.ExecuteTemplate(w, "services.gohtml", &servicepage{
			Page:     page,
			Services: services,
		})
	}
}

func (s *Server) handleContactsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contactspage, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		tmpl.ExecuteTemplate(w, "singlepage.gohtml", contactspage)
	}
}
