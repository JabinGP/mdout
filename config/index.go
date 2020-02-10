package config

import (
	"sync"

	"github.com/JabinGP/mdout/model"
	"github.com/spf13/viper"
)

var once sync.Once

// Viper 全局配置文件实例
var Viper *viper.Viper

// Obj 全局配置文件对应实体的实例
var Obj model.Config

func init() {
	once.Do(func() {
		InitConfigFile()
		initViper()
		initObj()
	})
}

func initObj() {
	Viper.Unmarshal(&Obj)
}
