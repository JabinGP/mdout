package cmd

import (
	// 读写文件
	"io"
	// 输出日志
	"log"

	// http包，用于请求文件
	"net/http"
	// 获取系统路径信息
	"os"

	// 命令行框架
	"github.com/spf13/cobra"

	dir "github.com/JabinGP/mdout/dir"
	ziper "github.com/JabinGP/mdout/ziper"
)

// InstallCmd 返回一个命令，在root.go中加入根命令
func InstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "下载运行时必要的资源到本地",
		Long:  "下载配置文件，主题包并保存到用户家目录下的binmdout文件夹",
		RunE:  installRunE,
	}
}

// installRunE install命令执行时运行
func installRunE(cmd *cobra.Command, args []string) error {
	// 资源文件所在地址
	url := `http://112.74.177.253:8000/f/382a2291375045cb81fa/?dl=1`

	// 获取文件保存地址和文件名
	homeDir := dir.HomeDir()
	downFilePath := homeDir + "/mdout-downolad.zip"
	log.Println("开始从" + url + "下载资源")

	// 打开文件
	file, err := os.Create(downFilePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	// 从url获取资源
	urlFileResp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer urlFileResp.Body.Close()

	// 将获取资源复制到文件
	_, err = io.Copy(file, urlFileResp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	// 校验文件大小
	// if n != 923637 {
	// 	log.Printf("下载文件出错，请重试，文件大小为%d", n)
	// }

	// 复制完成，关闭文件
	file.Close()
	urlFileResp.Body.Close()

	log.Println("下载成功，保存到" + downFilePath)
	log.Println("开始解压")

	// 解压缩到binmdout文件
	unzipFolder := homeDir + "/binmdout"
	err = ziper.UnZip(downFilePath, unzipFolder)
	if err != nil {
		log.Println("解压失败")
		log.Println(err)
		return err
	}
	log.Println("解压成功，解压到" + unzipFolder)

	log.Println("开始删除多余压缩文件" + downFilePath)

	if dir.IsExists(downFilePath) {
		log.Println("文件存在，正在删除" + downFilePath)
		err := os.Remove(downFilePath)
		if err != nil {
			log.Println("删除文件" + downFilePath + "失败，但不会影响使用")
			log.Println(err)
			return err
		}
	}

	log.Println("删除" + downFilePath + "成功")
	return nil
}
