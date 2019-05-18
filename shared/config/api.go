package config

import (
	"github.com/spf13/viper"
)

var configuration = "dev"

func Get(key string) interface{} {
	return viper.Get(configuration + "." + key)
}

func GetInt(key string) int {
	return viper.GetInt(configuration + "." + key)
}

func GetString(key string) string {
	return viper.GetString(configuration + "." + key)
}

func GetBool(key string) bool {
	return viper.GetBool(configuration + "." + key)
}
