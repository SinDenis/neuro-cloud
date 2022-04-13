package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"log"
	"nc-platform-back/internal/config"
	"nc-platform-back/internal/handler"
	"nc-platform-back/internal/interceptor"
	"net/http"
)

type Server struct {
	config          *config.Config
	router          *chi.Mux
	authHandler     *handler.AuthHandler
	registerHandler *handler.RegisterHandler
	imageHandler    *handler.ImageHandler
}

func (s *Server) configureRoutes() {
	s.router = chi.NewRouter()

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.URLFormat)
	s.router.Use(cors.AllowAll().Handler)

	s.router.Group(func(r chi.Router) {
		tokenAuth := jwtauth.New(s.config.JwtAlgo, []byte(s.config.JwtKey), nil)
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(interceptor.UserAuthMiddleware)
		r.Mount("/api/images", s.imageHandler.Routes())
	})

	s.router.Group(func(r chi.Router) {
		r.Mount("/api/login", s.authHandler.Routes())
		r.Mount("/api/register", s.registerHandler.Routes())
	})
}

func (s *Server) Run() {
	s.configureRoutes()
	log.Fatal(http.ListenAndServe(":8080", s.router))
}

func NewServer(
	config *config.Config,
	authHandler *handler.AuthHandler,
	registerHandler *handler.RegisterHandler,
	imageHandler *handler.ImageHandler,
) *Server {
	return &Server{
		config:          config,
		router:          chi.NewRouter(),
		authHandler:     authHandler,
		registerHandler: registerHandler,
		imageHandler:    imageHandler,
	}
}
