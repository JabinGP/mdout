package config

import (
	"sync"

	"github.com/JabinGP/mdout/model"
)

var once sync.Once

// Obj 全局配置文件对应实体的实例
var Obj model.Config

func init() {
	once.Do(func() {
		InitConfigFile()
		readConfig()
	})
}
