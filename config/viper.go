package config

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/JabinGP/mdout/model"
	"github.com/JabinGP/mdout/tool"
	"github.com/spf13/viper"
)

var once sync.Once

// NewViper 获取一个初始化完成的viper实例
func NewViper() *viper.Viper {
	viper := viper.New()
	initConfig(viper)
	setDefault(viper)
	readConfig(viper)
	return viper
}

func initConfig(v *viper.Viper) {
	// 添加扫描路径
	configPath := tool.GetHomeDir() + "/mdout"
	v.AddConfigPath(configPath)

	// 设置配置文件名称
	v.SetConfigName("conf")
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
	CheckAndInitBasicConfig()
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
