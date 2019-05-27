package config

import (
	"github.com/spf13/viper"
	"time"
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

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(configuration + "." + key)
}

func GetBool(key string) bool {
	return viper.GetBool(configuration + "." + key)
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {
	duration, err := time.ParseDuration(GetString(key))
	if err != nil {
		duration = defaultValue
	}
	return duration
}
