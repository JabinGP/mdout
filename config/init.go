package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
)

// CheckConfigFile 检查配置文件，不存在则自动创建
func CheckConfigFile() {
	if !tool.IsExists(static.ConfigFileFullName) {
		log.Infof("配置文件 %s 不存在，将自动创建默认配置文件。", static.ConfigFileFullName)
		CreateDefaultConfig()
		log.Infof("自动创建默认配置文件 %s 成功。", static.ConfigFileFullName)
	}
}

// CreateDefaultConfig 根据 static 中的 TomlConfig 自动创建配置文件
func CreateDefaultConfig() {
	autoCreateInfo := fmt.Sprintf("# 由 mdout-%s 自动创建于 %s\n", static.Version, time.Now().Format("2006-01-02 15:04:05"))
	err := ioutil.WriteFile(static.ConfigFileFullName, []byte(autoCreateInfo+static.TomlConfig), 0777)
	if err != nil {
		log.Errorln("创建默认配置文件失败！", err)
		panic(err)
	}
}

func readConfig() {
	if _, err := toml.DecodeFile(static.ConfigFileFullName, &Obj); err != nil {
		log.Errorln("读取配置文件失败！", err)
		panic(err)
	}
}

// ShowConfig 输出读取到的配置文件
func ShowConfig() {
	confBts, err := json.MarshalIndent(Obj, "", "\t")
	if err != nil {
		log.Errorln(err)
	}
	log.Debugln("---这是你的配置文件参数---")
	log.Debugln(string(confBts))
	log.Debugln("--------------------------")
}

// InitConfigFileFolder 初始化配置文件夹
func InitConfigFileFolder() {
	if err := tool.InitFolder(static.ConfigFolderFullName); err != nil {
		log.Errorln(err)
		panic(err)
	}
}

// InitThemeFolder 初始化主题文件夹
func InitThemeFolder() {
	if err := tool.InitFolder(static.ThemeFolderFullName); err != nil {
		log.Errorln(err)
		panic(err)
	}
}
