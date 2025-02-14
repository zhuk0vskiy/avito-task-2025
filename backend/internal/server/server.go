package server

import (
	"avito-task-2025/backend/config"
	"avito-task-2025/backend/internal/app"
	"avito-task-2025/backend/internal/controller"
	"fmt"

	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"

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
	handler := gin.New()
	con := controller.NewRouter(handler)

	// Set routes
	con.SetV1Routes(
		a.Logger,
		a.AuthSvcIntf,
		a.CoinSvcIntf,
		a.MerchSvcIntf,
		a.UserSvcIntf,
		a.JwtMngIntf,
	)

	// fmt.Println(cfg.Port)
	serverPort := fmt.Sprintf(":%s", cfg.Port)
	return &Server{
		httpServer: &http.Server{
			Addr:         serverPort,
			Handler:      handler,
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
