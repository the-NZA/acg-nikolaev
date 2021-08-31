package acg

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-pkgz/lgr"
	"github.com/the-NZA/acg-nikolaev/internal/app/store"
	"github.com/the-NZA/acg-nikolaev/internal/app/store/mongostore"
)

// Server contains all things to run website
type Server struct {
	config *Config
	logger *lgr.Logger
	router *chi.Mux
	store  store.Storer
}

// NewServer returns Server object with router, logger and config
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		router: chi.NewRouter(),
	}
}

func (s *Server) configureRouter() {
	// CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// CORS END

	// Pages Routes
	s.router.Get("/", s.handleHomePage())

	s.router.Get("/about", s.handleAboutPage())

	// s.router.Get("/materials")

	s.router.Route("/materials", func(r chi.Router) {
		r.Get("/", s.handleMaterialsPage())
		r.Get("/{matCatSlug:[a-z0-9_-]+}", s.handleMaterialGetAllBySlug())
	})

	s.router.Get("/services", s.handleServicesPage())

	s.router.Get("/contacts", s.handleContactsPage())

	s.router.Route("/posts", func(r chi.Router) {
		r.Get("/", s.handlePostsPage())
	})

	s.router.Route("/category", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			// * NOTE: Just redirect to /posts page for now
			http.Redirect(w, r, "/posts", http.StatusSeeOther)
		})

		r.Get("/{categorySlug:[a-z0-9_-]+}", s.handleSingleCategoryPage())

		r.Get("/{categorySlug:[a-z0-9_-]+}/{postSlug:[a-z0-9_-]+}", s.handleSinglePostPage())
	})

	s.router.Get("/404", func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusNotFound, map[string]string{
			"page": "not found",
			"you":  "must try another one",
		})

	})
	// Pages END

	// Not Found
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Logf("[DEBUG] 404 at %v\n", r.URL.Path)
		s.respond(w, r, http.StatusNotFound, map[string]string{
			"page": "not found",
			"you":  "must try another one",
		})
	})
	// Not Found END

	// API Routes
	s.router.Route("/api", func(r chi.Router) {
		// ! REMOVE BEFORE GOING LIVE
		if !s.config.LogDebug {
			r.Use(s.authMiddleware)
		}

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			s.respond(w, r, http.StatusOK, "This is API endpoint")
		})

		r.Post("/upload", s.handleUpload())

		r.Route("/category", func(r chi.Router) {
			r.Get("/", s.handleCategoryGetByID())
			r.Post("/", s.handleCategoryCreate())
			r.Put("/", s.handleCategoryUpdate())
			r.Delete("/", s.handleCategoryDelete())
			r.Get("/all", s.handleCategoryGetAll())
		})

		r.Route("/post", func(r chi.Router) {
			r.Get("/", s.handlePostGetByID())
			r.Post("/", s.handlePostCreate())
			r.Delete("/", s.handlePostDelete())
			r.Put("/", s.handlePostUpdate())
			r.Get("/all", s.handlePostGetAll())
			r.Get("/count", s.handlePostCount())
		})

		r.Route("/service", func(r chi.Router) {
			r.Get("/", s.handleServiceGetByID())
			r.Post("/", s.handleServiceCreate())
			r.Put("/", s.handleServiceUpdate())
			r.Delete("/", s.handleServiceDelete())
			r.Get("/all", s.handleServiceGetAll())
		})

		r.Route("/matcategory", func(r chi.Router) {
			r.Post("/", s.handleMatCategoryCreate())
			r.Get("/", s.handleMatCategoryGetByID())
			r.Delete("/", s.handleMatCategoryDelete())
			r.Put("/", s.handleMatCategoryUpdate())
			r.Get("/all", s.handleMatCategoryGetAll())
		})

		r.Route("/material", func(r chi.Router) {
			r.Post("/", s.handleMaterialCreate())
			r.Get("/", s.handleMaterialGetByID())
			r.Delete("/", s.handleMaterialDelete())
			r.Put("/", s.handleMaterialUpdate())
			r.Get("/all", s.handleMaterialGetAll())
			r.Get("/count", s.handleMaterialCount())
		})

		r.Route("/page", func(r chi.Router) {
			// r.Get("/", s.handlePageGetByURL())
			r.Get("/", s.handlePageGetByID())
			r.Post("/", s.handlePageCreate())
			r.Delete("/", s.handlePageDelete())
			r.Put("/", s.handlePageUpdate())
			r.Get("/all", s.handlePageGetAll())
		})

		r.Route("/user", func(r chi.Router) {
			r.Post("/", s.handleUserCreate())
			r.Delete("/", s.handleUserDelete())
		})
	})
	// API END

	// Auth Routes
	s.router.Route("/auth", func(r chi.Router) {
		r.Get("/", s.handleAuthRoot())
		r.Post("/login", s.handleAuthLogin())
		r.Post("/logout", s.handleAuthLogout())
	})
	// Auth END
}

// configureStore creates new Store and try to establish connection
func (s *Server) configureStore() error {
	st, err := mongostore.NewStore(s.config.DatabaseURL)
	if err != nil {
		return err
	}

	s.store = st
	return nil
}

// newLogger configure logger in DEBUG or PRODUCTION mode
// Possible log levels TRACE, DEBUG, INFO, WARN, ERROR, PANIC and FATAL
func (s *Server) configureLogger(dbg bool) *lgr.Logger {
	if dbg {
		return lgr.New(lgr.Msec, lgr.Debug, lgr.CallerFile, lgr.CallerFunc, lgr.LevelBraces)
	}

	return lgr.New(lgr.Msec, lgr.LevelBraces)
}

// Start performs pre-run configuration and starts server
func (s *Server) Start() error {
	s.logger = s.configureLogger(s.config.LogDebug)

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	s.logger.Logf("[INFO] Server is starting at %v...\n", s.config.BindAddr)

	return http.ListenAndServe(s.config.BindAddr, s.router)
}
