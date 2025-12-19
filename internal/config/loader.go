package config

import (
	"log"

	"github.com/spf13/viper"
)

var AppCfgs *AppConfigs
func Load() {
	viper.SetConfigFile("../../internal/config/config.yml")
	if err:=viper.ReadInConfig();err!=nil{
		log.Println("Viper config load error: ",err.Error())
		return
	}

	AppCfgs=&AppConfigs{}
	if err:=viper.Unmarshal(AppCfgs);err!=nil{
		log.Println("Config unmarshal error: ",err.Error())
		return
	}
	log.Println(AppCfgs)
}