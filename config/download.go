package config

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JabinGP/mdout/log"
	"github.com/JabinGP/mdout/static"
)

// DownloadConfig 从指定连接下载配置文件
func DownloadConfig(configURL string) error {
	// 获取
	configReso, err := http.Get(configURL)
	if err != nil {
		return err
	}
	if configReso.StatusCode != 200 {
		return fmt.Errorf("从 %s 下载配置文件失败，HTTP 状态码 %d。", configURL, configReso.StatusCode)
	}
	defer configReso.Body.Close()
	confBytes, err := ioutil.ReadAll(configReso.Body)
	if err != nil {
		return err
	}

	// 保存
	err = ioutil.WriteFile(static.ConfigFileFullName, confBytes, 0777)
	if err != nil {
		log.Errorf("创建配置文件 " + static.ConfigFileFullName + " 失败，请重新尝试或者手动创建！")
		return err
	}
	log.Infoln("创建配置文件 " + static.ConfigFileFullName + " 成功！")
	return nil
}
