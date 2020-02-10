package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/static"
	"github.com/JabinGP/mdout/tool"
	"github.com/spf13/viper"
)

// NewViper 获取一个初始化完成的viper实例
func initViper() {
	Viper = viper.New()
	addConfigPathAndName(Viper)
	setDefault(Viper)
	readConfig(Viper)
}

func addConfigPathAndName(v *viper.Viper) {
	// 添加扫描路径
	v.AddConfigPath(static.ConfigFolderFullName)
	// 设置配置文件名称
	v.SetConfigName(static.ConfigFileViperName)
}

func setDefault(v *viper.Viper) {
	v.SetDefault("Out", "")
	v.SetDefault("Type", "pdf")
	v.SetDefault("Theme", "default")
	v.SetDefault("PageFormat", "a4")
	v.SetDefault("PageOrientation", "portrait")
	v.SetDefault("PageMargin", "0.4")
}

func readConfig(v *viper.Viper) {
	if err := v.ReadInConfig(); err == nil {
		log.Println("读取配置文件成功:", v.ConfigFileUsed())
		ShowConfig(v)
	} else {
		log.Printf("读取配置文件失败: %s \n", err)
		panic(err)
	}
}

// ShowConfig 输出读取到的配置文件
func ShowConfig(v *viper.Viper) {
	conf := model.Config{}
	v.Unmarshal(&conf)
	confBts, err := json.Marshal(conf)
	if err != nil {
		log.Println(err)
	}
	log.Println("---这是你的配置文件参数---")
	log.Println(string(confBts))
	log.Println("--------------------------")
}

// InitConfigFile 初始化配置文件
func InitConfigFile() {
	if !tool.IsExists(static.ConfigFileFullName) {
		log.Println("配置文件 " + static.ConfigFileFullName + " 不存在，创建中...")
		confBytes := DownloadConfig("v1")

		if !tool.IsExists(static.ConfigFolderFullName) {
			log.Println("配置文件夹 " + static.ConfigFolderFullName + " 不存在，创建中...")
			err := os.Mkdir(static.ConfigFolderFullName, os.ModePerm)
			if err != nil {
				fmt.Printf("创建文件夹 " + static.ConfigFolderFullName + " 失败!\n")
				panic(err)
			}
			fmt.Printf("创建文件夹 " + static.ConfigFolderFullName + " 成功!\n")
		}

		err := ioutil.WriteFile(static.ConfigFileFullName, confBytes, 0777)
		if err != nil {
			log.Println("创建配置文件 " + static.ConfigFileFullName + " 失败，请重新尝试或者手动创建！")
			panic(err)
		}
		log.Println("创建配置文件 " + static.ConfigFileFullName + " 成功！")
	}
}
