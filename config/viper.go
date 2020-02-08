package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/JabinGP/mdout/tool"
	"github.com/JabinGP/mdout/types"
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
		fmt.Println("读取配置文件成功:", v.ConfigFileUsed())
		ShowConfig(v)
	} else {
		log.Printf("读取配置文件失败: %s \n", err)
		log.Println("将以系统预设值作为参数默认值")
	}
}

// ShowConfig 输出读取到的配置文件
func ShowConfig(v *viper.Viper) {
	confParmas := types.Parmas{}
	v.UnmarshalKey("Parmas", &confParmas)
	log.Println("---这是你的配置文件参数---")
	log.Printf("输出路径：%s\n", confParmas.Out)
	log.Printf("输出格式：%s\n", confParmas.Type)
	log.Printf("选择主题：%s\n", confParmas.Theme)
	log.Printf("打印页面格式：%s\n", confParmas.PageFormat)
	log.Printf("打印页面方向：%s\n", confParmas.PageOrientation)
	log.Printf("打印页面边距：%s\n", confParmas.PageMargin)
	log.Println("-----------------------")
}
