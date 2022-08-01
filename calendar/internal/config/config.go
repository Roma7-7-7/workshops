package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

// Application holds application configuration values
type Application struct {
	DB *Database
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  bool   `yaml:"sslmode"`
}

var instance *Application
var once sync.Once

func GetConfig() *Application {
	once.Do(func() {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		viper.SetDefault("db.host", "localhost")
		viper.SetDefault("db.port", "5432")
		viper.SetDefault("db.name", "gotest")
		viper.SetDefault("db.user", "gouser")
		viper.SetDefault("db.password", "gopassword")
		viper.SetDefault("db.sslmode", true)

		instance = &Application{}
		if err := viper.Unmarshal(instance); err != nil {
			panic(err)
		}

		overrideIfPresentS(&instance.DB.Host, "DB_HOST")
		overrideIfPresentS(&instance.DB.Port, "DB_PORT")
		overrideIfPresentS(&instance.DB.Name, "DB_NAME")
		overrideIfPresentS(&instance.DB.User, "DB_USER")
		overrideIfPresentS(&instance.DB.Password, "DB_PASSWORD")
		overrideIfPresentB(&instance.DB.SSLMode, "DB_SSL_MODE")

	})
	return instance
}

func (a *Application) DSN() string {
	res := fmt.Sprintf("user=%s password=%s dbname=%s", a.DB.User, a.DB.Password, a.DB.Name)
	if !a.DB.SSLMode {
		res = res + " sslmode=disable"
	}
	return res
}

func overrideIfPresentS(target *string, key string) {
	if viper.IsSet(key) {
		val := viper.GetString(key)
		*target = val
	}
}

func overrideIfPresentB(target *bool, key string) {
	if viper.IsSet(key) {
		val := viper.GetBool(key)
		*target = val
	}
}
