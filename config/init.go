package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/JabinGP/mdout/tool"
)

// CheckAndInitBasicConfig 检查配置文件，并初始化
func CheckAndInitBasicConfig() {
	confBytes := []byte(`[Parmas]
Out = ""
Type = "pdf"
Theme = "github"
PageFormat = "a4"
PageOrientation = "portrait"
PageMargin = 0.4`)

	confFolderName, err := filepath.Abs(tool.GetHomeDir() + "/mdout")
	if err != nil {
		log.Println("处理配置文件夹路径失败。")
		panic(err)
	}
	confFileName, err := filepath.Abs(confFolderName + "/conf.toml")
	if err != nil {
		log.Println("处理配置文件夹文件路径失败。")
		panic(err)
	}

	if !tool.IsExists(confFileName) {
		log.Println("配置文件 " + confFileName + " 不存在，创建中...")

		if !tool.IsExists(confFolderName) {
			log.Println("配置文件夹 " + confFolderName + " 不存在，创建中...")
			err := os.Mkdir(confFolderName, os.ModePerm)
			if err != nil {
				fmt.Printf("创建文件夹 " + confFolderName + " 失败!\n")
				panic(err)
			}
			fmt.Printf("创建文件夹 " + confFolderName + " 成功!\n")
		}

		err = ioutil.WriteFile(confFileName, confBytes, os.ModeAppend)
		if err != nil {
			log.Println("创建配置文件 " + confFileName + " 失败，请重新尝试或者手动创建！")
			panic(err)
		}
		log.Println("创建配置文件 " + confFileName + " 成功！")
	}
}
