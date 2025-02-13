package intgr

import "avito-task-2025/backend/config"



func NewTestConfig() *config.Config {

	return &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Database: "test_shop",
				Port:     5438,
				Host:     "localhost",
				User:     "admin",
				Password: "avito",
				Driver:   "postgres",
			},
		},
		Jwt: config.JwtConfig{
			Key: "testAvtitoKey",
			ExpTimeHour: 1,
		},
	}
	// return &config.DatabaseConfig{

	// }
}

func TestRun(t *testing.T) {
	var address string
	var port string
	go BeforeAll(&address, &port)
	baseUrl := fmt.Sprintf("%s:%s", address, port)

	t.Run("SignUp Success 01", func(t *testing.T) {
	
		SignInSuccess_01(t, baseUrl)
	})
}


