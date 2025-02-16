package main

import (
	"avito-task-2025/backend/config"
	"avito-task-2025/backend/internal/app"
	"avito-task-2025/backend/internal/server"
	"avito-task-2025/backend/internal/storage/postgres"
	"avito-task-2025/backend/pkg/logger"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

const envPath = "./.env"

func main() {
	c, err := config.NewConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.Jwt.Key)
	loggerFile, err := os.OpenFile(
		c.Logger.File,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Fatal(err)
	}
	l := logger.New(c.Logger.Level, loggerFile)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbConn, err := postgres.NewDbConn(ctx, &c.Database.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	tokenAuth := jwtauth.New("HS256", []byte(c.Jwt.Key), nil)

	a := app.NewApp(c, l, dbConn)


	s := server.NewChiServer(c.Http, tokenAuth, a)

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	go s.Start()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.Stop(ctx)
	if err != nil {
		log.Fatal("server shutdown:", err)
	}

	<-ctx.Done()

	log.Println("server exiting")
}
