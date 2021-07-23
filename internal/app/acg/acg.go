package acg

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	// Pages Routes
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("home page"))
	})

	s.router.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("about page"))
	})

	s.router.Get("/materials", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("materials page"))
	})

	s.router.Get("/services", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("services page"))
	})

	s.router.Get("/contacts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("contacts page"))
	})

	s.router.Route("/posts", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("posts"))
		})

		r.Get("/{postSlug:[a-z-]+}", func(w http.ResponseWriter, r *http.Request) {
			slug := chi.URLParam(r, "postSlug")

			s.logger.Logf("INFO %v\n", slug)

			w.Write([]byte(slug))

		})
	})

	s.router.Route("/category", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("category"))
		})

		r.Get("/{categorySlug:[a-z-]+}", func(w http.ResponseWriter, r *http.Request) {
			slug := chi.URLParam(r, "categorySlug")

			s.logger.Logf("INFO %v\n", slug)

			w.Write([]byte(slug))
		})
	})
	// Pages END

	// API Routes
	s.router.Route("/api", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// s.logger.Logf("[INFO] This is message from middleware\n")
				// s.logger.Logf("[DEBUG] %v \n", r.Body)
				cookie, err := r.Cookie("TKN")
				if err != nil {
					s.logger.Logf("[ERROR] %v\n", err)
					s.error(w, r, http.StatusInternalServerError, err)
					return
				}

				s.logger.Logf("[DEBUG] %v\n", cookie)

				// try to verify token and maybe go to 'next'

				next.ServeHTTP(w, r)
			})
		})

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is api"))
		})

		r.Route("/category", func(r chi.Router) {
			r.Get("/", s.handleCategoryGetBySlug())
			r.Post("/", s.handleCategoryCreate())
			r.Delete("/", s.handleCategoryDelete())
			r.Get("/all", s.handleCategoryGetAll())
		})

		r.Route("/post", func(r chi.Router) {
			r.Get("/", s.handlePostGetBySlug())
			r.Post("/", s.handlePostCreate())
			r.Delete("/", s.handlePostDelete())
			r.Get("/all", s.handlePostGetAll())
		})

		r.Route("/service", func(r chi.Router) {
			r.Get("/", s.handleServiceGetByID())
			r.Post("/", s.handleServiceCreate())
			r.Delete("/", s.handleServiceDelete())
			r.Get("/all", s.handleServiceGetAll())
		})

		r.Route("/matcategory", func(r chi.Router) {
			r.Post("/", s.handleMatCategoryCreate())
			r.Get("/", s.handleMatCategoryGetByID())
			r.Delete("/", s.handleMatCategoryDelete())
			r.Get("/all", s.handleMatCategoryGetAll())
		})

		r.Route("/material", func(r chi.Router) {
			r.Post("/", s.handleMaterialCreate())
			r.Get("/", s.handleMaterialGetByID())
			r.Delete("/", s.handleMaterialDelete())
			r.Get("/all", s.handleMaterialGetAll())
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
