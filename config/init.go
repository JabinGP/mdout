package config

import (
	"encoding/json"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
)

func initConfigByDefault() {
	Obj.Parmas.OutType = "pdf"
	Obj.Parmas.ThemeName = "github"
	Obj.Parmas.PageFormat = "a4"
	Obj.Parmas.PageOrientation = "portrait"
	Obj.Parmas.PageMargin = "0.4"

	Obj.Runtime.EditorPath = "code"
	Obj.Runtime.StdoutLogLevel = "debug"
	Obj.Runtime.FileLogLevel = "debug"
	Obj.Runtime.EnableXHTMLOutput = true
	Obj.Runtime.EnableHTMLTag = true
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
	if !tool.IsExists(static.ConfigFolderFullName) {
		log.Infoln("配置文件夹 " + static.ConfigFolderFullName + " 不存在，创建中...")
		err := os.Mkdir(static.ConfigFolderFullName, os.ModePerm)
		if err != nil {
			log.Errorf("创建文件夹 " + static.ConfigFolderFullName + " 失败!\n")
			panic(err)
		}
		log.Infoln("创建文件夹 " + static.ConfigFolderFullName + " 成功!\n")
	}
}
