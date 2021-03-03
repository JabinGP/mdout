package config

import (
	"sync"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
)

var once sync.Once

// Obj 全局配置文件对应实体的实例
var Obj model.Config

func init() {
	once.Do(func() {
		InitConfigFileFolder()
		InitThemeFolder()
		if !tool.IsExists(static.ConfigFileFullName) {
			log.Infof("配置文件 %s 不存在，使用内置默认参数。", static.ConfigFileFullName)
			initConfigByDefault()
			return
		}
		readConfig()
	})
}
