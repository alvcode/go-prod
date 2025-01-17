package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	//IsDevelopment bool `yaml:"is_development" env:"IS_DEVELOPMENT" env-default:"false"`
	HTTP struct {
		IP           string        `yaml:"ip" env:"HTTP_IP" env-default:"127.0.0.1"`
		Port         int           `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
		ReadTimeout  time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"15"`
		WriteTimeout time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"15"`
		CORS         struct {
			AllowedMethods     []string `yaml:"allowed_methods" env:"HTTP_CORS_ALLOWED_METHODS"`
			AllowedOrigins     []string `yaml:"allowed_origins" env:"HTTP_CORS_ALLOWED_ORIGINS"`
			AllowedHeaders     []string `yaml:"allowed_headers" env:"HTTP_CORS_ALLOWED_HEADERS"`
			AllowCredentials   bool     `yaml:"allow_credentials" env:"HTTP_CORS_ALLOW_CREDENTIALS"`
			OptionsPassthrough bool     `yaml:"options_passthrough" env:"HTTP_CORS_OPTIONS_PASSWORD"`
			ExposedHeaders     []string `yaml:"exposed_headers" env:"HTTP_CORS_EXPOSED_HEADERS"`
			Debug              bool     `yaml:"debug" env:"HTTP_CORS_DEBUG" env-default:"false"`
		} `yaml:"cors"`
	} `yaml:"http"`
	AppConfig struct {
		IsDebug   bool   `yaml:"is_debug" env:"IS_DEBUG" env-default:"false"`
		LogLevel  string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
		AdminUser struct {
			Email    string `yaml:"email" env:"ADMIN_EMAIL" env-default:"admin@example.com"`
			Password string `yaml:"password" env:"ADMIN_PASSWORD" env-default:"admin"`
		} `yaml:"admin_user"`
	} `yaml:"app_config"`
	PostgreSQL struct {
		Host     string `yaml:"host" env:"PGSQL_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"PGSQL_PORT" env-required:"true"`
		Username string `yaml:"username" env:"PGSQL_USER" env-required:"true"`
		Password string `yaml:"password" env:"PGSQL_PASSWORD" env-required:"true"`
		Database string `yaml:"database" env:"PGSQL_DB" env-required:"true"`
	} `yaml:"postgresql"`
}

const (
	EnvConfigPathName  = "CONFIG_PATH"
	FlagConfigPathName = "config"
)

var configPath string
var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		flag.StringVar(&configPath, FlagConfigPathName, "config/config.local.yaml", "this is app config file")
		flag.Parse()

		log.Print("Config init")

		if configPath == "" {
			configPath = os.Getenv(EnvConfigPathName)
		}

		if configPath == "" {
			log.Fatal("config file path is empty")
		}

		instance = &Config{}

		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			helpText := "ALVCODE - prod app"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
