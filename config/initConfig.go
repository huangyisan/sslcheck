package config

import "github.com/spf13/viper"

var Conf *viper.Viper

func init() {
	conf := viper.New()
	conf.AddConfigPath("config")
	conf.SetConfigName("conf.yml")
	conf.SetConfigType("yaml")
	if err := conf.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	Conf = conf
}
