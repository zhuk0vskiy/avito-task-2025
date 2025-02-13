package e2e

import (
	"avito-task-2025/backend/config"
	"avito-task-2025/backend/internal/app"
	"avito-task-2025/backend/internal/server"
	"avito-task-2025/backend/internal/storage/postgres"
	loggerMock "avito-task-2025/backend/pkg/logger/mocks"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

const envPath = "./.env"

func BeforeAll(address *string, port *string) {
	c, err := config.NewConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}

	*address = c.Http.Address
	*port = c.Http.Port

	mockLogger := new(loggerMock.Interface)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbConn, err := postgres.NewDbConn(ctx, &c.Database.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	tokenAuth := jwtauth.New("HS256", []byte(c.Jwt.Key), nil)

	a := app.NewApp(c, mockLogger, dbConn)

	s := server.NewServer(c.Http, tokenAuth, a)

	s.Start()
}

func TestRunAllTest(t *testing.T) {
	var address string
	var port string
	go BeforeAll(&address, &port)
	baseUrl := fmt.Sprintf("%s:%s", address, port)

	t.Run("SignUp Success 01", func(t *testing.T) {
	
		SignInSuccess_01(t, baseUrl)
	})
}
