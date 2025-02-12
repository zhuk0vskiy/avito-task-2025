package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Jwt      JwtConfig      `yaml:"jwt"`
	Http     HTTPConfig     `yaml:"http"`
	Context  ContextConfig  `yaml:"context"`
}

type JwtConfig struct {
	Key         string `yaml:"key" env:"JWT_KEY"`
	ExpTimeHour int    `yaml:"expTimeHour" env:"JWT_EXP_TIME_HOUR"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
	File  string `yaml:"file" env:"LOGGER_FILE" env-default:"out.log"`
}

type HTTPConfig struct {
	Address      string `yaml:"address" env:"HTTP_ADDRESS_" env-default:"localhost"`
	Port         string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout  int    `yaml:"readTimeout" env:"READ_TIMEOUT" env-default:"5"`
	WriteTimeout int    `yaml:"writeTimeout" env:"WRITE_TIMEOUT" env-default:"5"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT" env-default:"5438"`
	User     string
	Password string
	Database string
	Driver   string `yaml:"driver" env:"POSTGRES_DRIVER" env-default:"postgres"`
}

type ContextConfig struct {
	TimeoutSec int `yaml:"timeoutSec"`
}

func NewConfig(envPath string) (*Config, error) {
	var err error
	var config Config

	err = godotenv.Load(envPath)
	if err != nil {
		return nil, err
	}

	configPath := os.Getenv("CONFIG_PATH")

	// viper.SetConfigFile(configPath)

	// err = viper.ReadInConfig()
	// if err != nil {
	// 	return nil, err
	// }
	// err = viper.Unmarshal(&config)
	// if err != nil {
	// 	return nil, err
	// }

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, err
	}

	config.Jwt.Key = os.Getenv("JWT_KEY")
	config.Database.Postgres.User = os.Getenv("POSTGRES_USER")
	config.Database.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	config.Database.Postgres.Database = os.Getenv("POSTGRES_DATABASE")

	return &config, nil

}
