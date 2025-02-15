package config

import (
	"fmt"
	"log"
	"os"

	// "log"

	// "github.com/joho/godotenv"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Jwt      JwtConfig      `yaml:"jwt"`
	Http     HTTPConfig     `yaml:"http"`
}

type JwtConfig struct {
	Key         string `yaml:"key" env:"JWT_KEY"`
	ExpTimeHour int    `yaml:"expTimeHour"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOGS_LEVEL" env-default:"warn"`
	File  string `yaml:"file" env:"LOGS_FILE"`
}

type HTTPConfig struct {
	Address      string `yaml:"address" env:"HTTP_ADDRESS" env-default:"localhost"`
	Port         string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT"`
	User     string `yaml:"user" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
	Database string `yaml:"database" env:"POSTGRES_DATABASE"`
	Driver   string `yaml:"driver" env:"POSTGRES_DRIVER" env-default:"postgres"`
}

func NewConfig(envPath string) (*Config, error) {
	var err error
	var config Config

	err = godotenv.Load(envPath)
	if err != nil {
		log.Println(".env not find, using os environment")
	}

	configPath := os.Getenv("CONFIG_PATH")

	fmt.Println(configPath)

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

	// config.Jwt.Key = os.Getenv("JWT_KEY")
	// config.Database.Postgres.User = os.Getenv("POSTGRES_USER")
	// config.Database.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	// config.Database.Postgres.Database = os.Getenv("POSTGRES_DATABASE")
	// config.Database.Postgres.Host = os.Getenv("POSTGRES_HOST")

	// postgresPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	// if err != nil {
	// 	return nil, err
	// }
	// config.Database.Postgres.Port = postgresPort
	// fmt.Println(config.Database.Postgres.Port)
	return &config, nil

}
