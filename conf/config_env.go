package conf

import (
	"github.com/caarlos0/env/v6"
)

// AppConfig presents app conf
type AppConfig struct {
	Port    string `env:"PORT" envDefault:"5050"`
	EnvName string `env:"ENV_NAME" envDefault:"dev"`

	//DB CONFIG
	DBHost        string   `env:"DB_HOST" envDefault:"localhost"`
	DBPort        string   `env:"DB_PORT" envDefault:"3306"`
	DBUser        string   `env:"DB_USER" envDefault:"doan"`
	DBPass        string   `env:"DB_PASS" envDefault:"doan"`
	DBName        string   `env:"DB_NAME" envDefault:"doan"`
	DBSchema      string   `env:"DB_SCHEMA" envDefault:"public"`
	Env           string   `env:"ENV" envDefault:"dev"`
	DebugPort     int      `env:"DEBUG_PORT" envDefault:"7070"`
	ReadTimeout   int      `env:"READ_TIMEOUT" envDefault:"15"`
	EnableProfile bool     `env:"ENABLE_PROFILE" envDefault:"true"` // enable profile listener
	EnableDB      bool     `env:"ENABLE_DB" envDefault:"true"`
	TrustedProxy  []string `env:"TRUSTED_PROXY" envSeparator:"," envDefault:"127.0.0.1,10.0.0.0/8,192.168.0.0/16"`
	Debug         bool     `env:"DEBUG" envDefault:"true"`
	DB            *Config

	// server config
	MovieOnAi string `env:"MOVIE_ON_AI" envDefault:"http://localhost:8000"`
}

var config AppConfig

func SetEnv() {
	_ = env.Parse(&config)
}

func GetEnv() AppConfig {
	return config
}
