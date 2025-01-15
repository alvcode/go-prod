package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEVELOPMENT" env-default:"false"`
	Listen        struct {
		Type       string `env:"LISTEN_TYPE" env-default:"port" env-description:"port|sock"`
		BindIP     string `env:"LISTEN_BIND_IPP" env-default:"0.0.0.0"`
		Port       string `env:"LISTEN_PORTT" env-default:"8098"`
		SocketFile string `env:"LISTEN_SOCKET_FILE" env-default:"app.sock"`
	}
	AppConfig struct {
		LogLevel  string `env:"LOG_LEVEL" env-default:"info"`
		AdminUser struct {
			Email    string `env:"ADMIN_EMAIL" env-default:"admin@example.com"`
			Password string `env:"ADMIN_PASSWORD" env-default:"admin"`
		}
	}
	PostgreSQL struct {
		Username string `env:"PGSQL_USER" env-required:"true"`
		Password string `env:"PGSQL_PASSWORD" env-required:"true"`
		Host     string `env:"PGSQL_HOST" env-required:"true"`
		Port     string `env:"PGSQL_PORT" env-required:"true"`
		Database string `env:"PGSQL_DB" env-required:"true"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("Gather config env")

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "ALVCODE - prod app"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
