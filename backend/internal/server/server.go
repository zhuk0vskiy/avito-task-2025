package server

import (
	"avito-task-2025/backend/config"
	"avito-task-2025/backend/internal/app"
	"avito-task-2025/backend/internal/controller/http/v1/chi/auth"
	"avito-task-2025/backend/internal/controller/http/v1/chi/coin"
	"avito-task-2025/backend/internal/controller/http/v1/chi/merch"
	"avito-task-2025/backend/internal/controller/http/v1/chi/user"
	ginServer "avito-task-2025/backend/internal/controller/http/v1/gin"

	// "avito-task-2025/backend/internal/controller"
	"fmt"

	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	// "github.com/go-chi/cors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"

	// "github.com/rs/zerolog/log"
	"log"
)

type Server struct {
	httpServer *http.Server
}

func NewGinServer(cfg config.HTTPConfig, jwtAuth *jwtauth.JWTAuth, a *app.App) *Server {
	handler := gin.New()
	con := ginServer.NewRouter(handler)

	// handler.Use(func(c *gin.Context) {

	// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 	ctx = context.WithValue(ctx, "Authorization", c.GetHeader("Authorization"))
	// 	defer cancel()

	// 	c.Request = c.Request.WithContext(ctx)
	// 	// c.Request.Context()
	// 	c.Next()
	// })

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

func NewChiServer(cfg config.HTTPConfig, jwtAuth *jwtauth.JWTAuth, a *app.App) *Server {
	mux := chi.NewMux()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	// mux.Use(middleware.Logger)
	mux.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/auth", auth.SignInHandler(a))
		})
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwtAuth))
			// r.Use(jwtauth.Authenticator(jwtAuth))
			r.Get("/info", user.GetUserInfoHandler(a))
			r.Post("/sendCoin", coin.SendCoinHandler(a))
			r.Get("/buy/{item}", merch.BuyMerchHandler(a))
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
		log.Println("error while listening:", err)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("stopping server")
	return s.httpServer.Shutdown(ctx)
}
