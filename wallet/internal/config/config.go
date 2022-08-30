package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"strings"
	"sync"
)

// Application holds application configuration values
type Application struct {
	DB     *Database `env:",prefix=DB_"`
	Bcrypt string    `env:"BCRYPT"`
	JWT    string    `env:"JWT"`
}

type Database struct {
	DSN string `env:"DSN"`
}

var instance *Application
var once sync.Once

func GetConfig() *Application {
	once.Do(initApplicationConfig)
	return instance
}

func initApplicationConfig() {
	instance = &Application{
		DB: &Database{},
	}

	if err := envconfig.Process("", instance); err != nil {
		panic(err)
	}

	missing := make([]string, 0)

	if strings.TrimSpace(instance.DB.DSN) == "" {
		missing = append(missing, "DB_DSN")
	}
	if strings.TrimSpace(instance.Bcrypt) == "" {
		missing = append(missing, "BCRYPT")
	}
	if strings.TrimSpace(instance.JWT) == "" {
		missing = append(missing, "JWT")
	}

	if len(missing) > 0 {
		zap.L().Fatal("missing environment variables", zap.Strings("envs", missing))
	}
}
