package config

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/JabinGP/mdout/static"
)

// DownloadConfig 从github下载配置文件
func DownloadConfig(version string) []byte {
	var configURL = static.ConfigGithubURL
	configURL = strings.Replace(configURL, "{version}", version, 1)
	configReso, err := http.Get(configURL)
	if err != nil {
		PublicLogger.Errorln("从 " + configURL + " 下载配置文件失败！")
		panic(err)
	}
	defer configReso.Body.Close()
	configBts, err := ioutil.ReadAll(configReso.Body)
	if err != nil {
		PublicLogger.Errorln("从 " + configURL + " 读取响应内容失败！")
		panic(err)
	}
	PublicLogger.Infoln("从 " + configURL + " 下载配置文件成功！")
	return configBts
}
