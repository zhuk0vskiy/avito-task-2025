package intgr

import "avito-task-2025/backend/config"

func NewTestConfig() *config.Config {
	return &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Database: "shop",
				Port:     5438,
				Host:     "localhost",
				User:     "admin",
				Password: "avito",
				Driver:   "postgres",
			},
		},
		JwtKey: "oh god i love avito",
	}
	// return &config.DatabaseConfig{

	// }
}
