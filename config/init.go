package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/JabinGP/mdout/tool"
)

// CheckAndInitBasicConfig 检查配置文件，并初始化
func CheckAndInitBasicConfig() {
	confBytes := []byte(`[Parmas]
Out = ""                      # 输出路径，一般通过命令行指定
Type = "pdf"                  # 输出类型默认PDF
Theme = "github"              # 输出默认主题
PageFormat = "a4"             # 输出默认尺寸
PageOrientation = "portrait"  # 输出默认方向
PageMargin = 0.4              # 输出默认边距
ExecPath = ""                 # chrome执行文件路径
[Runtime]
EditorPath = "code"           # 打开配置文件的编辑器路径或命令`)

	confFolderName, err := tool.Abs(tool.GetHomeDir() + "/mdout")
	if err != nil {
		log.Println("处理配置文件夹路径失败。")
		panic(err)
	}
	confFileName, err := tool.Abs(confFolderName + "/conf.toml")
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

		err = ioutil.WriteFile(confFileName, confBytes, 0777)
		if err != nil {
			log.Println("创建配置文件 " + confFileName + " 失败，请重新尝试或者手动创建！")
			panic(err)
		}
		log.Println("创建配置文件 " + confFileName + " 成功！")
	}
}
