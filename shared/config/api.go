package config

import (
	"github.com/spf13/viper"
)

var configuration = "dev"

func Get(key string) interface{} {
	return viper.Get(configuration + "." + key)
}
