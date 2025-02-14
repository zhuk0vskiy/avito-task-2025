package helper

import (
	"avito-task-2025/backend/config"
	"fmt"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewTestConfig(migrationsDir string) (*migrate.Migrate, *config.Config) {

	database := "test_shop"
	host := "localhost"
	port := 9999
	user := "test_admin"
	password := "test_avito"
	driver := "postgres"

	m, err := migrate.New(
		migrationsDir,
		fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", driver, user, password, host, strconv.Itoa(port), database))
	if err != nil {
		panic(err)
	}
	return m, &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Database: database,
				Port:     port,
				Host:     host,
				User:     user,
				Password: password,
				Driver:   driver,
			},
		},
		Jwt: config.JwtConfig{
			Key:         "testAvtitoKey",
			ExpTimeHour: 1,
		},
	}

}
