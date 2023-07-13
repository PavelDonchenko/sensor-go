package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config describes all app configuration
type Config struct {
	HTTP struct {
		Host        string        `env-required:"true" yaml:"ip" env:"SERVER_HOST"`
		Port        string        `env-required:"true" yaml:"port" env:"SERVER_PORT"`
		ReadTimeOut time.Duration `env-required:"true" yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	} `yaml:"http"`
	Postgres struct {
		Password    string `env-default:"secret" env-required:"true" yaml:"password" env:"DB_PASSWORD"`
		Username    string `env-default:"root" env-required:"true" yaml:"username" env:"DB_USERNAME"`
		Host        string `env-default:"localhost" env-required:"true" yaml:"host" env:"DB_HOST"`
		Port        string `env-default:"5432" env-required:"true" yaml:"port" env:"DB_PORT"`
		Database    string `env-default:"coa" env-required:"true" yaml:"database" env:"DB_DATABASE"`
		MaxAttempts int    `env-default:"5" env-required:"true" yaml:"attempts" env:"ATTEMPTS"`
	} `yaml:"postgresql"`
	GroupNames         string `env-default:"Alpha, Beta, Gamma" env-required:"true" yaml:"group_names" env:"GROUP_NAMES"`
	CountSensorInGroup int    `env-default:"5" env-required:"true" yaml:"sensors_count" env:"SENSORS_COUNT"`
	SensorGenerated    bool   `env-default:"false" yaml:"sensors_generated" env:"SENSORS_GENERATED"`
}

func GetConfig() *Config {
	log.Print("config init")

	c := &Config{}

	if err := cleanenv.ReadConfig("config.yaml", c); err != nil {
		log.Fatalf("error read config: %v", err)
	}

	return c
}
