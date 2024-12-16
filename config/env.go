package config

import (
	"github.com/spf13/viper"
)

func Env(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return viper.GetString(key)
}