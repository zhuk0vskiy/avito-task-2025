package server

import (
	"avito-task-2025/backend/config"
	"avito-task-2025/backend/internal/app"
	v1 "avito-task-2025/backend/internal/controller/http/v1"
	"fmt"

	// httpv1 "avito-task-2025/backend/internal/controller/http/v1"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	// "github.com/go-chi/cors"

	// "github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"

	// "github.com/rs/zerolog/log"
	"log"
)

type Server struct {
	Mux        *chi.Mux
	httpServer *http.Server
}

func NewServer(cfg config.HTTPConfig, jwtAuth *jwtauth.JWTAuth, a *app.App) *Server {
	mux := chi.NewMux()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))
	// mux.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"https://*", "http://*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: false,
	// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
	// }))

	mux.Use(middleware.Logger)

	mux.Route("/api", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Post("/auth", v1.SignInHandler(a))
		})

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwtAuth))
			// r.Use(jwtauth.Authenticator(jwtAuth))

			r.Get("/info", v1.GetUserInfoHandler(a))
			r.Post("/sendCoin", v1.SendCoinHandler(a))
			r.Get("/buy/{item}", v1.BuyMerchHandler(a))
		})
		
	})

	serverPort := fmt.Sprintf(":%s", cfg.Port)
	return &Server{
		httpServer: &http.Server{
			Addr:         serverPort,
			Handler:      mux,
			ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		},
	}
}

func (s *Server) Start() {
	log.Println("server start at port", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Println("error while starting server:", err)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("stopping server")
	return s.httpServer.Shutdown(ctx)
}
