package config

import (
	"github.com/spf13/viper"
)

func Get(key string) interface{} {
	return viper.Get(key)
}
